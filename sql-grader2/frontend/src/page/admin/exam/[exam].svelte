<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate, Link } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import {
		ArrowLeftIcon,
		BookOpenIcon,
		DatabaseIcon,
		UsersIcon,
		HelpCircleIcon,
		CalendarIcon,
		ClockIcon,
		InfoIcon,
		Loader2Icon,
		FileTextIcon,
		BarChart3Icon,
		EditIcon,
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import Tab from '$/component/ui/Tab.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { formatDate, formatDateTime } from '$/util/format.ts'
	import type {
		PayloadExam,
		PayloadClass,
		PayloadCollection,
		PayloadExamAttemptCount,
	} from '$/util/backend/backend.ts'
	import TableStructureDialog from '../collection/dialog/TableStructureDialog.svelte'
	import ExamQuestions from './components/ExamQuestions.svelte'
	import ExamStudents from './components/ExamStudents.svelte'
	import ExamEditDialog from './dialog/ExamEditDialog.svelte'

	export let exam: number

	let examData: PayloadExam
	let classData: PayloadClass
	let collectionData: PayloadCollection
	let attemptCount: PayloadExamAttemptCount
	let loading = true
	let activeTab: 'questions' | 'students' = 'questions'
	let showTableDialog = false
	let editDialogOpen = false

	const loadExamDetail = () => {
		loading = true
		backend.admin
			.examDetail({ examId: exam })
			.then((response) => {
				if (response.success && response.data) {
					examData = response.data.exam
					classData = response.data.class
					collectionData = response.data.collection
					attemptCount = response.data.attemptCount
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleTabChange = (tab: 'questions' | 'students') => {
		activeTab = tab
	}

	const getTotalRows = () => {
		if (!collectionData?.metadata?.structure) return 0
		return collectionData.metadata.structure.reduce((total: number, table: any) => total + (table.rowCount || 0), 0)
	}

	const handleExamUpdated = () => {
		loadExamDetail()
		editDialogOpen = false
	}

	onMount(() => {
		loadExamDetail()
	})
</script>

<Container>
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else if !examData}
		<div class="flex min-h-[400px] flex-col items-center justify-center">
			<InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
			<h3 class="mb-2 text-lg font-semibold">Exam not found</h3>
			<p class="text-muted-foreground mb-4">The exam you're looking for doesn't exist</p>
			<Button onclick={() => navigate('/admin')}>
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back to Admin
			</Button>
		</div>
	{:else}
		<div class="flex justify-between">
			<div class="mb-6 flex flex-col gap-4">
				<button
					class="text-muted-foreground hover:text-primary flex items-center gap-2 hover:cursor-pointer"
					onclick={() => navigate(`/admin/class/${examData.classId}`)}
				>
					<ArrowLeftIcon size={16} />
					<span class="text-xs font-medium tracking-wide uppercase">{classData.name}</span>
				</button>
				<div class="flex items-center justify-between">
					<PageTitle title={examData.name} description={classData.name} />
				</div>
			</div>
			<Button variant="outline" size="sm" onclick={() => (editDialogOpen = true)}>
				<EditIcon class="mr-2 h-4 w-4" />
				Edit Exam
			</Button>
		</div>

		<!-- Info Cards Row -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
			<!-- Collection Info Card -->
			<Card class="transition-shadow hover:shadow-lg">
				<CardHeader class="pb-3">
					<CardTitle class="flex items-center gap-2">
						<DatabaseIcon class="h-5 w-5 text-green-600" />
						Collection Information
					</CardTitle>
					<CardDescription>Database collection details</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						<div>
							<h4 class="font-semibold">{collectionData.name}</h4>
							<p class="text-sm text-gray-600">{collectionData.questionCount} questions available</p>
						</div>

						{#if collectionData.metadata?.structure?.length > 0}
							<div class="flex items-center justify-between">
								<div class="text-sm text-gray-600">
									<p>{collectionData.metadata.structure.length} tables • {getTotalRows()} rows</p>
									{#if collectionData.metadata.schemaFilename}
										<p class="text-xs text-gray-500">{collectionData.metadata.schemaFilename}</p>
									{/if}
								</div>
								<Button variant="ghost" size="sm" onclick={() => (showTableDialog = true)}>
									Show Details
								</Button>
							</div>
						{/if}

						<div class="flex items-center justify-between border-t pt-3">
							<span class="text-sm text-gray-500">Created</span>
							<span class="text-sm">{formatDate(collectionData.createdAt)}</span>
						</div>
					</div>
				</CardContent>
			</Card>

			<!-- Exam Info Card -->
			<Card class="transition-shadow hover:shadow-lg">
				<CardHeader class="pb-3">
					<CardTitle class="flex items-center gap-2">
						<BookOpenIcon class="h-5 w-5 text-blue-600" />
						Exam Information
					</CardTitle>
					<CardDescription>Exam configuration and statistics</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						<!-- Access Code -->
						{#if examData.accessCode}
							<div class="flex items-center justify-between rounded-lg bg-blue-50 p-3">
								<span class="text-sm font-medium text-gray-700">Access Code</span>
								<span
									class="cursor-pointer font-mono text-lg font-bold text-blue-600 blur-sm transition-all hover:blur-none"
								>
									{examData.accessCode}
								</span>
							</div>
						{/if}

						<div class="flex items-center gap-2 text-sm">
							<HelpCircleIcon class="h-4 w-4" />
							<span class="font-medium">
								{examData.questionCount} / {collectionData.questionCount} questions
							</span>
						</div>

						<!-- Attempt Statistics -->
						<div class="rounded-lg bg-gray-50 p-3">
							<h5 class="mb-2 text-sm font-medium text-gray-700">Attempt Statistics</h5>
							<div class="grid grid-cols-3 gap-3 text-center">
								<div>
									<div class="text-lg font-semibold text-blue-600">{attemptCount.openedCount}</div>
									<div class="text-xs text-gray-500">Opened</div>
								</div>
								<div>
									<div class="text-lg font-semibold text-yellow-600">{attemptCount.startedCount}</div>
									<div class="text-xs text-gray-500">Started</div>
								</div>
								<div>
									<div class="text-lg font-semibold text-green-600">{attemptCount.finishedCount}</div>
									<div class="text-xs text-gray-500">Finished</div>
								</div>
							</div>
						</div>

						<div class="space-y-2 border-t pt-3 text-sm">
							<div class="flex items-center justify-between">
								<span class="flex items-center gap-1 text-gray-500">
									<CalendarIcon class="h-3 w-3" />
									Created
								</span>
								<span>{formatDate(examData.createdAt)}</span>
							</div>
							<div class="flex items-center justify-between">
								<span class="flex items-center gap-1 text-gray-500">
									<ClockIcon class="h-3 w-3" />
									Opens
								</span>
								<span>{examData.openedAt ? formatDateTime(examData.openedAt) : 'Not set'}</span>
							</div>
							<div class="flex items-center justify-between">
								<span class="flex items-center gap-1 text-gray-500">
									<ClockIcon class="h-3 w-3" />
									Closes
								</span>
								<span>{examData.closedAt ? formatDateTime(examData.closedAt) : 'Not set'}</span>
							</div>
						</div>
					</div>
				</CardContent>
			</Card>
		</div>

		<!-- Tab Navigation -->
		<Tab
			tabs={[
				{ id: 'questions', label: 'Questions', icon: HelpCircleIcon },
				{ id: 'students', label: 'Students', icon: UsersIcon },
			]}
			{activeTab}
			on:change={(e) => handleTabChange(e.detail)}
		/>

		<!-- Tab Content -->
		{#if activeTab === 'questions'}
			<ExamQuestions examId={exam} {collectionData} />
		{:else if activeTab === 'students'}
			<ExamStudents examId={exam} />
		{/if}
	{/if}
</Container>

<TableStructureDialog bind:open={showTableDialog} structure={collectionData?.metadata?.structure || []} />

<ExamEditDialog bind:open={editDialogOpen} {examData} on:updated={handleExamUpdated} />
