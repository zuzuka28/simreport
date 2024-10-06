from src.domain.model.document import Document
from src.domain.service.parser_service import DocumentParserService
from src.domain.service.preprocessor_service import DocumentPreprocessorService
from src.domain.service.vectorizer_service import DocumentVectorizerService


class PreprocessorService(DocumentPreprocessorService):
    def __init__(
        self,
        parser_service: DocumentParserService,
        vectorizer_service: DocumentVectorizerService,
    ):
        self._parser = parser_service
        self._vectorizer = vectorizer_service

    async def preprocess(self, source: bytes) -> Document:
        doc = await self._parser.parse_document(source)
        enriched = await self._vectorizer.vectorize_document(doc)

        return enriched
