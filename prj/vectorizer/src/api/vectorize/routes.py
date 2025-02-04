from dependency_injector.wiring import Provide, inject
from fastapi import APIRouter, Depends

from src.api.vectorize.model import (
    VectorizeTextReponse,
    VectorizeTextRequest,
)
from src.container import Container
from src.domain.service.vectorizer_service import DocumentVectorizerService

vectorize_router = APIRouter()


@vectorize_router.post(
    "/vectorize/text",
    response_model=VectorizeTextReponse,
)
@inject
async def vectorize_text(
    request: VectorizeTextRequest,
    service: DocumentVectorizerService = Depends(Provide[Container.vectorizer_service]),
):
    contents = request.text
    result = await service.vectorize_text(contents)

    return VectorizeTextReponse(**{"vector": result.tolist()})
