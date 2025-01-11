import asyncio
import os
import json
import logging
import contextlib
import signal

import nats
import nats.micro

from nats.aio.client import Client as NATS
from elasticsearch import AsyncElasticsearch
from pydantic.json import pydantic_encoder

import src.service as service
import src.model as model


class NotFoundException(BaseException):
    pass


class InternalException(BaseException):
    pass


class ServiceHandler:
    nats_status_header = "Nats-Service-Error-Code"

    def __init__(
        self,
        nc: NATS,
        service: service.Service,
    ):
        self.nc = nc
        self.service = service

        self.document_by_id_subject = "document.byid"

    def similarity_handler(self):
        async def callback(msg: nats.micro.request.Request):
            id = msg.data.decode()

            try:
                doc = await self._fetch_document(id)
            except NotFoundException:
                await msg.respond_error("404", "document not found")
                return
            except BaseException:
                await msg.respond_error("500", "internal error")
                return

            result = await self.service.search_similar(doc)

            await msg.respond(json.dumps(result, default=pydantic_encoder).encode())

        return callback

    async def _fetch_document(self, id: str) -> model.Document:
        resp = await self.nc.request(
            self.document_by_id_subject, id.encode(), timeout=60
        )

        if resp.headers and resp.headers.get(self.nats_status_header, None) is not None:
            if resp.headers.get(self.nats_status_header, None) == "404":
                raise NotFoundException()
            else:
                raise InternalException()

        raw = json.loads(resp.data)

        return model.Document(**raw)


async def main():
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    es_host = os.getenv("ES_HOST", "http://localhost:9200")

    es = AsyncElasticsearch(hosts=[es_host])

    root = logging.getLogger()
    root.setLevel(logging.INFO)

    srv = service.Service(
        es,
    )

    await srv.start()

    quit_event = asyncio.Event()

    loop = asyncio.get_event_loop()
    for sig in (signal.Signals.SIGINT, signal.Signals.SIGTERM):
        loop.add_signal_handler(sig, lambda *_: quit_event.set())

    async with contextlib.AsyncExitStack() as stack:
        nc = await stack.enter_async_context(await nats.connect(nats_url))

        srv_handler = ServiceHandler(nc, srv)

        ncs = await stack.enter_async_context(
            await nats.micro.add_service(nc, name="semanticindex", version="0.0.1")
        )

        group = ncs.add_group(name="similarity.semanticindex")

        await group.add_endpoint(
            name="search",
            handler=srv_handler.similarity_handler(),
        )

        logging.info("nats handlers started")

        await quit_event.wait()


if __name__ == "__main__":
    asyncio.run(main())
