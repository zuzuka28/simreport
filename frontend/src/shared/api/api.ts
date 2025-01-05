import { paths } from "./api-types";
import axios, { AxiosInstance } from "axios";

export type UploadDocumentResponse =
    paths["/document/upload"]["post"]["responses"]["200"]["content"]["application/json"];
export type SearchDocumentsRequest =
    paths["/document/search"]["post"]["requestBody"]["content"]["application/json"];
export type SearchDocumentsResponse =
    paths["/document/search"]["post"]["responses"]["200"]["content"]["application/json"];
export type SimilarDocumentsResponse =
    paths["/analyze/{document_id}/similar"]["get"]["responses"]["200"]["content"]["application/json"];

export class DocumentApi {
    private axiosInstance: AxiosInstance;

    constructor(axiosInstance: AxiosInstance) {
        this.axiosInstance = axiosInstance;
    }

    async uploadDocument(formData: FormData): Promise<UploadDocumentResponse> {
        const response = await this.axiosInstance.post<UploadDocumentResponse>(
            "/document/upload",
            formData,
            {
                headers: { "Content-Type": "multipart/form-data" },
            },
        );
        return response.data;
    }

    async downloadDocument(documentId: string): Promise<Blob> {
        const response = await this.axiosInstance.get<Blob>(
            `/document/${documentId}/download`,
            {
                responseType: "blob",
            },
        );
        return response.data;
    }

    async searchDocuments(
        request: SearchDocumentsRequest,
    ): Promise<SearchDocumentsResponse> {
        const response = await this.axiosInstance.post<SearchDocumentsResponse>(
            "/document/search",
            request,
        );
        return response.data;
    }
}

export class AnalyzeApi {
    private axiosInstance: AxiosInstance;

    constructor(axiosInstance: AxiosInstance) {
        this.axiosInstance = axiosInstance;
    }

    async findSimilarDocuments(id: string): Promise<SimilarDocumentsResponse> {
        const response = await this.axiosInstance.get<SimilarDocumentsResponse>(
            `/analyze/${id}/similar`,
        );
        return response.data;
    }
}

export class ApiClient {
    public document: DocumentApi;
    public analyze: AnalyzeApi;

    constructor(baseURL: string, token?: string) {
        const axiosInstance = axios.create({
            baseURL,
            headers: token ? { Authorization: `Bearer ${token}` } : {},
        });

        this.document = new DocumentApi(axiosInstance);
        this.analyze = new AnalyzeApi(axiosInstance);
    }
}
