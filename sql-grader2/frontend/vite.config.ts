import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
	plugins: [svelte(), tailwindcss()],
	build: {
		rollupOptions: {
			output: {
				manualChunks: (id) => {
					if (id.includes('@sentry')) {
						return 'sentry'
					}
					return null
				},
			},
		},
	},
	resolve: {
		alias: {
			$: path.resolve('./src'),
		},
	},
	css: {
		preprocessorOptions: {
			scss: {
				api: 'modern',
			},
		},
	},
	server: {
		host: '0.0.0.0',
		port: 5173,
		proxy: {
			'/api/': {
				target: 'http://localhost:3000/',
				changeOrigin: true,
			},
			'/file/': {
				target: 'http://localhost:3000/',
				changeOrigin: true,
			},
		},
		allowedHosts: true,
	},
})
