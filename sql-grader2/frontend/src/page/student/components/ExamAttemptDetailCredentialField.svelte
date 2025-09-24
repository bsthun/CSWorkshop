<script lang="ts">
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { CopyIcon, CheckIcon } from 'lucide-svelte'

	export let name: string
	export let value: string
	export let blur: boolean = false
	export let copy: boolean = false

	let preCopied = false
	let copied = false

	const onCopy = (value: string) => {
		navigator.clipboard.writeText(value).then(() => {
			copied = true
			preCopied = true
			setTimeout(() => {
				copied = false
			}, 2000)
			setTimeout(() => {
				preCopied = false
			}, 1500)
		})
	}
</script>

<div class="group flex items-center justify-between">
	<span class="text-gray-500">{name}:</span>
	<div class="flex items-center gap-2">
		{#if blur}
			<span class="cursor-pointer font-medium transition-all" class:blur-sm={blur} class:hover:blur-none={blur}>
				{value}
			</span>
		{:else}
			<span class="font-medium">{value}</span>
		{/if}
		{#if copy}
			<Button
				variant="ghost"
				size="sm"
				class="h-6 w-6 p-0 transition-opacity group-hover:opacity-100 {preCopied ? 'opacity-100' : 'opacity-0'}"
				onclick={() => onCopy(value)}
			>
				{#if copied}
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
