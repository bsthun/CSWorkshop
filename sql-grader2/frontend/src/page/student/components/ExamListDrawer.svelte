<script lang="ts">
	import { navigate } from 'svelte-navigator'
	import { ScrollArea } from '$/lib/shadcn/components/ui/scroll-area'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Loader2Icon, FileTextIcon, CalendarIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend'
	import { toast } from 'svelte-sonner'
	import type { PayloadClassExamListItem, PayloadStudentClassListItem } from '$/util/backend/backend.ts'
	import { PayloadClassExamListItemStatusEnum } from '$/util/backend/backend.ts'
	import ExamCard from './ExamCard.svelte'
	import {
		Drawer,
		DrawerContent,
		DrawerDescription,
		DrawerHeader,
		DrawerTitle,
	} from '$/lib/shadcn/components/ui/drawer'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'

	export let open = false
	export let selectedClass: PayloadStudentClassListItem | null = null

	let exams: PayloadClassExamListItem[] = []
	let loading = false
	let accessCodeDialogOpen = false
	let selectedExam: PayloadClassExamListItem | null = null
	let accessCode = ''
	let attemptLoading = false

	$: if (open && selectedClass) {
		loadExams()
	}

	const loadExams = () => {
		if (!selectedClass) return

		loading = true
		backend.student
			.classExamList({ classId: selectedClass.class.id })
			.then((response) => {
				exams = response.data.exams
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleExamClick = (exam: PayloadClassExamListItem) => {
		if (exam.status === PayloadClassExamListItemStatusEnum.Attempted) {
			// Go directly to exam if already attempted
			open = false
			navigate(`/student/exam/${exam.examAttempt.id}/detail`)
		} else if (exam.status === PayloadClassExamListItemStatusEnum.Opened) {
			selectedExam = exam
			accessCode = ''
			accessCodeDialogOpen = true
		}
	}

	const handleAccessCodeSubmit = () => {
		if (!selectedExam || !accessCode.trim()) return

		attemptLoading = true
		backend.student
			.classExamAttempt({
				examId: selectedExam!.exam.id as number,
				accessCode: accessCode.trim(),
			})
			.then((response) => {
				toast.success('successfully attempted the exam')
				accessCodeDialogOpen = false
				open = false
				navigate(`/student/exam/${response!.data.examAttempt.id}/detail`)
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				attemptLoading = false
			})
	}

	const handleAccessCodeDialogClose = () => {
		if (!attemptLoading) {
			accessCodeDialogOpen = false
			accessCode = ''
			selectedExam = null
		}
	}

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter' && !attemptLoading) {
			handleAccessCodeSubmit()
		}
	}
</script>

<Drawer bind:open direction="right">
	<DrawerContent class="w-full min-w-2xl">
		<DrawerHeader>
			<DrawerTitle>
				{#if selectedClass}
					{selectedClass.class.code} - Exams
				{:else}
					Select a Class
				{/if}
			</DrawerTitle>
			<DrawerDescription>
				{#if selectedClass}
					<div class="flex items-center gap-2 text-sm">
						<CalendarIcon class="h-4 w-4" />
						{selectedClass.semester.name}
					</div>
					<div class="mt-1">
						{selectedClass.class.name}
					</div>
				{:else}
					Choose a class to view available exams
				{/if}
			</DrawerDescription>
		</DrawerHeader>

		<div class="flex-1 overflow-visible px-4">
			{#if loading}
				<div class="flex min-h-[300px] items-center justify-center">
					<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
				</div>
			{:else if exams.length === 0}
				<div class="flex min-h-[300px] flex-col items-center justify-center">
					<FileTextIcon class="mb-4 h-16 w-16 text-gray-400" />
					<h3 class="mb-2 text-lg font-semibold">No exams available</h3>
					<p class="text-muted-foreground text-center text-sm">There are no exams for this class yet.</p>
				</div>
			{:else}
				<ScrollArea class="h-[calc(100vh-200px)]">
					<div class="flex flex-col space-y-4 pb-4">
						{#each exams as exam}
							<button onclick={() => handleExamClick(exam)}>
								<ExamCard {exam} />
							</button>
						{/each}
					</div>
				</ScrollArea>
			{/if}
		</div>
	</DrawerContent>
</Drawer>

<Dialog bind:open={accessCodeDialogOpen} onOpenChange={handleAccessCodeDialogClose}>
	<DialogContent class="sm:max-w-[425px]">
		<DialogHeader>
			<DialogTitle>Enter Exam Access Code</DialogTitle>
			<DialogDescription>
				{#if selectedExam}
					Enter the access code for "{selectedExam.exam.name}" to start the exam.
				{:else}
					Enter the access code provided by your instructor.
				{/if}
			</DialogDescription>
		</DialogHeader>

		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="access-code">Access Code</Label>
				<Input
					id="access-code"
					placeholder="Enter access code"
					bind:value={accessCode}
					disabled={attemptLoading}
					onkeydown={handleKeyDown}
				/>
			</div>
		</div>

		<DialogFooter>
			<Button variant="outline" onclick={handleAccessCodeDialogClose} disabled={attemptLoading}>Cancel</Button>
			<Button onclick={handleAccessCodeSubmit} disabled={attemptLoading || !accessCode.trim()}>
				{#if attemptLoading}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
					Starting Exam...
				{:else}
					Start Exam
				{/if}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
