import { paths } from "./api-types";
import { AxiosInstance } from "axios";

export type UploadFileResponse =
  paths["/upload"]["post"]["responses"]["200"]["content"]["application/json"];

export class AnysaveApi {
  private axiosInstance: AxiosInstance;

  constructor(axiosInstance: AxiosInstance) {
    this.axiosInstance = axiosInstance;
  }

  async uploadFile(formData: FormData): Promise<UploadFileResponse> {
    const response = await this.axiosInstance.post<UploadFileResponse>(
      "/upload",
      formData,
      {
        headers: { "Content-Type": "multipart/form-data" },
      },
    );
    return response.data;
  }

  async downloadDocument(documentId: string): Promise<Blob> {
    const response = await this.axiosInstance.get<Blob>(
      `/${documentId}/download`,
      {
        responseType: "blob",
      },
    );
    return response.data;
  }
}
