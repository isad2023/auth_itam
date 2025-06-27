import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

const apiUrl = process.env.VITE_API_URL || 'http://localhost:8080';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/auth/api': {
        target: 'http://109.73.202.151:8080',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/auth\/api/, '/auth/api'),
      },
    },
  },
})
