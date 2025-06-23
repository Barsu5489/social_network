import devtoolsJson from 'vite-plugin-devtools-json';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  base: '/',
  /*css: {
    preprocessorOptions: {
      css: {
        additionalData: `@import './src/app.css';`
      }
    }
  },*/
  plugins: [sveltekit(), devtoolsJson()]
});
