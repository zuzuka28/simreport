import pydantic

import numpy as np


class ImageHashes(pydantic.BaseModel):
    class Config:
        arbitrary_types_allowed = True

    phash: str
    phash_vector: np.ndarray

    dhash: str
    dhash_vector: np.ndarray

    ahash: str
    ahash_vector: np.ndarray

    whash: str
    whash_vector: np.ndarray


class Image(pydantic.BaseModel):
    class Config:
        arbitrary_types_allowed = True

    fname: str
    source_bytes: bytes
    sha256: str

    hashes: ImageHashes | None = None
    clip_image_vector: np.ndarray | None = None
