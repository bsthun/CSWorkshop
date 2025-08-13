<script lang="ts">
	import { onMount } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { SectionIcon, GraduationCapIcon, PlusIcon, Loader2Icon, BookOpenIcon, HashIcon } from 'lucide-svelte'
	import { Link } from 'svelte-navigator'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import Container from '$/component/layout/Container.svelte'
	import FloatingAction from '$/component/ui/FloatingAction.svelte'
	import SearchPaginate from '$/component/ui/floating/SearchPaginate.svelte'
	import CreateCollectionDialog from '../dialog/CreateCollectionDialog.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadCollection } from '$/util/backend/backend.ts'

	let activeTab: 'collection' | 'class' = 'collection'
	let collections: PayloadCollection[] = []
	let loading = true
	let showCreateDialog = false
	let creating = false

	let nameFilter = ''
	let currentPage = 1
	let totalCount = 0
	const itemsPerPage = 24

	const loadCollections = (name: string, page: number) => {
		loading = true
		const offset = (page - 1) * itemsPerPage

		backend.admin
			.collectionList({
				limit: itemsPerPage,
				offset,
				name,
			})
			.then((response) => {
				collections = response.data.collections!
				totalCount = response.data.count!
				currentPage = page
				nameFilter = name
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleCollectionCreated = () => {
		loadCollections('', 1)
	}

	const handleSearch = (event: CustomEvent<{ query: string; page: number }>) => {
		loadCollections(event.detail.query, event.detail.page)
	}

	const handlePaginate = (event: CustomEvent<{ page: number }>) => {
		loadCollections(nameFilter, event.detail.page)
	}

	const mount = () => {
		loadCollections('', 1)
	}

	onMount(mount)
</script>

<Container>
	<div class="mb-6 flex items-center justify-between">
		<PageTitle description="Manage collections and classes" title="Admin Panel" />
		{#if activeTab === 'collection'}
			<Button class="gap-2" onclick={() => (showCreateDialog = true)}>
				<PlusIcon class="h-4 w-4" />
				Add Collection
			</Button>
		{/if}
	</div>

	<!-- Tab Navigation -->
	<div class="mb-6 flex space-x-1 rounded-lg bg-gray-100 p-1">
		<button
			class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === 'collection'
				? 'bg-white text-gray-900 shadow-sm'
				: 'text-gray-500 hover:text-gray-900'}"
			onclick={() => (activeTab = 'collection')}
		>
			<div class="flex items-center gap-2">
				<SectionIcon size={16} />
				Collections
			</div>
		</button>
		<button
			class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === 'class'
				? 'bg-white text-gray-900 shadow-sm'
				: 'text-gray-500 hover:text-gray-900'}"
			onclick={() => (activeTab = 'class')}
		>
			<div class="flex items-center gap-2">
				<GraduationCapIcon size={16} />
				Classes
			</div>
		</button>
	</div>

	{#if activeTab === 'collection'}
		{#if loading}
			<div class="flex min-h-[400px] items-center justify-center">
				<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
			</div>
		{:else if collections.length === 0}
			<div class="flex min-h-[400px] flex-col items-center justify-center">
				<SectionIcon class="mb-4 h-16 w-16 text-gray-400" />
				<h3 class="mb-2 text-lg font-semibold">No collections yet</h3>
				<p class="text-muted-foreground mb-4">Create your first collection to get started</p>
				<Button class="gap-2" onclick={() => (showCreateDialog = true)}>
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
	{:else}
		{#await import('../class/index.svelte') then ClassList}
			<ClassList.default />
		{:catch}
			<div class="flex min-h-[400px] flex-col items-center justify-center">
				<GraduationCapIcon class="mb-4 h-16 w-16 text-gray-400" />
				<h3 class="mb-2 text-lg font-semibold">Classes</h3>
				<p class="text-muted-foreground">Loading class management...</p>
			</div>
		{/await}
	{/if}
</Container>

{#if activeTab === 'collection'}
	<FloatingAction>
		<SearchPaginate
			{itemsPerPage}
			placeholder="Search collections..."
			bind:query={nameFilter}
			bind:currentPage
			bind:totalItems={totalCount}
			on:search={handleSearch}
			on:paginate={handlePaginate}
		/>
	</FloatingAction>
{/if}

<CreateCollectionDialog bind:creating bind:open={showCreateDialog} on:created={handleCollectionCreated} />
