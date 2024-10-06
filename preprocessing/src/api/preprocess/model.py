import pydantic
import base64

from src.domain.model.document import Document, Image
from src.domain.model.image import ImageHashes


class ImageHashesResponse(pydantic.BaseModel):
    phash: str
    phash_vector: list[float]

    dhash: str
    dhash_vector: list[float]

    ahash: str
    ahash_vector: list[float]

    whash: str
    whash_vector: list[float]


class ImageResponse(pydantic.BaseModel):
    fname: str
    source_bytes: str
    sha256: str

    hashes: ImageHashesResponse | None = None
    clip_image_vector: list[float] | None = None


class DocumentResponse(pydantic.BaseModel):
    id: str
    source_bytes: str
    sha256: str

    text_content: str
    sbert_text_vector: list[float] | None = None

    images: list[ImageResponse] = []


def map_image_hashes_to_response(
    source: ImageHashes | None,
) -> ImageHashesResponse | None:
    if source is None:
        return None

    return ImageHashesResponse(**source.model_dump())


def map_image_to_response(source: Image) -> ImageResponse:
    return ImageResponse(
        **{
            "fname": source.fname,
            "source_bytes": base64.encodebytes(source.source_bytes),
            "sha256": source.sha256,
            "hashes": map_image_hashes_to_response(source.hashes),
            "clip_image_vector": source.clip_image_vector.tolist()
            if source.clip_image_vector is not None
            else None,
        }
    )


def map_document_to_response(source: Document) -> DocumentResponse:
    return DocumentResponse(
        **{
            "id": source.id,
            "source_bytes": base64.encodebytes(source.source_bytes),
            "sha256": source.sha256,
            "text_content": source.text_content,
            "images": [map_image_to_response(img) for img in source.images],
            "sbert_text_vector": source.sbert_text_vector,
        }
    )
