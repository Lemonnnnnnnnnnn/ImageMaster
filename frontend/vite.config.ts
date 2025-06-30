import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'
// @ts-expect-error
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [tailwindcss(),svelte(),]
})
