<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '$/lib/shadcn/components/ui/dialog'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Loader2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { toast } from 'svelte-sonner'
	import type { PayloadExam } from '$/util/backend/backend.ts'

	export let open = false
	export let examData: PayloadExam

	const dispatch = createEventDispatcher<{
		updated: void
		close: void
	}>()

	let updating = false
	let name = ''
	let openedAt = ''
	let closedAt = ''

	$: if (open && examData) {
		name = examData.name
		openedAt = examData.openedAt ? formatDateTimeForInput(examData.openedAt) : ''
		closedAt = examData.closedAt ? formatDateTimeForInput(examData.closedAt) : ''
	}

	const formatDateTimeForInput = (datetime: string) => {
		if (!datetime) return ''
		const date = new Date(datetime)
		const year = date.getFullYear()
		const month = String(date.getMonth() + 1).padStart(2, '0')
		const day = String(date.getDate()).padStart(2, '0')
		const hours = String(date.getHours()).padStart(2, '0')
		const minutes = String(date.getMinutes()).padStart(2, '0')
		return `${year}-${month}-${day}T${hours}:${minutes}`
	}

	const updateExam = () => {
		if (!name.trim()) {
			toast.error('Please enter an exam name')
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

		updating = true
		backend.admin
			.examEdit({
				examId: examData.id,
				name: name.trim(),
				openedAt: openDate.toISOString(),
				closedAt: closeDate.toISOString()
			})
			.then((response) => {
				if (response.success) {
					toast.success('Exam updated successfully')
					open = false
					dispatch('updated')
				} else {
					toast.error(response.message || 'Failed to update exam')
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				updating = false
			})
	}

	const resetDialog = () => {
		dispatch('close')
	}

	const handleOpenChange = (isOpen: boolean) => {
		open = isOpen
		if (!isOpen) {
			resetDialog()
		}
	}

	$: isFormValid = () => {
		return name.trim() && openedAt && closedAt && !updating
	}
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent class="sm:max-w-[500px]">
		<DialogHeader>
			<DialogTitle>Edit Exam</DialogTitle>
			<DialogDescription>Update exam details and schedule.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Exam Name</Label>
				<Input bind:value={name} disabled={updating} id="name" placeholder="Enter exam name" />
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="grid gap-2">
					<Label for="openedAt">Opens At</Label>
					<Input bind:value={openedAt} disabled={updating} id="openedAt" type="datetime-local" />
				</div>
				<div class="grid gap-2">
					<Label for="closedAt">Closes At</Label>
					<Input bind:value={closedAt} disabled={updating} id="closedAt" type="datetime-local" />
				</div>
			</div>
		</div>
		<DialogFooter>
			<Button disabled={updating} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={!isFormValid()} onclick={updateExam}>
				{#if updating}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Update Exam
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>