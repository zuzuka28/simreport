from abc import ABC, abstractmethod

from src.domain.model.document import Document


class DocumentParserService(ABC):
    @abstractmethod
    async def parse_document(self, source: bytes) -> Document:
        """Parse Document from bytes"""
