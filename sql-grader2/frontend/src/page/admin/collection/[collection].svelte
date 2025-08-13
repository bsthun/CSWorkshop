<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import {
		ArrowLeftIcon,
		BookOpenIcon,
		DatabaseIcon,
		FileIcon,
		InfoIcon,
		Loader2Icon,
		UploadIcon,
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadCollection, PayloadCollectionDetailResponse } from '$/util/backend/backend.ts'
	import TableStructureDialog from './dialog/TableStructureDialog.svelte'
	import SchemaUploadDialog from './dialog/SchemaUploadDialog.svelte'
	import Questions from './components/Questions.svelte'

	export let collection: number

	let collectionData: PayloadCollection | null = null
	let metadata: any = null
	let loading = true
	let showTableDialog = false
	let showSchemaDialog = false

	const loadCollection = () => {
		loading = true
		backend.admin
			.collectionDetail({ collectionId: collection })
			.then((response) => {
				if (response.success && response.data) {
					collectionData = response.data.collection
					metadata = collectionData?.metadata || null
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleSchemaUploaded = () => {
		loadCollection()
	}

	const getTotalRows = () => {
		if (!metadata?.structure) return 0
		return metadata.structure.reduce((total: number, table: any) => total + (table.rowCount || 0), 0)
	}

	onMount(() => {
		loadCollection()
	})
</script>

<Container>
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else if !collectionData}
		<div class="flex min-h-[400px] flex-col items-center justify-center">
			<InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
			<h3 class="mb-2 text-lg font-semibold">Collection not found</h3>
			<p class="text-muted-foreground mb-4">The collection you're looking for doesn't exist</p>
			<Button onclick={() => navigate('/admin/collection')}>
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back to Collections
			</Button>
		</div>
	{:else}
		<div class="mb-6 flex flex-col gap-4">
			<button
				class="text-muted-foreground hover:text-primary flex items-center gap-2 hover:cursor-pointer"
				on:click={() => navigate('/admin/collection')}
			>
				<ArrowLeftIcon size={16} />
				<span class="text-xs font-medium tracking-wide uppercase">COLLECTION</span>
			</button>
			<h1 class="text-3xl font-bold">{collectionData.name}</h1>
		</div>

		<!-- Info Cards -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-3">
			<!-- Question Count Card -->
			<Card>
				<CardContent class="px-6">
					<div class="flex items-center gap-4">
						<div class="rounded-lg bg-blue-100 p-2">
							<BookOpenIcon class="h-6 w-6 text-blue-600" />
						</div>
						<div>
							<p class="text-2xl font-bold">{collectionData.questionCount}</p>
							<p class="text-muted-foreground text-sm">Questions</p>
						</div>
					</div>
				</CardContent>
			</Card>

			<!-- Database Structure Card -->
			<Card>
				<CardContent class="px-6">
					<div class="flex items-center gap-4">
						<div class="rounded-lg bg-green-100 p-2">
							<DatabaseIcon class="h-6 w-6 text-green-600" />
						</div>
						<div>
							{#if metadata?.structure?.length > 0}
								<p class="text-md font-medium">
									{metadata.structure.length} tables, {getTotalRows()} rows
								</p>
								<Button
									variant="ghost"
									size="sm"
									class="mt-1 h-auto p-0 text-blue-600"
									onclick={() => (showTableDialog = true)}
								>
									View Details
								</Button>
							{:else}
								<p class="text-muted-foreground text-sm">No information</p>
							{/if}
						</div>
					</div>
				</CardContent>
			</Card>

			<!-- Schema File Card -->
			<Card>
				<CardContent class="px-6">
					<div class="flex items-center gap-4">
						<div class="rounded-lg bg-purple-100 p-2">
							<FileIcon class="h-6 w-6 text-purple-600" />
						</div>
						<div>
							{#if metadata?.schemaFilename}
								<p class="text-md font-medium">{metadata.schemaFilename}</p>
							{:else}
								<p class="text-muted-foreground text-sm">No schema file</p>
							{/if}
							<Button
								variant="ghost"
								size="sm"
								class="mt-1 h-auto p-0 text-blue-600"
								onclick={() => (showSchemaDialog = true)}
							>
								Upload Schema
							</Button>
						</div>
					</div>
				</CardContent>
			</Card>
		</div>

		<!-- Questions Section -->
		<Questions collectionId={collection} />
	{/if}
</Container>

<TableStructureDialog bind:open={showTableDialog} structure={metadata?.structure || []} />

<SchemaUploadDialog
	bind:open={showSchemaDialog}
	{collection}
	schemaFilename={metadata?.schemaFilename}
	on:uploaded={handleSchemaUploaded}
/>
