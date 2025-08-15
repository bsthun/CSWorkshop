<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import type { ComponentType } from 'svelte'

	interface TabItem {
		id: string
		label: string
		icon?: ComponentType
	}

	export let tabs: TabItem[] = []
	export let activeTab: string = ''

	const dispatch = createEventDispatcher<{
		change: string
	}>()

	const handleTabClick = (tabId: string) => {
		if (activeTab !== tabId) {
			dispatch('change', tabId)
		}
	}
</script>

<div class="mb-6 flex space-x-1 rounded-lg bg-gray-100 p-1">
	{#each tabs as tab}
		<button
			class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === tab.id
				? 'bg-white text-gray-900 shadow-sm'
				: 'text-gray-500 hover:text-gray-900'}"
			onclick={() => handleTabClick(tab.id)}
		>
			<span class="flex items-center gap-2">
				{#if tab.icon}
					<svelte:component this={tab.icon} size={16} />
				{/if}
				{tab.label}
			</span>
		</button>
	{/each}
</div>