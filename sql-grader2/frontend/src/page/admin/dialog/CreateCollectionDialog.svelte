<script lang="ts">
	import { createEventDispatcher } from 'svelte'
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
	import { Loader2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { toast } from 'svelte-sonner'

	export let open = false
	export let creating = false

	const dispatch = createEventDispatcher<{
		created: void
		close: void
	}>()

	let newCollectionName = ''

	const createCollection = () => {
		if (!newCollectionName.trim()) {
			toast.error('Please enter a collection name')
			return
		}

		creating = true
		backend.admin
			.collectionCreate({
				name: newCollectionName,
			})
			.then(() => {
				toast.success('Collection created successfully')
				open = false
				newCollectionName = ''
				dispatch('created')
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				creating = false
			})
	}

	const resetDialog = () => {
		newCollectionName = ''
		dispatch('close')
	}

	const handleOpenChange = (isOpen: boolean) => {
		open = isOpen
		if (!isOpen) {
			resetDialog()
		}
	}
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent class="sm:max-w-[425px]">
		<DialogHeader>
			<DialogTitle>Create New Collection</DialogTitle>
			<DialogDescription>Create a new question collection with database schema.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Collection Name</Label>
				<Input
					bind:value={newCollectionName}
					disabled={creating}
					id="name"
					placeholder="Enter collection name"
				/>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={creating} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={creating || !newCollectionName.trim()} onclick={createCollection}>
				{#if creating}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Create
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
