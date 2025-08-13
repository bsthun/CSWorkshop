<script lang="ts">
	import { cn } from '$/lib/shadcn/utils.js'

	let className: string | undefined = undefined
	export { className as class }
	export let value: string = ''
	export let placeholder: string = ''
	export let disabled: boolean = false
	export let rows: number = 3
</script>

<textarea
	class={cn(
		'border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50',
		className
	)}
	bind:value
	{placeholder}
	{disabled}
	{rows}
	{...$$restProps}
	on:blur
	on:change
	on:click
	on:focus
	on:focusin
	on:focusout
	on:keydown
	on:keypress
	on:keyup
	on:mousedown
	on:mouseenter
	on:mouseleave
	on:mouseover
	on:mouseup
	on:paste
	on:input
></textarea>
