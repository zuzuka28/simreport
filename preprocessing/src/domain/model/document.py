import pydantic

import numpy as np

from src.domain.model.image import Image


class Document(pydantic.BaseModel):
    class Config:
        arbitrary_types_allowed = True

    id: str
    source_bytes: bytes
    sha256: str

    text_content: str
    sbert_text_vector: np.ndarray | None = None

    images: list[Image] = []
