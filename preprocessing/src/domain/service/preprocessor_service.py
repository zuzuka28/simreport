from abc import ABC, abstractmethod

from src.domain.model.document import Document


class DocumentPreprocessorService(ABC):
    @abstractmethod
    async def preprocess(self, source: bytes) -> Document:
        """Preprocess Document from bytes"""
