import { ApiClient } from "./api.ts";

const client = new ApiClient(
  import.meta.env.VITE_APP_API_URL || "http://localhost:8080",
);

export default client;
