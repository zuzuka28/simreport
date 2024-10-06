from dependency_injector.wiring import Provide, inject
from fastapi import APIRouter, Depends, UploadFile
from src.container import Container

from src.api.preprocess.model import DocumentResponse, map_document_to_response
from src.service.preprocessor_service import (
    DocumentPreprocessorService,
)

preprocess_router = APIRouter()


@preprocess_router.post(
    "/preprocess/doc",
    response_model=DocumentResponse,
)
@inject
async def process_item(
    file: UploadFile,
    service: DocumentPreprocessorService = Depends(
        Provide[Container.preprocessor_service]
    ),
):
    contents = await file.read()
    result = await service.preprocess(contents)

    return map_document_to_response(result)
