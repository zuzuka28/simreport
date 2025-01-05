import asyncio
import os
import json
import logging
import contextlib
import signal

import nats
import nats.micro

from nats.aio.client import Client as NATS

import service
import model


class ServiceHandler:
    def __init__(
        self,
        nc: NATS,
        service: service.Service,
    ):
        self.nc = nc
        self.service = service

        self.document_by_id_subject = "document.byid"

    def similarity_handler(self):
        async def callback(msg):
            id = msg.subject.split(".")[:-1]

            doc = await self._fetch_document(id)

            result = self.service.search_similar(doc)

            await msg.respond(result)

        return callback

    def indexer_handler(self):
        async def callback(msg):
            id: bytes = msg.data

            doc = await self._fetch_document(id.decode())

            self.service.index_content(doc)

            await msg.ack_sync()

        return callback

    async def _fetch_document(self, id: str) -> model.Document:
        doc = await self.nc.request(
            self.document_by_id_subject, id.encode(), timeout=60
        )

        raw = json.loads(doc.data)

        return model.Document(**raw)


async def main():
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    redis_url = os.getenv("REDIS_DSN", "")

    root = logging.getLogger()
    root.setLevel(logging.INFO)

    srv = service.Service(redis_url)

    quit_event = asyncio.Event()

    loop = asyncio.get_event_loop()
    for sig in (signal.Signals.SIGINT, signal.Signals.SIGTERM):
        loop.add_signal_handler(sig, lambda *_: quit_event.set())

    async with contextlib.AsyncExitStack() as stack:
        nc = await stack.enter_async_context(await nats.connect(nats_url))

        srv_handler = ServiceHandler(nc, srv)

        ncs = await stack.enter_async_context(
            await nats.micro.add_service(nc, name="shingleindex", version="0.0.1")
        )

        group = ncs.add_group(name="similarity.shingleindex")

        await group.add_endpoint(
            name="search",
            handler=srv_handler.similarity_handler(),
        )

        await quit_event.wait()


if __name__ == "__main__":
    asyncio.run(main())
