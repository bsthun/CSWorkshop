<script lang="ts">
	import { Toaster } from 'svelte-sonner'
	import { writable } from 'svelte/store'
	import { onMount, setContext } from 'svelte'
	import { LoaderCircleIcon } from 'lucide-svelte'
	import { scrollTop } from '$/util/scroll'
	import type { Setup } from '$/util/type/setup'
	import WebviewNotice from '$/component/screen/WebviewNotice.svelte'
	import { backend, catcher } from '$/util/backend.ts'

	scrollTop()

	let properties = {
		webview:
			navigator.userAgent.includes('wv') ||
			(navigator.userAgent.includes('Mobile/') && !navigator.userAgent.includes('Safari/')),
	}

	let setup = writable<Setup>({
		profile: {} as any,
		initialized: false,
		reload: async () => {
			mount()
		},
	})
	setContext('setup', setup)

	const mount = () => {
		backend.state
			.state()
			.then((res) => {
				if (res.success) {
					setup.update((value) => ({
						...value,
						profile: {
							id: res.data.userId.toString(),
							name: res.data.displayName,
							email: res.data.email,
							avatar: res.data.photoUrl,
						},
						initialized: true,
					}))
				} else {
					setup.update((value) => ({
						...value,
						profile: {},
						initialized: true,
					}))
				}
			})
			.catch(() => {
				// TODO: disable catcher for eliminating toast on project session
				// if (err.message !== 'canceled') catcher(err)
				setup.update((value) => {
					return {
						...value,
						profile: {},
						initialized: true,
					}
				})
			})
	}

	onMount(mount)
</script>

<main>
	{#if !$setup.initialized}
		<div class="flex h-screen items-center justify-center">
			<LoaderCircleIcon class="text-gray-600-dark animate-spin" />
		</div>
	{:else if properties.webview}
		<WebviewNotice />
	{:else}
		<slot />
	{/if}

	<Toaster duration={5000} position="bottom-right" richColors />
</main>
