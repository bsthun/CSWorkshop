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
	import type { PayloadClass, PayloadSemesterInfo } from '$/util/backend/backend.ts'

	export let open = false
	export let classData: PayloadClass
	export let semester: PayloadSemesterInfo

	const dispatch = createEventDispatcher<{
		edited: void
		close: void
	}>()

	let editing = false
	let name = ''
	let code = ''

	$: if (classData && open) {
		name = classData.name
		code = classData.code
	}

	const editClass = () => {
		if (!name.trim()) {
			toast.error('Please enter a class name')
			return
		}
		if (!code.trim()) {
			toast.error('Please enter a class code')
			return
		}

		editing = true
		backend.admin
			.classEdit({
				id: classData.id,
				name: name.trim(),
				code: code.trim(),
				registerCode: classData.registerCode,
			})
			.then(() => {
				toast.success('Class updated successfully')
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
		if (classData) {
			name = classData.name
			code = classData.code
		}
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
			<DialogTitle>Edit Class</DialogTitle>
			<DialogDescription>Update class information. Semester: {semester?.name}</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Class Name</Label>
				<Input
					bind:value={name}
					disabled={editing}
					id="name"
					placeholder="Enter class name"
				/>
			</div>
			<div class="grid gap-2">
				<Label for="code">Class Code</Label>
				<Input
					bind:value={code}
					disabled={editing}
					id="code"
					placeholder="Enter class code"
				/>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={editing} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={editing || !name.trim() || !code.trim()} onclick={editClass}>
				{#if editing}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Save Changes
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>