import devtoolsJson from 'vite-plugin-devtools-json';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  base: '/',
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        secure: false
      }
    }
  },
  optimizeDeps: {
    exclude: ['svelte-spa-router']
  },
  /*css: {
    preprocessorOptions: {
      css: {
        additionalData: `@import './src/app.css';`
      }
    }
  },*/
  plugins: [sveltekit(), devtoolsJson()]
});
