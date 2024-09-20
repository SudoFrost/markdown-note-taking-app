import { defineConfig } from "vite";

export default defineConfig({
  build: {
    manifest: true,
    outDir: 'dist',
    rollupOptions: {
      input: {
        js: 'resources/js/main.js',
        css: 'resources/css/app.css',
      }
    },
  },
})