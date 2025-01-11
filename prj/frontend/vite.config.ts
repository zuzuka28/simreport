import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      app: "/src/app",
      asset: "/src/asset",
      entity: "/src/entity",
      feature: "/src/feature",
      page: "/src/page",
      shared: "/src/shared",
      widget: "/src/widget",
    },
  },
});
