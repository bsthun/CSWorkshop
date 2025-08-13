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
	import type { PayloadSemester } from '$/util/backend/backend.ts'

	export let open = false
	export let semester: PayloadSemester

	const dispatch = createEventDispatcher<{
		edited: void
		close: void
	}>()

	let editing = false
	let name = ''

	$: if (semester && open) {
		name = semester.name
	}

	const editSemester = () => {
		if (!name.trim()) {
			toast.error('Please enter a semester name')
			return
		}

		editing = true
		backend.admin
			.semesterEdit({
				id: semester.id,
				name: name.trim(),
			})
			.then(() => {
				toast.success('Semester updated successfully')
				open = false
				dispatch('edited')
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				editing = false
			})
	}

	const resetDialog = () => {
		name = semester.name
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
			<DialogTitle>Edit Semester</DialogTitle>
			<DialogDescription>Update the semester name.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Semester Name</Label>
				<Input
					bind:value={name}
					disabled={editing}
					id="name"
					placeholder="Enter semester name"
				/>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={editing} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={editing || !name.trim()} onclick={editSemester}>
				{#if editing}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Save
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>