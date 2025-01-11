import pydantic


class Document(pydantic.BaseModel):
    id: str
    text: bytes


class SimilarDocument(pydantic.BaseModel):
    id: str
    rate: float
    highlights: list[str]
    similarImages: list[str]
