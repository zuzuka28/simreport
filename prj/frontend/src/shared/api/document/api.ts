import { paths } from "./api-types";
import { AxiosInstance } from "axios";

export type UploadDocumentRequest =
  paths["/document/upload"]["post"]["requestBody"]["content"]["application/json"];
export type UploadDocumentResponse =
  paths["/document/upload"]["post"]["responses"]["200"]["content"]["application/json"];

export type SearchDocumentsRequest =
  paths["/document/search"]["post"]["requestBody"]["content"]["application/json"];
export type SearchDocumentsResponse =
  paths["/document/search"]["post"]["responses"]["200"]["content"]["application/json"];

export type SimilarDocumentsResponse =
  paths["/analyze/{document_id}/similar"]["get"]["responses"]["200"]["content"]["application/json"];

export type FetchAttributeRequest =
  paths["/attribute"]["post"]["requestBody"]["content"]["application/json"];
export type FetchAttributeResponse =
  paths["/attribute"]["post"]["responses"]["200"]["content"]["application/json"];

export class DocumentApi {
  private axiosInstance: AxiosInstance;

  constructor(axiosInstance: AxiosInstance) {
    this.axiosInstance = axiosInstance;
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

  async uploadDocument(
    request: UploadDocumentRequest,
  ): Promise<UploadDocumentResponse> {
    const response = await this.axiosInstance.post<UploadDocumentResponse>(
      "/document/upload",
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

export class AttributeApi {
  private axiosInstance: AxiosInstance;

  constructor(axiosInstance: AxiosInstance) {
    this.axiosInstance = axiosInstance;
  }

  async fetchValues(
    request: FetchAttributeRequest,
  ): Promise<FetchAttributeResponse> {
    const response = await this.axiosInstance.post<FetchAttributeResponse>(
      "/attribute",
      request,
    );

    return response.data;
  }
}
