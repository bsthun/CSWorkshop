<script lang="ts">
	import { Card, CardContent } from '$/lib/shadcn/components/ui/card'
	import Container from '$/component/layout/Container.svelte'
	import { AlertTriangleIcon, EllipsisVerticalIcon } from 'lucide-svelte'
</script>

<Container class="px-6">
	<div class="flex min-h-[80vh] items-center justify-center">
		<Card class="w-full max-w-md bg-white shadow-md">
			<CardContent class="pt-6">
				<div class="flex flex-col items-center text-center">
					<div class="mb-4 rounded-full bg-amber-100 p-3">
						<AlertTriangleIcon class="h-10 w-10 text-amber-500" />
					</div>
					<h1 class="mb-2 text-2xl font-semibold">Webview is not supported</h1>
					<p class="text-muted-foreground mb-6">
						Please open this website in the native browser. You can click on the
						<EllipsisVerticalIcon
							class="border-gray-300-dark inline animate-pulse border p-0.5"
							size="1.2em"
						/>
						button and choose to open in browser.
					</p>
				</div>
			</CardContent>
		</Card>
	</div>
</Container>
