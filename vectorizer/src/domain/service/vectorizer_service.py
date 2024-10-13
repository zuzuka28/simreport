from abc import ABC, abstractmethod

import numpy as np

from src.domain.model.document import Document
from src.domain.model.image import ImageHashes


class DocumentVectorizerService(ABC):
    @abstractmethod
    async def vectorize_document(self, doc: Document) -> Document:
        """vectorize content in document"""

    @abstractmethod
    async def vectorize_text(self, text: str) -> np.ndarray:
        """vectorize text document"""

    @abstractmethod
    async def vectorize_image(self, img: bytes) -> np.ndarray:
        """vectorize image document"""

    @abstractmethod
    async def retrieve_image_hashes(self, img: bytes) -> ImageHashes:
        """retrieve image hashes"""
