<script lang="ts">
	import { Route, Router } from 'svelte-navigator'
	import AppLayout from '$/component/layout/AppLayout.svelte'
	import router from '$/util/router'
	import Wrapper from '$/component/layout/Wrapper.svelte'
	import '$/style/style.scss'
	import '$/style/tailwind.css'
</script>

<Router>
	<Wrapper>
		<Route path="/*">
			<AppLayout>
				{#each router as { path, component }}
					<Route {path} {component} />
				{/each}
			</AppLayout>
		</Route>
		<Route path="/entry/*">
			{#each router as { path, component }}
				{#if path.startsWith('/entry')}
					<Route path={path.replace('/entry', '')} {component} />
				{/if}
			{/each}
		</Route>
	</Wrapper>
</Router>
