import re
import base64

from elasticsearch import AsyncElasticsearch
from elasticsearch.exceptions import NotFoundError
from sentence_transformers import SentenceTransformer
import json
import spacy

import src.model as model


class Service:
    def __init__(
        self,
        es: AsyncElasticsearch,
    ):
        self.es = es
        self.index_name = "semantic_index"

        self.nlp = spacy.load("ru_core_news_sm", disable=["parser", "tagger", "ner"])
        self.nlp_stops = self.nlp.Defaults.stop_words

        self.model = SentenceTransformer("all-mpnet-base-v2")

    async def start(self):
        mappings = {
            "properties": {
                "text_vector": {
                    "type": "dense_vector",
                    "dims": 768,
                    "index": True,
                    "similarity": "cosine",
                    "index_options": {
                        "type": "int8_hnsw",
                        "m": 16,
                        "ef_construction": 100,
                    },
                }
            }
        }

        settings = {"mappings": mappings}

        try:
            await self.es.indices.get(index=self.index_name)
        except NotFoundError:
            await self.es.indices.create(index=self.index_name, body=settings)

        await self.es.indices.put_mapping(index=self.index_name, body=mappings)

    async def search_similar(self, doc: model.Document) -> list[model.SimilarDocument]:
        text = base64.b64decode(doc.text.decode()).decode()

        text = self._normalize(text)

        text = self._vectorize(text)

        resp = await self.es.search(
            index=self.index_name,
            knn={"field": "text_vector", "query_vector": text, "k": 100},
        )

        return [
            model.SimilarDocument(
                **{
                    "id": item["_id"],
                    "rate": item["_score"],
                    "highlights": [],
                    "similarImages": [],
                }
            )
            for item in resp.body["hits"]["hits"]
        ]

    async def index_content(self, doc: model.Document):
        text = base64.b64decode(doc.text.decode()).decode()

        text = self._normalize(text)

        text = self._vectorize(text)

        await self.es.index(
            index=self.index_name, document={"text_vector": text}, id=doc.id
        )

    def _vectorize(self, text: str):
        return self.model.encode([text])[0]

    def _normalize(self, text: str) -> str:
        text = re.sub(r"[^\w\s.]", "", text)
        text = re.sub(r"\s+", " ", text)

        doc = self.nlp(text)

        lemmatized = list()

        for word in doc:
            lemma = word.lemma_.strip()
            if lemma and lemma not in self.nlp_stops:
                lemmatized.append(lemma)

        return " ".join(lemmatized)
