import asyncio
import os
import json
import logging

from nats.aio.client import Client as NATS
import pydantic


class Document(pydantic.BaseModel):
    text: bytes


class SimilarDocument(pydantic.BaseModel):
    id: str
    rate: float
    highlights: list[str]
    similarImages: list[str]


class Service:
    def __init__(
        self,
        redis_url: str,
    ):
        pass

    def search_similar(self, doc: Document) -> list[SimilarDocument]:
        print(f"got document on search_similar: {doc}")
        return []

    def index_content(self, doc: Document):
        print(f"got document on index_content: {doc}")
        pass


class Handler:
    def __init__(
        self,
        nats_url: str,
        service: Service,
    ):
        self.nats_url = nats_url
        self.nc = NATS()

        self.service = service

        self.document_by_id_prefix = "document.byid"

        self.document_status_stream = "documentstatus"
        self.document_indexer_subject = "documentstatus.document_parsed"

        self.document_similarity_subject = "document.similarity.*"

    async def start(self):
        await self.nc.connect(
            self.nats_url,
        )

        await self.nc.subscribe(
            self.document_similarity_subject, cb=self.similarity_handler()
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

    async def _fetch_document(self, id: str) -> Document:
        doc = await self.nc.request(self.document_by_id_prefix, id.encode(), timeout=5)

        raw = json.loads(doc.data)

        return Document(**raw)


async def main(loop):
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    redis_url = os.getenv("REDIS_DSN", "")

    root = logging.getLogger()
    root.setLevel(logging.INFO)

    service = Service(redis_url)

    nats_handler = Handler(nats_url, service)

    await nats_handler.start()


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main(loop))
    try:
        loop.run_forever()
    finally:
        loop.close()
