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

	const dispatch = createEventDispatcher<{
		created: void
		close: void
	}>()

	let creating = false
	let name = ''

	const createSemester = () => {
		if (!name.trim()) {
			toast.error('Please enter a semester name')
			return
		}

		creating = true
		backend.admin
			.semesterCreate({
				name: name.trim(),
			})
			.then(() => {
				toast.success('Semester created successfully')
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
	}

	const resetDialog = () => {
		resetForm()
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
			<DialogTitle>Create New Semester</DialogTitle>
			<DialogDescription>Create a new semester to organize your classes.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Semester Name</Label>
				<Input
					bind:value={name}
					disabled={creating}
					id="name"
					placeholder="Enter semester name (e.g., Fall 2024)"
				/>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={creating} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={creating || !name.trim()} onclick={createSemester}>
				{#if creating}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Create
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>