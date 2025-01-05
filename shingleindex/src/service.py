import model


class Service:
    def __init__(
        self,
        redis_url: str,
    ):
        pass

    def search_similar(self, doc: model.Document) -> list[model.SimilarDocument]:
        print(f"got document on search_similar: {doc}")
        return []

    def index_content(self, doc: model.Document):
        print(f"got document on index_content: {doc}")
        pass
