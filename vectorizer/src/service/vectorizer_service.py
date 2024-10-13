import io

import imagehash
import numpy as np
from PIL import Image
from sentence_transformers import SentenceTransformer

from src.domain.model.document import Document
from src.domain.model.image import ImageHashes
from src.domain.service.vectorizer_service import (
    DocumentVectorizerService as AVectorizerService,
)


class VectorizerService(AVectorizerService):
    _text_encoder = "all-mpnet-base-v2"
    _image_encoder = "clip-ViT-B-32"

    def __init__(self):
        pass

    async def vectorize_document(self, doc: Document) -> Document:
        doc.sbert_text_vector = await self.vectorize_text(doc.text_content)

        for img in doc.images:
            img.hashes = await self.retrieve_image_hashes(img.source_bytes)
            img.clip_image_vector = await self.vectorize_image(img.source_bytes)

        return doc

    async def vectorize_text(self, text: str) -> np.ndarray:
        model = SentenceTransformer(self._text_encoder)
        embeddings = model.encode([text])[0]

        return embeddings

    async def vectorize_image(self, img: bytes) -> np.ndarray:
        model = SentenceTransformer(self._image_encoder)
        pilimg = Image.open(io.BytesIO(img))
        embeddings = model.encode([pilimg])[0]

        return embeddings

    async def retrieve_image_hashes(self, img: bytes) -> ImageHashes:
        pilimg = Image.open(io.BytesIO(img))

        phash = imagehash.phash(pilimg)
        dhash = imagehash.dhash(pilimg)
        ahash = imagehash.average_hash(pilimg)
        whash = imagehash.whash(pilimg)

        return ImageHashes(
            **{
                "phash": str(phash),
                "phash_vector": phash.hash.flatten(),
                "dhash": str(dhash),
                "dhash_vector": dhash.hash.flatten(),
                "ahash": str(ahash),
                "ahash_vector": ahash.hash.flatten(),
                "whash": str(whash),
                "whash_vector": whash.hash.flatten(),
            }
        )
