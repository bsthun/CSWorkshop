<script lang="ts">
	export let title: string
	export let description: string | undefined = undefined
</script>

<div>
	<h1 class="text-3xl font-bold tracking-tight">{title}</h1>
	{#if description}
		<p class="text-muted-foreground mt-2">{description}</p>
	{/if}
</div>
