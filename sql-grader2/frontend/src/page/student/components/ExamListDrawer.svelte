<script lang="ts">
	import { navigate } from 'svelte-navigator'
	import { ScrollArea } from '$/lib/shadcn/components/ui/scroll-area'
	import { Loader2Icon, FileTextIcon, CalendarIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend'
	import type { PayloadClassExamListItem, PayloadStudentClassListItem } from '$/util/backend/backend.ts'
	import ExamCard from './ExamCard.svelte'
	import {
		Drawer,
		DrawerContent,
		DrawerDescription,
		DrawerHeader,
		DrawerTitle,
	} from '$/lib/shadcn/components/ui/drawer'

	export let open = false
	export let selectedClass: PayloadStudentClassListItem | null = null

	let exams: PayloadClassExamListItem[] = []
	let loading = false

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

	const handleExamClick = (examId: number) => {
		open = false
		navigate(`/student/exam/${examId}/detail`)
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

		<div class="flex-1 px-4">
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
					<div class="space-y-4 pb-4">
						{#each exams as exam}
							<div onclick={() => handleExamClick(exam.exam.id)}>
								<ExamCard {exam} />
							</div>
						{/each}
					</div>
				</ScrollArea>
			{/if}
		</div>
	</DrawerContent>
</Drawer>
