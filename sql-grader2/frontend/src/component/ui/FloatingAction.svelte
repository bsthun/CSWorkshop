<script lang="ts">
	import { fly } from 'svelte/transition'

	export let visible: boolean = true
</script>

{#if visible}
	<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 transform" transition:fly={{ y: 50, duration: 200 }}>
		<div class="rounded-lg border bg-white p-4 shadow-lg">
			<slot />
		</div>
	</div>
{/if}
