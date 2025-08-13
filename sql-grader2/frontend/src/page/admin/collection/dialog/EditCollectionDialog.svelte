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
	import type { PayloadCollection } from '$/util/backend/backend.ts'

	export let open = false
	export let saving = false
	export let collection: PayloadCollection | null = null

	const dispatch = createEventDispatcher<{
		saved: void
		close: void
	}>()

	let editedName = ''

	const editCollection = () => {
		if (!editedName.trim()) {
			toast.error('Please enter a collection name')
			return
		}

		if (!collection) {
			toast.error('No collection selected')
			return
		}

		saving = true
		backend.admin
			.collectionEdit({
				id: collection.id,
				name: editedName,
			})
			.then(() => {
				toast.success('Collection updated successfully')
				open = false
				dispatch('saved')
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				saving = false
			})
	}

	const resetDialog = () => {
		editedName = collection?.name || ''
		dispatch('close')
	}

	const handleOpenChange = (isOpen: boolean) => {
		open = isOpen
		if (!isOpen) {
			resetDialog()
		}
	}

	$: if (collection && open) {
		editedName = collection.name || ''
	}
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent class="sm:max-w-[425px]">
		<DialogHeader>
			<DialogTitle>Edit Collection</DialogTitle>
			<DialogDescription>Update the collection name.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Collection Name</Label>
				<Input
					bind:value={editedName}
					disabled={saving}
					id="name"
					placeholder="Enter collection name"
				/>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={saving} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={saving || !editedName.trim()} onclick={editCollection}>
				{#if saving}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Save Changes
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>