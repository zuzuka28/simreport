import pydantic


class Document(pydantic.BaseModel):
    text: bytes


class SimilarDocument(pydantic.BaseModel):
    id: str
    rate: float
    highlights: list[str]
    similarImages: list[str]
