import NotFound from '$/component/screen/NotFound.svelte'

// Import all route components
const modules = import.meta.glob('/src/page/**/[a-z[]*.svelte', { eager: true }) as Record<string, { default: any }>

const routes = Object.keys(modules)
	.filter((item) => item.endsWith('.svelte'))
	.map((route) => {
		const path = route
			.replace(/\[\.{3}.+]/, '*')
			.replace(/\[(.+)]/, ':$1')
			.replace(/\/$/, '')
			.replace(/\/src\/page\//g, '/')
			.replace(/.svelte$/, '')
			.replace(/index/, '')

		return { path: path, component: modules[route].default }
	}, {})

export default [...routes, { path: '*', component: NotFound }]
