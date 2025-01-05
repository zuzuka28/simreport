import src.model as model
import datasketch
import re

from typing import Callable


class Service:
    def __init__(
        self,
        redis_host: str,
        redis_port: int,
    ):
        self.shingle_size = 4
        self.minhash_permutations = 128
        self.similarity_threshold = 0.8

        self._index = datasketch.MinHashLSH(
            threshold=self.similarity_threshold,
            num_perm=self.minhash_permutations,
            storage_config={
                "type": "redis",
                "basename": b"shingleindex",
                "redis": {"host": redis_host, "port": redis_port},
            },
        )

    def search_similar(self, doc: model.Document) -> list[model.SimilarDocument]:
        text = doc.text.decode()

        mh = self._minhash(text)

        results = self._index.query(mh)

        return [
            model.SimilarDocument(
                **{
                    "id": i,
                    "rate": self.similarity_threshold,
                    "highlights": [],
                    "similarImages": [],
                }
            )
            for i in results
        ]

    def index_content(self, doc: model.Document):
        text = doc.text.decode()

        mh = self._minhash(text)

        self._index.insert(doc.id, mh, check_duplication=False)

    def _minhash(self, text: str) -> datasketch.MinHash:
        text = self._normalize(text)

        shingles = self._shingle(
            text, shingle_size=self.shingle_size, hashfunc=self._hash_func
        )

        mh = datasketch.MinHash(self.minhash_permutations)

        for shingle in shingles:
            mh.update(shingle.encode("utf8"))

        return mh

    def _normalize(self, text: str) -> str:
        text = re.sub(r"[^\w\s.]", "", text)
        text = re.sub(r"\s+", " ", text)
        text = text.lower()

        return text

    def _shingle(
        self, text: str, shingle_size=4, hashfunc: Callable = lambda x: x
    ) -> set[str]:
        words = list(text.split())

        shings = [
            words[i : i + shingle_size] for i in range(len(words) - shingle_size + 1)
        ]

        shings = [" ".join(i) for i in shings]

        shings = [hashfunc(i) for i in shings]

        return set(shings)

    @staticmethod
    def _hash_func(s: str, salt: int = 1) -> int:
        return hash(s + str(salt))
