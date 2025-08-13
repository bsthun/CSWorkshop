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
	import type { PayloadSemester } from '$/util/backend/backend.ts'

	export let open = false

	const dispatch = createEventDispatcher<{
		created: void
		close: void
	}>()

	let creating = false
	let semesters: PayloadSemester[] = []
	let loadingSemesters = false
	let name = ''
	let code = ''
	let selectedSemesterId: string[]

	const loadSemesters = () => {
		loadingSemesters = true
		backend.admin
			.semesterList({})
			.then((response) => {
				semesters = response.data.semesters!
				selectedSemesterId = [semesters[0]?.id + '']
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loadingSemesters = false
			})
	}

	const createClass = () => {
		if (!name.trim()) {
			toast.error('Please enter a class name')
			return
		}
		if (!code.trim()) {
			toast.error('Please enter a class code')
			return
		}
		if (!selectedSemesterId) {
			toast.error('Please select a semester')
			return
		}

		creating = true
		backend.admin
			.classCreate({
				name: name.trim(),
				code: code.trim(),
				semesterId: selectedSemesterId[0] as any,
			})
			.then(() => {
				toast.success('Class created successfully')
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
		code = ''
		selectedSemesterId = []
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
			loadSemesters()
		}
	}

	$: isFormValid = () => {
		return name.trim() !== '' && code.trim() !== '' && selectedSemesterId.length > 0
	}

	onMount(() => {
		loadSemesters()
	})
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent class="sm:max-w-[425px]">
		<DialogHeader>
			<DialogTitle>Create New Class</DialogTitle>
			<DialogDescription>Create a new class and assign it to a semester.</DialogDescription>
		</DialogHeader>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="name">Class Name</Label>
				<Input bind:value={name} disabled={creating} id="name" placeholder="Enter class name" />
			</div>
			<div class="grid gap-2">
				<Label for="code">Class Code</Label>
				<Input bind:value={code} disabled={creating} id="code" placeholder="Enter class code (e.g., CS101)" />
			</div>
			<div class="grid gap-2">
				<Label for="semester">Semester</Label>
				{#if loadingSemesters}
					<div class="flex items-center justify-center py-2">
						<Loader2Icon class="h-4 w-4 animate-spin" />
					</div>
				{:else}
					<Select bind:value={selectedSemesterId} disabled={creating}>
						<SelectTrigger>
							{semesters.find((s) => s.id + '' === selectedSemesterId[0])?.name || 'Select Semester'}
						</SelectTrigger>
						<SelectContent>
							{#each semesters as semester}
								<SelectItem value={semester.id + ''}>{semester.name}</SelectItem>
							{/each}
						</SelectContent>
					</Select>
				{/if}
			</div>
		</div>
		<DialogFooter>
			<Button disabled={creating} onclick={() => (open = false)} variant="outline">Cancel</Button>
			<Button disabled={!isFormValid()} onclick={createClass}>
				{#if creating}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Create
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
