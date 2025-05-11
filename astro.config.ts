import { defineConfig } from "astro/config";

export default defineConfig({
  compressHTML: false,
  build: {
    format: "file",
  },
});
