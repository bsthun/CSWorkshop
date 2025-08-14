<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { SectionIcon, PlusIcon, Loader2Icon, BookOpenIcon, HashIcon } from 'lucide-svelte'
	import { Link } from 'svelte-navigator'
	import type { PayloadCollection } from '$/util/backend/backend.ts'

	export let collections: PayloadCollection[] = []
	export let loading = false

	const dispatch = createEventDispatcher<{
		createCollection: void
	}>()
</script>

{#if loading}
	<div class="flex min-h-[400px] items-center justify-center">
		<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
	</div>
{:else if collections.length === 0}
	<div class="flex min-h-[400px] flex-col items-center justify-center">
		<SectionIcon class="mb-4 h-16 w-16 text-gray-400" />
		<h3 class="mb-2 text-lg font-semibold">No collections yet</h3>
		<p class="text-muted-foreground mb-4">Create your first collection to get started</p>
		<Button class="gap-2" onclick={() => dispatch('createCollection')}>
			<PlusIcon class="h-4 w-4" />
			Create Collection
		</Button>
	</div>
{:else}
	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each collections as collection}
			<Link to="/admin/collection/{collection.id}">
				<Card class="transition-shadow hover:shadow-lg">
					<CardHeader class="gap-2">
						<CardTitle class="flex items-center gap-1">
							<SectionIcon size={12} class="text-muted-foreground" />
							<span class="text-muted-foreground text-xs font-medium tracking-wide uppercase">
								COLLECTION
							</span>
						</CardTitle>
						<h3 class="text-lg font-semibold">{collection.name}</h3>
					</CardHeader>
					<CardContent>
						<div class="space-y-2">
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<BookOpenIcon class="h-4 w-4" />
								<span>{collection.questionCount} questions</span>
							</div>
							<div class="text-muted-foreground flex items-center gap-2 text-sm">
								<HashIcon class="h-4 w-4" />
								<span>{collection.id}</span>
							</div>
						</div>
					</CardContent>
				</Card>
			</Link>
		{/each}
	</div>
{/if}
