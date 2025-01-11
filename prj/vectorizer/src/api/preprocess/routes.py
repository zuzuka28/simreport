from dependency_injector.wiring import Provide, inject
from fastapi import APIRouter, Depends, File, UploadFile
from src.container import Container

from src.api.preprocess.model import (
    ImageHashesResponse,
    VectorizeImageReponse,
    VectorizeTextReponse,
    VectorizeTextRequest,
    map_image_hashes_to_response,
)
from src.domain.service.vectorizer_service import DocumentVectorizerService

preprocess_router = APIRouter()


@preprocess_router.post(
    "/vectorize/text",
    response_model=VectorizeTextReponse,
)
@inject
async def vectorize_text(
    request: VectorizeTextRequest,
    service: DocumentVectorizerService = Depends(
        Provide[Container.vectorizer_service]),
):
    contents = request.text
    result = await service.vectorize_text(contents)

    return VectorizeTextReponse(**{"vector": result.tolist()})


@preprocess_router.post(
    "/vectorize/image",
    response_model=VectorizeImageReponse,
)
@inject
async def vectorize_image(
    file: UploadFile = File(...),
    service: DocumentVectorizerService = Depends(
        Provide[Container.vectorizer_service]),
):
    contents = await file.read()

    result = await service.vectorize_image(contents)

    return VectorizeImageReponse(**{"vector": result.tolist()})


@preprocess_router.post(
    "/hash/image",
    response_model=ImageHashesResponse,
)
@inject
async def hash_image(
    file: UploadFile = File(...),
    service: DocumentVectorizerService = Depends(
        Provide[Container.vectorizer_service]),
):
    contents = await file.read()

    result = await service.retrieve_image_hashes(contents)

    return map_image_hashes_to_response(result)
