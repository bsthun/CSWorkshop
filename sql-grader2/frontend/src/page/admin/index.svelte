<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate, useLocation } from 'svelte-navigator'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { SectionIcon, GraduationCapIcon, PlusIcon } from 'lucide-svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import Container from '$/component/layout/Container.svelte'
	import FloatingAction from '$/component/ui/FloatingAction.svelte'
	import SearchPaginate from '$/component/ui/floating/SearchPaginate.svelte'
	import CreateCollectionDialog from './dialog/CreateCollectionDialog.svelte'
	import CreateClassDialog from './class/dialog/CreateClassDialog.svelte'
	import CreateSemesterDialog from './class/dialog/CreateSemesterDialog.svelte'
	import EditSemesterDialog from './class/dialog/EditSemesterDialog.svelte'
	import CollectionContent from './collection/index.svelte'
	import ClassContent from './class/index.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadCollection, PayloadSemester } from '$/util/backend/backend.ts'

	const location = useLocation()
	let activeTab: 'collection' | 'class'
	
	// Collection state
	let collections: PayloadCollection[] = []
	let collectionsLoading = true
	let showCreateCollectionDialog = false
	let nameFilter = ''
	let currentPage = 1
	let totalCount = 0
	const itemsPerPage = 24

	// Class state
	let semesters: PayloadSemester[] = []
	let classesLoading = true
	let showCreateClassDialog = false
	let showCreateSemesterDialog = false
	let showEditSemesterDialog = false
	let editingSemester: PayloadSemester | null = null

	const loadCollections = (name: string, page: number) => {
		collectionsLoading = true
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
				collectionsLoading = false
			})
	}

	const loadSemesters = () => {
		classesLoading = true
		backend.admin
			.semesterList({})
			.then((response) => {
				semesters = response.data.semesters!
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				classesLoading = false
			})
	}

	const handleCollectionCreated = () => {
		loadCollections('', 1)
	}

	const handleClassCreated = () => {
		loadSemesters()
	}

	const handleSemesterCreated = () => {
		loadSemesters()
	}

	const handleSemesterEdit = (semester: PayloadSemester) => {
		editingSemester = semester
		showEditSemesterDialog = true
	}

	const handleSemesterEdited = () => {
		loadSemesters()
		editingSemester = null
	}

	const handleSearch = (event: CustomEvent<{ query: string; page: number }>) => {
		loadCollections(event.detail.query, event.detail.page)
	}

	const handlePaginate = (event: CustomEvent<{ page: number }>) => {
		loadCollections(nameFilter, event.detail.page)
	}

	const switchTab = (tab: 'collection' | 'class') => {
		activeTab = tab
		const fragment = tab === 'collection' ? '#collection' : '#class'
		navigate(`/admin${fragment}`, { replace: true })
		if (tab === 'collection' && collections.length === 0) {
			loadCollections('', 1)
		} else if (tab === 'class' && semesters.length === 0) {
			loadSemesters()
		}
	}

	const initializeTab = () => {
		const hash = $location.hash.slice(1)
		if (hash === 'class') {
			activeTab = 'class' as const
		} else {
			activeTab = 'collection' as const
		}
	}

	onMount(() => {
		initializeTab()
		loadCollections('', 1)
		loadSemesters()
	})
</script>

<Container>
	<div class="mb-6 flex items-center justify-between">
		<PageTitle description="Manage collections and classes" title="Admin Panel" />
		<div class="flex gap-2">
			{#if activeTab === 'collection'}
				<Button class="gap-2" onclick={() => (showCreateCollectionDialog = true)}>
					<PlusIcon class="h-4 w-4" />
					Add Collection
				</Button>
			{:else}
				<Button variant="outline" class="gap-2" onclick={() => (showCreateSemesterDialog = true)}>
					<PlusIcon class="h-4 w-4" />
					Create Semester
				</Button>
				<Button class="gap-2" onclick={() => (showCreateClassDialog = true)}>
					<PlusIcon class="h-4 w-4" />
					Create Class
				</Button>
			{/if}
		</div>
	</div>

	<!-- Tab Navigation -->
	<div class="mb-6 flex space-x-1 rounded-lg bg-gray-100 p-1">
		<button
			class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === 'collection'
				? 'bg-white text-gray-900 shadow-sm'
				: 'text-gray-500 hover:text-gray-900'}"
			onclick={() => switchTab('collection')}
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
			onclick={() => switchTab('class')}
		>
			<div class="flex items-center gap-2">
				<GraduationCapIcon size={16} />
				Classes
			</div>
		</button>
	</div>

	<!-- Collections Tab Content -->
	{#if activeTab === 'collection'}
		<CollectionContent 
			{collections} 
			loading={collectionsLoading} 
			on:createCollection={() => (showCreateCollectionDialog = true)} 
		/>
	{/if}

	<!-- Classes Tab Content -->
	{#if activeTab === 'class'}
		<ClassContent 
			{semesters} 
			loading={classesLoading} 
			on:createClass={() => (showCreateClassDialog = true)}
			on:createSemester={() => (showCreateSemesterDialog = true)}
			on:editSemester={(e) => handleSemesterEdit(e.detail.semester)}
		/>
	{/if}
</Container>

<!-- Floating Search/Pagination for Collections -->
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

<!-- Dialogs -->
<CreateCollectionDialog bind:open={showCreateCollectionDialog} on:created={handleCollectionCreated} />
<CreateClassDialog bind:open={showCreateClassDialog} {semesters} on:created={handleClassCreated} />
<CreateSemesterDialog bind:open={showCreateSemesterDialog} on:created={handleSemesterCreated} />
{#if editingSemester}
	<EditSemesterDialog
		bind:open={showEditSemesterDialog}
		semester={editingSemester}
		on:edited={handleSemesterEdited}
	/>
{/if}