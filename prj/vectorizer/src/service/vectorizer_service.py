import numpy as np
from sentence_transformers import SentenceTransformer

from src.domain.service.vectorizer_service import (
    DocumentVectorizerService as AVectorizerService,
)


class VectorizerService(AVectorizerService):
    _text_encoder = "all-mpnet-base-v2"

    def __init__(self):
        pass

    async def vectorize_text(self, text: str) -> np.ndarray:
        model = SentenceTransformer(self._text_encoder)
        embeddings = model.encode([text])[0]

        return embeddings
