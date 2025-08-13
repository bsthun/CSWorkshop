<script lang="ts">
	let className = ''
	export { className as class }
</script>

<div class="m-auto max-w-screen-xl px-6 pt-24 pb-8 max-lg:px-4 max-lg:pt-16 max-md:px-2 {className}">
	<slot />
</div>
