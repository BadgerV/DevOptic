import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@app': path.resolve(__dirname, './src/app'),
      '@domains': path.resolve(__dirname, './src/domains'),
      '@shared': path.resolve(__dirname, './src/shared'),
      '@pages': path.resolve(__dirname, './src/pages'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@styles': path.resolve(__dirname, './src/styles'),
    },
  },
  server: {
    port: 3000,
    open: true,
    // host: '0.0.0.0',
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'domain-endpoints': ['./src/domains/endpoint-monitoring'],
          'domain-rbac': ['./src/domains/rbac'],
          'domain-gitlab': ['./src/domains/gitlab'],
          // 'domain-dashboard': ['./src/domains/dashboard'],/
          'shared-charts': ['recharts', 'd3'],
          'shared-ui': ['antd', '@ant-design/icons'],
          'shared-realtime': ['socket.io-client']
        }
      }
    }
  }
})
