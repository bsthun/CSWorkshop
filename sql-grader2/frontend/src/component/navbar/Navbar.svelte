<script lang="ts">
	import { Link, navigate } from 'svelte-navigator'
	import { getContext } from 'svelte'
	import type { Writable } from 'svelte/store'
	import type { Setup } from '$/util/type/setup'
	import { onMount } from 'svelte'

	let scrolled = false

	const setup = getContext<Writable<Setup>>('setup')

	onMount(() => {
		const handleScroll = () => {
			scrolled = window.scrollY > 20
		}

		window.addEventListener('scroll', handleScroll)

		return () => {
			window.removeEventListener('scroll', handleScroll)
		}
	})
</script>

<nav
	class="fixed start-0 top-0 z-20 flex h-16 w-full justify-center transition-all duration-300 max-lg:h-12"
	class:scrolled
>
	<div class="flex w-full max-w-screen-2xl items-center justify-between px-8">
		<Link to="/" class="flex items-center gap-[18px]">
			<p class="text-[18px] font-medium">SQL Grader</p>
		</Link>
		<div class="flex items-center gap-2"></div>
	</div>
</nav>

<style lang="postcss">
	@reference '$/style/tailwind.css';

	nav {
		@apply border-b border-gray-200 bg-gray-50 shadow-sm;
	}

	nav::after {
		content: '';
		@apply absolute inset-0 -z-10 backdrop-blur-sm transition-all duration-300;
	}

	nav.scrolled {
		@apply border-gray-200 bg-gray-100 shadow-md;
	}

	nav.scrolled::after {
		@apply backdrop-blur-md;
	}
</style>
