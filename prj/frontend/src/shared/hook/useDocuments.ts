import client from "shared/api";
import { components } from "shared/api/document/api-types";
import { createDataHook } from "./createDataHook";

type DocumentResult = components["responses"]["SearchResult"]["content"]["application/json"];

type DocumentsQuery = components["schemas"]["SearchRequest"];

export const useDocuments = createDataHook<DocumentsQuery, DocumentResult>(
    (q: DocumentsQuery) => { return client.document.searchDocuments(q)} ,
)