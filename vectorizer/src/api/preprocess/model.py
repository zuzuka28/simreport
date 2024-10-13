from src.domain.model.image import ImageHashes
import pydantic


class ImageHashesResponse(pydantic.BaseModel):
    phash: str
    phash_vector: list[float]

    dhash: str
    dhash_vector: list[float]

    ahash: str
    ahash_vector: list[float]

    whash: str
    whash_vector: list[float]


class VectorizeTextRequest(pydantic.BaseModel):
    text: str


class VectorizeTextReponse(pydantic.BaseModel):
    vector: list[float]


class VectorizeImageReponse(pydantic.BaseModel):
    vector: list[float]


def map_image_hashes_to_response(
    source: ImageHashes | None,
) -> ImageHashesResponse | None:
    if source is None:
        return None

    return ImageHashesResponse(**source.model_dump())
