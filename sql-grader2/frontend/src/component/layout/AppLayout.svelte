<script lang="ts">
	import Navbar from '$/component/navbar/Navbar.svelte'
	import { getContext } from 'svelte'
	import type { Writable } from 'svelte/store'
	import type { Setup } from '$/util/type/setup'
	import { useLocation } from 'svelte-navigator'

	const setup = getContext<Writable<Setup>>('setup')
	const location = useLocation()

	$: if (!$setup.profile.id) {
		let redirect = true
		if ($location.pathname.match(/^\/project\/.+/)) {
			redirect = false
		}

		if (redirect) {
			window.location.href = '/entry/login/'
		}
	}
</script>

<div>
	<Navbar />
	<div>
		<slot />
	</div>
</div>
