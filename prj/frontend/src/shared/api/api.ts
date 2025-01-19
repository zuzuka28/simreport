import axios from "axios";
import { DocumentApi, AnalyzeApi, AttributeApi } from "shared/api/document";
import { AnysaveApi } from "shared/api/anysave";

export class ApiClient {
  public document: DocumentApi;
  public analyze: AnalyzeApi;
  public anysave: AnysaveApi;
  public attribute: AttributeApi;

  constructor(baseURL: string, token?: string) {
    const axiosInstance = axios.create({
      baseURL,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });

    this.document = new DocumentApi(axiosInstance);
    this.analyze = new AnalyzeApi(axiosInstance);
    this.attribute = new AttributeApi(axiosInstance);
    this.anysave = new AnysaveApi(axiosInstance);
  }
}
