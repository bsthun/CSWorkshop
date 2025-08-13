<script lang="ts">
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Card, CardContent } from '$/lib/shadcn/components/ui/card'
	import Container from '$/component/layout/Container.svelte'
	import { ArrowLeftIcon } from 'lucide-svelte'
	import { Link } from 'svelte-navigator'
</script>

<Container class="px-6">
	<div class="flex min-h-[80vh] items-center justify-center">
		<Card class="w-full max-w-md bg-white shadow-md">
			<CardContent class="pt-6">
				<div class="flex flex-col items-center text-center">
					<div class="bg-muted mb-4 rounded-full px-8 py-6">
						<span class="text-6xl font-bold">404</span>
					</div>
					<h1 class="mb-2 text-2xl font-semibold">Page Not Found</h1>
					<p class="text-muted-foreground mb-6">The page you are looking for doesn't exist.</p>
					<Link to="/">
						<Button class="gap-2">
							<ArrowLeftIcon class="h-4 w-4" />
							Back to Home
						</Button>
					</Link>
				</div>
			</CardContent>
		</Card>
	</div>
</Container>
