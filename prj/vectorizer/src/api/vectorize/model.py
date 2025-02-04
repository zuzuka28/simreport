import pydantic


class VectorizeTextRequest(pydantic.BaseModel):
    text: str


class VectorizeTextReponse(pydantic.BaseModel):
    vector: list[float]
