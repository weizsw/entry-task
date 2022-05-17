import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // server: {
  //   proxy: {
  //     '/': {
  //       target: 'http://localhost:8080',
  //       changeOrigin: true,
  //       secure: false,
  //       ws: true,
  //     }
  //   },
  //   // cors: true,
  //   // '/api': {
  //   //   target: 'http://localhost:8080',
  //   //   changeOrigin: true,
  //   //   secure: false,
  //   //   ws: true,
  //   //   rewrite: (path) => path.replace(/^\/api/, '')
  //   // },
  //   // '/api': 'http://localhost:8080',
  // }
})
