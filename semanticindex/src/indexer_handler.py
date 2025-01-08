import asyncio
import os
import json
import logging

from nats.aio.client import Client as NATS
from elasticsearch import AsyncElasticsearch

import src.service as service
import src.model as model


class IndexerHandler:
    def __init__(
        self,
        nats_url: str,
        service: service.Service,
    ):
        self.nats_url = nats_url
        self.nc = NATS()

        self.service = service

        self.document_by_id_subject = "document.byid"

        self.document_status_stream = "documentstatus"
        self.document_indexer_subject = "documentstatus.document_parsed"

    async def start(self):
        await self.nc.connect(
            self.nats_url,
        )

        await self.nc.jetstream().subscribe(
            self.document_indexer_subject,
            "shingleindex_injest",
            cb=self.indexer_handler(),
        )

        logging.info("nats handlers started")

    async def stop(self):
        await self.nc.close()
        logging.info("nats handlers stopped")

    def indexer_handler(self):
        async def callback(msg):
            id: bytes = msg.data

            doc = await self._fetch_document(id.decode())

            await self.service.index_content(doc)

            await msg.ack_sync()

        return callback

    async def _fetch_document(self, id: str) -> model.Document:
        doc = await self.nc.request(self.document_by_id_subject, id.encode(), timeout=5)

        raw = json.loads(doc.data)

        return model.Document(**raw)


async def main(loop):
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    es_host = os.getenv("ES_HOST", "http://localhost:9200")

    es = AsyncElasticsearch(hosts=[es_host])

    root = logging.getLogger()
    root.setLevel(logging.INFO)

    srv = service.Service(
        es,
    )

    await srv.start()

    nats_handler = IndexerHandler(nats_url, srv)

    await nats_handler.start()


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(loop))
    try:
        loop.run_forever()
    finally:
        loop.close()
