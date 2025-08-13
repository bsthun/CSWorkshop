<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Loader2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import Container from '$/component/layout/Container.svelte'

	let loading = false

	const handleLogin = () => {
		loading = true
		backend.public
			.loginRedirect()
			.then((res) => {
				window.location.href = res.data.redirectUrl
			})
			.catch((err) => {
				catcher(err)
				loading = false
			})
	}
</script>

<Container class="bg-gray-50-dark flex min-h-dvh items-center justify-center">
	<Card class="border-gray-200-dark dark:bg-gray-100-dark mx-4 w-full max-w-md rounded-lg bg-white shadow-md">
		<CardHeader>
			<CardTitle>Login</CardTitle>
			<CardDescription>Sign in to your account</CardDescription>
		</CardHeader>
		<CardContent>
			<Button class="w-full" disabled={loading} onclick={handleLogin} variant="default">
				{#if loading}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Continue with OAuth
			</Button>
		</CardContent>
	</Card>
</Container>
