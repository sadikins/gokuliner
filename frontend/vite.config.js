import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path';


// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss()
  ],
  resolve: {
    alias: {
      // Example: Alias 'src' to the project's 'src' directory
      '@': path.resolve(__dirname, './src'),
      // You can define multiple aliases for specific directories
    },
  },
})
