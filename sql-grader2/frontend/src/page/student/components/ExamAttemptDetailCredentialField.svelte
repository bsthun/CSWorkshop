<script lang="ts">
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { CopyIcon, CheckIcon } from 'lucide-svelte'

	export let name: string
	export let value: string
	export let blur: boolean = false
	export let copy: boolean = false
	export let copiedField: string | null = null
	export let onCopy: (value: string, field: string) => void = () => {}

	let visible = false
	$: fieldId = name.toLowerCase()
	$: isBlurred = blur && !visible
</script>

<div class="group flex items-center justify-between">
	<span class="text-gray-500">{name}:</span>
	<div class="flex items-center gap-2">
		{#if blur}
			<span
				class="font-medium transition-all cursor-pointer"
				class:blur-sm={isBlurred}
				class:hover:blur-none={isBlurred}
				onclick={() => visible = !visible}
			>
				{value}
			</span>
		{:else}
			<span class="font-medium">{value}</span>
		{/if}
		{#if copy}
			<Button
				variant="ghost"
				size="sm"
				class="h-6 w-6 p-0 opacity-0 transition-opacity group-hover:opacity-100"
				onclick={() => onCopy(value, fieldId)}
			>
				{#if copiedField === fieldId}
					<CheckIcon size={12} class="text-green-600" />
				{:else}
					<CopyIcon size={12} />
				{/if}
			</Button>
		{:else}
			<div class="h-6 w-6"></div>
		{/if}
	</div>
</div>