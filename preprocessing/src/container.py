import sys

from dependency_injector.containers import DeclarativeContainer, WiringConfiguration
from dependency_injector.providers import Singleton

from src.service.preprocessor_service import PreprocessorService
from src.service.vectorizer_service import VectorizerService
from src.service.parser_service import ParserService


class Container(DeclarativeContainer):
    wiring_config = WiringConfiguration(modules=[".api.preprocess.routes"])

    parser_service = Singleton(ParserService)
    vectorizer_service = Singleton(VectorizerService)
    preprocessor_service = Singleton(
        PreprocessorService,
        parser_service=parser_service,
        vectorizer_service=vectorizer_service,
    )
