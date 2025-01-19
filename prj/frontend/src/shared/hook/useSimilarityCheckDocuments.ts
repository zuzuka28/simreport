import client from "shared/api";
import { components } from "shared/api/document/api-types";
import { createDataHook } from "./createDataHook";

type DocumentSimilarityMatchResult = components["responses"]["SimilaritySearchResult"]["content"]["application/json"];

type DocumentSimilarityQuery = string;

export const useSimilarityCheckDocuments = createDataHook<DocumentSimilarityQuery, DocumentSimilarityMatchResult>(
    (q: DocumentSimilarityQuery) => { return client.analyze.findSimilarDocuments(q)},
)
