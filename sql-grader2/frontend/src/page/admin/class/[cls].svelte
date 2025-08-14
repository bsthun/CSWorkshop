<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import {
		ArrowLeftIcon,
		UsersIcon,
		BookOpenIcon,
		EditIcon,
		InfoIcon,
		Loader2Icon,
		CopyIcon,
		PlusIcon,
		HelpCircleIcon,
		CheckCircleIcon,
		CalendarIcon,
		ClockIcon,
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { toast } from 'svelte-sonner'
	import { formatDate, formatDateTime } from '$/util/format.ts'
	import type {
		PayloadClass,
		PayloadClassJoinee,
		PayloadExamListItem,
		PayloadSemesterInfo,
	} from '$/util/backend/backend.ts'
	import EditClassDialog from './dialog/EditClassDialog.svelte'
	import CreateExamDialog from './dialog/CreateExamDialog.svelte'

	export let cls: number

	let classData: PayloadClass
	let semester: PayloadSemesterInfo
	let joinees: PayloadClassJoinee[]
	let exams: PayloadExamListItem[]
	let loading = true
	let loadingExams = false
	let showEditDialog = false
	let showCreateExamDialog = false
	let activeTab: 'students' | 'exams' = 'students'

	const loadClass = () => {
		loading = true
		backend.admin
			.classDetail({ classId: cls })
			.then((response) => {
				if (response.success && response.data) {
					classData = response.data.class
					semester = response.data.semester
					joinees = response.data.joinees
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const loadExams = () => {
		if (activeTab !== 'exams') return

		loadingExams = true
		backend.admin
			.examList({ classId: cls })
			.then((response) => {
				if (response.success && response.data) {
					exams = response.data.exams
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loadingExams = false
			})
	}

	const handleClassEdited = () => {
		loadClass()
	}

	const handleExamCreated = () => {
		loadExams()
	}

	const handleTabChange = (tab: 'students' | 'exams') => {
		activeTab = tab
		if (tab === 'exams') {
			loadExams()
		}
	}

	const copyRegisterCode = () => {
		if (classData?.registerCode) {
			navigator.clipboard.writeText(classData.registerCode)
			toast.success('Register code copied to clipboard')
		}
	}


	onMount(() => {
		loadClass()
	})
</script>

<Container>
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else if !classData}
		<div class="flex min-h-[400px] flex-col items-center justify-center">
			<InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
			<h3 class="mb-2 text-lg font-semibold">Class not found</h3>
			<p class="text-muted-foreground mb-4">The class you're looking for doesn't exist</p>
			<Button onclick={() => navigate('/admin')}>
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back to Admin
			</Button>
		</div>
	{:else}
		<div class="mb-6 flex flex-col gap-4">
			<button
				class="text-muted-foreground hover:text-primary flex items-center gap-2 hover:cursor-pointer"
				onclick={() => navigate('/admin')}
			>
				<ArrowLeftIcon size={16} />
				<span class="text-xs font-medium tracking-wide uppercase">CLASS</span>
			</button>
			<div class="flex items-center justify-between">
				<PageTitle title={classData.name} description={`${classData.code} • ${semester.name}`} />
				<Button variant="outline" onclick={() => (showEditDialog = true)}>
					<EditIcon class="mr-2 h-4 w-4" />
					Edit
				</Button>
			</div>
		</div>

		<div class="mb-6 flex w-fit items-center gap-3 rounded-md border bg-gray-50 px-4 py-2">
			<span class="text-sm text-gray-600">Register Code</span>
			<code class="font-mono font-semibold">{classData.registerCode}</code>
			<Button size="sm" variant="ghost" onclick={copyRegisterCode}>
				<CopyIcon class="h-4 w-4" />
			</Button>
		</div>

		<!-- Tab Navigation -->
		<div class="mb-6 flex items-center justify-between">
			<div class="flex space-x-1 rounded-lg bg-gray-100 p-1">
				<button
					class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === 'students'
						? 'bg-white text-gray-900 shadow-sm'
						: 'text-gray-500 hover:text-gray-900'}"
					onclick={() => handleTabChange('students')}
				>
					<span class="flex items-center gap-2">
						<UsersIcon size={16} />
						Students
					</span>
				</button>
				<button
					class="rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab === 'exams'
						? 'bg-white text-gray-900 shadow-sm'
						: 'text-gray-500 hover:text-gray-900'}"
					onclick={() => handleTabChange('exams')}
				>
					<span class="flex items-center gap-2">
						<BookOpenIcon size={16} />
						Exams
					</span>
				</button>
			</div>

			{#if activeTab === 'exams'}
				<Button class="gap-2" onclick={() => (showCreateExamDialog = true)}>
					<PlusIcon class="h-4 w-4" />
					Create Exam
				</Button>
			{/if}
		</div>

		<!-- Tab Content -->
		{#if activeTab === 'students'}
			<div class="space-y-4">
				<h3 class="text-lg font-semibold">Students ({joinees.length})</h3>
				{#if joinees.length === 0}
					<div class="flex items-center justify-center py-8">
						<div class="text-center">
							<UsersIcon class="mx-auto mb-2 h-12 w-12 text-gray-400" />
							<p class="text-gray-500">No students enrolled yet</p>
						</div>
					</div>
				{:else}
					<div class="space-y-2">
						{#each joinees as joinee}
							<Card>
								<CardContent class="py-3">
									<div class="flex items-center justify-between">
										<div>
											<p class="font-medium">{joinee.user.firstname} {joinee.user.lastname}</p>
											<p class="text-sm text-gray-600">{joinee.user.email}</p>
										</div>
										<div class="text-right text-sm text-gray-500">
											<p>User ID: {joinee.user.id}</p>
										</div>
									</div>
								</CardContent>
							</Card>
						{/each}
					</div>
				{/if}
			</div>
		{:else if activeTab === 'exams'}
			<div class="space-y-4">
				<h3 class="text-lg font-semibold">Exams</h3>
				{#if loadingExams}
					<div class="flex items-center justify-center py-8">
						<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
					</div>
				{:else if exams.length === 0}
					<div class="flex items-center justify-center py-8">
						<div class="text-center">
							<BookOpenIcon class="mx-auto mb-2 h-12 w-12 text-gray-400" />
							<p class="mb-4 text-gray-500">No exams created yet</p>
							<Button class="gap-2" onclick={() => (showCreateExamDialog = true)}>
								<PlusIcon class="h-4 w-4" />
								Create First Exam
							</Button>
						</div>
					</div>
				{:else}
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each exams as examItem}
							<Card class="cursor-pointer transition-shadow hover:shadow-lg">
								<CardHeader class="gap-2">
									<CardTitle class="flex items-center gap-1">
										<BookOpenIcon size={12} class="text-muted-foreground" />
										<span class="text-muted-foreground text-xs font-medium tracking-wide uppercase">
											EXAM
										</span>
									</CardTitle>
									<h3 class="text-lg font-semibold">{examItem.exam.name}</h3>
									<div class="flex items-center gap-1 text-sm text-gray-600">
										<span class="text-muted-foreground text-xs font-medium tracking-wide uppercase">
											COLLECTION
										</span>
									</div>
									<h4 class="font-semibold">{examItem.collection.name}</h4>
								</CardHeader>
								<CardContent>
									<div class="space-y-3 text-sm">
										<div class="flex items-center gap-2 text-gray-600">
											<HelpCircleIcon class="h-4 w-4" />
											<span class="font-medium">
												{examItem.questionCount} / {examItem.collection.questionCount} questions
											</span>
										</div>

										<div class="grid grid-cols-2 gap-3 text-gray-500">
											<!-- Created -->
											<div class="flex flex-col gap-1">
												<div class="flex items-center gap-1">
													<CalendarIcon class="h-3 w-3" />
													<span class="text-xs font-medium text-gray-400 uppercase"
														>Created</span
													>
												</div>
												<span class="text-sm">{formatDate(examItem.exam.createdAt)}</span>
											</div>

											<!-- Opens -->
											<div class="flex flex-col gap-1">
												<div class="flex items-center gap-1">
													<ClockIcon class="h-3 w-3" />
													<span class="text-xs font-medium text-gray-400 uppercase"
														>Opens</span
													>
												</div>
												<span class="text-sm">
													{examItem.exam.openedAt
														? formatDateTime(examItem.exam.openedAt)
														: 'Not set'}
												</span>
											</div>

											<!-- Updated -->
											<div class="flex flex-col gap-1">
												<div class="flex items-center gap-1">
													<CalendarIcon class="h-3 w-3" />
													<span class="text-xs font-medium text-gray-400 uppercase"
														>Updated</span
													>
												</div>
												<span class="text-sm">{formatDate(examItem.exam.updatedAt)}</span>
											</div>

											<!-- Closes -->
											<div class="flex flex-col gap-1">
												<div class="flex items-center gap-1">
													<ClockIcon class="h-3 w-3" />
													<span class="text-xs font-medium text-gray-400 uppercase"
														>Closes</span
													>
												</div>
												<span class="text-sm">
													{examItem.exam.closedAt
														? formatDateTime(examItem.exam.closedAt)
														: 'Not set'}
												</span>
											</div>
										</div>
									</div>
								</CardContent>
							</Card>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	{/if}
</Container>

{#if classData}
	<EditClassDialog bind:open={showEditDialog} {classData} {semester} on:edited={handleClassEdited} />
	<CreateExamDialog bind:open={showCreateExamDialog} classId={cls} on:created={handleExamCreated} />
{/if}
