<script lang="ts">
	import { navigate, useLocation } from 'svelte-navigator'
	import { getContext, onMount } from 'svelte'
	import type { Writable } from 'svelte/store'
	import type { Setup } from '$/util/type/setup.ts'
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { CheckCircle, Loader2Icon, XCircle } from 'lucide-svelte'
	import { toast } from 'svelte-sonner'
	import { backend } from '$/util/backend.ts'
	import Container from '$/component/layout/Container.svelte'

	const location = useLocation()
	const params = new URLSearchParams($location.search)
	const code = params.get('code')
	const setup = getContext<Writable<Setup>>('setup')

	let loading = true
	let success = false
	let error = ''

	const processCallback = () => {
		if (!code) {
			error = 'Missing required parameters'
			loading = false
			setTimeout(() => {
				navigate('/entry/login')
			}, 3000)
			return
		}

		backend.public
			.loginCallback({
				code: code,
			})
			.then((res) => {
				if (res.success) {
					success = true
					toast.success('Successfully logged in')
					$setup.reload().then(() => {
						setTimeout(() => {
							navigate('/')
						}, 1000)
					})
				} else {
					error = res.message || 'Authentication failed'
					setTimeout(() => {
						navigate('/entry/login')
					}, 3000)
				}
			})
			.catch((err) => {
				if (err.response?.data) {
					error = err.response.data.message || 'Authentication failed'
					toast.error(err.response.data.message, {
						description: err.response.data.error || undefined,
					})
				} else {
					error = err.message
					toast.error(err.message)
				}

				setTimeout(() => {
					navigate('/entry/login')
				}, 3000)
			})
			.finally(() => {
				loading = false
			})
	}

	onMount(() => {
		processCallback()
	})
</script>

<Container class="bg-gray-50-dark flex min-h-dvh items-center justify-center">
	<Card class="border-gray-200-dark dark:bg-gray-100-dark mx-4 w-full max-w-md rounded-lg bg-white shadow-md">
		<CardHeader>
			<CardTitle>
				{#if loading}
					Processing Login
				{:else if success}
					Login Successful
				{:else}
					Login Failed
				{/if}
			</CardTitle>
			<CardDescription>
				{#if loading}
					Please wait while we verify your credentials
				{:else if success}
					You will be redirected shortly
				{:else}
					{error}. Redirecting to login page...
				{/if}
			</CardDescription>
		</CardHeader>
		<CardContent class="text-center">
			<div class="mb-4 flex justify-center">
				{#if loading}
					<Loader2Icon class="text-primary h-12 w-12 animate-spin" />
				{:else if success}
					<CheckCircle class="h-12 w-12 text-green-600" />
				{:else}
					<XCircle class="h-12 w-12 text-red-600" />
				{/if}
			</div>
		</CardContent>
	</Card>
</Container>
