<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$/lib/shadcn/components/ui/select'
	import { Loader2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { toast } from 'svelte-sonner'
	import type { PayloadCollection } from '$/util/backend/backend.ts'

	export let open = false
	export let classId: number

	const dispatch = createEventDispatcher<{
		created: void
		close: void
	}>()

	let creating = false
	let collections: PayloadCollection[] = []
	let loadingCollections = false
	let name = ''
	let selectedCollectionId = ''
	let openedAt = ''
	let closedAt = ''

	const loadCollections = () => {
		loadingCollections = true
		backend.admin
			.collectionList({
				limit: 100,
				offset: 0,
				name: '',
			})
			.then((response) => {
				collections = response.data.collections!
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loadingCollections = false
			})
	}

	const createExam = () => {
		if (!name.trim()) {
			toast.error('Please enter an exam name')
			return
		}
		if (!selectedCollectionId) {
			toast.error('Please select a collection')
			return
		}
		if (!openedAt) {
			toast.error('Please select an open date and time')
			return
		}
		if (!closedAt) {
			toast.error('Please select a close date and time')
			return
		}

		const openDate = new Date(openedAt)
		const closeDate = new Date(closedAt)

		if (closeDate <= openDate) {
			toast.error('Close time must be after open time')
			return
		}

		creating = true
		backend.admin
			.examCreate({
				classId: classId,
				name: name.trim(),
				collectionId: parseInt(selectedCollectionId),
				openedAt: openDate.toISOString(),
				closedAt: closeDate.toISOString(),
			})
			.then(() => {
				toast.success('Exam created successfully')
				resetForm()
				open = false
				dispatch('created')
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				creating = false
			})
	}

	const resetForm = () => {
		name = ''
		selectedCollectionId = ''
		openedAt = ''
		closedAt = ''
	}

	const resetDialog = () => {
		resetForm()
		dispatch('close')
	}

	const handleOpenChange = (isOpen: boolean) => {
		open = isOpen
		if (!isOpen) {
			resetDialog()
		} else {
			loadCollections()
		}
	}

	const isFormValid = () => {
		return name.trim() && selectedCollectionId && openedAt && closedAt && !creating
	}

	const formatDateTimeForInput = (date: Date) => {
		const year = date.getFullYear()
		const month = String(date.getMonth() + 1).padStart(2, '0')
		const day = String(date.getDate()).padStart(2, '0')
		const hours = String(date.getHours()).padStart(2, '0')
		const minutes = String(date.getMinutes()).padStart(2, '0')
		return `${year}-${month}-${day}T${hours}:${minutes}`
	}

	// Set default dates when dialog opens
	$: if (open && !openedAt && !closedAt) {
		const now = new Date()
		const defaultOpen = new Date(now.getTime() + 24 * 60 * 60 * 1000)
		const defaultClose = new Date(defaultOpen.getTime() + 2 * 60 * 60 * 1000)

		openedAt = formatDateTimeForInput(defaultOpen)
		closedAt = formatDateTimeForInput(defaultClose)
	}

	onMount(() => {
		if (open) {
			loadCollections()
		}
	})
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent class="sm:max-w-[500px]">
		<DialogHeader>
			<DialogTitle>Create New Exam</DialogTitle>
			<DialogDescription>Create a new exam and select a question collection.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Exam Name</Label>
				<Input bind:value={name} disabled={creating} id="name" placeholder="Enter exam name" />
			</div>

			<div class="grid gap-2">
				<Label for="collection">Collection</Label>
				{#if loadingCollections}
					<div class="flex items-center justify-center py-2">
						<Loader2Icon class="h-4 w-4 animate-spin" />
					</div>
				{:else}
					<Select bind:value={selectedCollectionId} disabled={creating}>
						<SelectTrigger>Select a collection</SelectTrigger>
						<SelectContent>
							{#each collections as collection}
								<SelectItem value={collection.id.toString()}>
									{collection.name} ({collection.questionCount} questions)
								</SelectItem>
							{/each}
						</SelectContent>
					</Select>
				{/if}
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="grid gap-2">
					<Label for="openedAt">Opens At</Label>
					<Input bind:value={openedAt} disabled={creating} id="openedAt" type="datetime-local" />
				</div>
				<div class="grid gap-2">
					<Label for="closedAt">Closes At</Label>
					<Input bind:value={closedAt} disabled={creating} id="closedAt" type="datetime-local" />
				</div>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={creating} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={!isFormValid()} onclick={createExam}>
				{#if creating}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Create Exam
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
