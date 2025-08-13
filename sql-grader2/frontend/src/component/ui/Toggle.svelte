<script lang="ts">
	export let checked = false
	export let disabled = false
	export let size: 'sm' | 'md' | 'lg' = 'md'
	export let ariaLabel = ''
	export let onclick = (_: Event) => {
		if (!disabled) {
			checked = !checked
		}
	}

	let className = ''
	export { className as class }

	const sizeConfig = {
		sm: {
			button: 'h-5 w-9',
			thumb: 'h-3 w-3',
			translateOn: 'translate-x-5',
			translateOff: 'translate-x-1',
		},
		md: {
			button: 'h-6 w-11',
			thumb: 'h-4 w-4',
			translateOn: 'translate-x-6',
			translateOff: 'translate-x-1',
		},
		lg: {
			button: 'h-7 w-12',
			thumb: 'h-5 w-5',
			translateOn: 'translate-x-7',
			translateOff: 'translate-x-1',
		},
	}

	$: config = sizeConfig[size]
</script>

<button
	aria-label={ariaLabel}
	aria-pressed={checked}
	class="relative inline-flex {config.button} items-center rounded-full transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 {checked
		? 'bg-blue-600'
		: 'bg-gray-200'} {disabled ? 'cursor-not-allowed opacity-50' : ''} {className}"
	{disabled}
	on:click={onclick}
>
	<span
		class="inline-block {config.thumb} transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {checked
			? config.translateOn
			: config.translateOff}"
	></span>
</button>
