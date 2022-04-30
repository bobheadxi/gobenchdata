import { fileURLToPath, URL } from "url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  base: "./",
  build: {
    minify: false,
  },
  css: {
    preprocessorOptions: {
      scss: {
        // @/ is an alias to src/
        additionalData: `@import "@/styles/variables.scss";`,
      },
    },
  },
});
