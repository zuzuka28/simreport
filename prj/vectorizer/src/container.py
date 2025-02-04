import sys

from dependency_injector.containers import DeclarativeContainer, WiringConfiguration
from dependency_injector.providers import Singleton

from src.service.vectorizer_service import VectorizerService


class Container(DeclarativeContainer):
    wiring_config = WiringConfiguration(modules=[".api.vectorize.routes"])

    vectorizer_service = Singleton(VectorizerService)
