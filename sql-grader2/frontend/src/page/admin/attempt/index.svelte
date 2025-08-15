<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Avatar, AvatarImage, AvatarFallback } from '$/lib/shadcn/components/ui/avatar'
	import {
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow
	} from '$/lib/shadcn/components/ui/table'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import {
		ArrowLeftIcon,
		UserIcon,
		BookOpenIcon,
		FileTextIcon,
		CheckCircleIcon,
		XCircleIcon,
		AlertTriangleIcon,
		MinusCircleIcon,
		EyeIcon,
		ClockIcon,
		CalendarIcon,
		InfoIcon,
		Loader2Icon
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { formatDateTime } from '$/util/format.ts'
	import type { PayloadSubmissionListItem } from '$/util/backend/backend.ts'
	import SubmissionDetailDialog from './dialog/SubmissionDetailDialog.svelte'

	// Get query parameters
	const urlParams = new URLSearchParams(window.location.search)
	const examAttemptId = urlParams.get('examAttemptId')
	const examQuestionId = urlParams.get('examQuestionId')

	let submissions: PayloadSubmissionListItem[] = []
	let loading = true
	let selectedSubmission: PayloadSubmissionListItem | null = null
	let showDetailDialog = false

	// Student and exam info (from first submission)
	let studentInfo: any = null
	let examInfo: any = null

	const loadSubmissions = () => {
		if (!examAttemptId) {
			navigate('/admin')
			return
		}

		loading = true
		const params: any = { examAttemptId: parseInt(examAttemptId) }
		if (examQuestionId) {
			params.examQuestionId = parseInt(examQuestionId)
		}

		backend.admin
			.submissionList(params)
			.then((response) => {
				if (response.success && response.data) {
					submissions = response.data.submissions || []
					
					// Extract student and exam info from first submission
					if (submissions.length > 0) {
						studentInfo = submissions[0].student
						examInfo = submissions[0].exam
					}
				}
			})
			.catch(catcher)
			.finally(() => {
				loading = false
			})
	}

	const getStatusIcon = (submission: any) => {
		if (submission.checkPromptPassed && submission.checkQueryPassed) {
			return { icon: CheckCircleIcon, color: 'text-green-600', bg: 'bg-green-100' }
		} else if (submission.checkPromptPassed === false || submission.checkQueryPassed === false) {
			return { icon: XCircleIcon, color: 'text-red-600', bg: 'bg-red-100' }
		} else if (submission.answer) {
			return { icon: AlertTriangleIcon, color: 'text-yellow-600', bg: 'bg-yellow-100' }
		} else {
			return { icon: MinusCircleIcon, color: 'text-gray-400', bg: 'bg-gray-100' }
		}
	}

	const getStatusText = (submission: any) => {
		if (submission.checkPromptPassed && submission.checkQueryPassed) {
			return 'Passed'
		} else if (submission.checkPromptPassed === false || submission.checkQueryPassed === false) {
			return 'Failed'
		} else if (submission.answer) {
			return 'Submitted'
		} else {
			return 'Not Submitted'
		}
	}

	const getStatusVariant = (submission: any): any => {
		if (submission.checkPromptPassed && submission.checkQueryPassed) {
			return 'default'
		} else if (submission.checkPromptPassed === false || submission.checkQueryPassed === false) {
			return 'destructive'
		} else if (submission.answer) {
			return 'secondary'
		} else {
			return 'outline'
		}
	}

	const handleViewDetail = (submission: PayloadSubmissionListItem) => {
		selectedSubmission = submission
		showDetailDialog = true
	}

	const getInitials = (firstname: string, lastname: string) => {
		return `${firstname.charAt(0)}${lastname.charAt(0)}`.toUpperCase()
	}

	onMount(() => {
		loadSubmissions()
	})
</script>

<Container>
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else if !examAttemptId}
		<div class="flex min-h-[400px] flex-col items-center justify-center">
			<InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
			<h3 class="mb-2 text-lg font-semibold">Missing Parameters</h3>
			<p class="text-muted-foreground mb-4">Exam attempt ID is required</p>
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
				<span class="text-xs font-medium tracking-wide uppercase">ATTEMPT DETAILS</span>
			</button>
			<div class="flex items-center justify-between">
				{#if studentInfo && examInfo}
					<PageTitle 
						title="{studentInfo.firstname} {studentInfo.lastname}"
						description="Exam: {examInfo.name}"
					/>
				{:else}
					<PageTitle title="Submission Details" description="Loading..." />
				{/if}
			</div>
		</div>

		{#if studentInfo && examInfo}
			<!-- Student & Exam Info Cards -->
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
				<!-- Student Info -->
				<Card class="transition-shadow hover:shadow-lg">
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<UserIcon class="h-5 w-5 text-blue-600" />
							Student Information
						</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="flex items-center gap-4 mb-4">
							<Avatar class="h-12 w-12">
								<AvatarImage src={studentInfo.pictureUrl} alt="" />
								<AvatarFallback>
									{getInitials(studentInfo.firstname, studentInfo.lastname)}
								</AvatarFallback>
							</Avatar>
							<div>
								<h4 class="font-semibold">{studentInfo.firstname} {studentInfo.lastname}</h4>
								<p class="text-sm text-gray-600">{studentInfo.email}</p>
							</div>
						</div>
						<div class="space-y-2 text-sm">
							<div class="flex items-center justify-between">
								<span class="text-gray-500">Student ID</span>
								<span>{studentInfo.id}</span>
							</div>
							{#if submissions[0]?.attempt}
								<div class="flex items-center justify-between">
									<span class="text-gray-500">Started</span>
									<span>{submissions[0].attempt.startedAt ? formatDateTime(submissions[0].attempt.startedAt) : 'Not started'}</span>
								</div>
								<div class="flex items-center justify-between">
									<span class="text-gray-500">Finished</span>
									<span>{submissions[0].attempt.finishedAt ? formatDateTime(submissions[0].attempt.finishedAt) : 'In progress'}</span>
								</div>
							{/if}
						</div>
					</CardContent>
				</Card>

				<!-- Exam Info -->
				<Card class="transition-shadow hover:shadow-lg">
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<BookOpenIcon class="h-5 w-5 text-green-600" />
							Exam Information
						</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="space-y-3">
							<div>
								<h4 class="font-semibold">{examInfo.name}</h4>
								<p class="text-sm text-gray-600">{examInfo.questionCount} questions</p>
							</div>
							<div class="space-y-2 text-sm">
								<div class="flex items-center justify-between">
									<span class="flex items-center gap-1 text-gray-500">
										<CalendarIcon class="h-3 w-3" />
										Opens
									</span>
									<span>{examInfo.openedAt ? formatDateTime(examInfo.openedAt) : 'Not set'}</span>
								</div>
								<div class="flex items-center justify-between">
									<span class="flex items-center gap-1 text-gray-500">
										<ClockIcon class="h-3 w-3" />
										Closes
									</span>
									<span>{examInfo.closedAt ? formatDateTime(examInfo.closedAt) : 'Not set'}</span>
								</div>
								<div class="flex items-center justify-between">
									<span class="text-gray-500">Access Code</span>
									<span class="font-mono text-xs">{examInfo.accessCode}</span>
								</div>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		{/if}

		<!-- Submissions Table -->
		<Card>
			<CardHeader class="pb-3">
				<CardTitle class="flex items-center gap-2">
					<FileTextIcon class="h-5 w-5" />
					Submissions ({submissions.length})
				</CardTitle>
				<CardDescription>
					{examQuestionId ? 'Showing submissions for a specific question' : 'All submissions for this attempt'}
				</CardDescription>
			</CardHeader>
			<CardContent>
				{#if submissions.length === 0}
					<div class="flex flex-col items-center justify-center py-12">
						<FileTextIcon class="mb-4 h-16 w-16 text-gray-400" />
						<h3 class="mb-2 text-lg font-semibold">No Submissions</h3>
						<p class="text-muted-foreground text-center">
							No submissions found for this attempt.
						</p>
					</div>
				{:else}
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Question</TableHead>
								<TableHead>Status</TableHead>
								<TableHead>Answer Preview</TableHead>
								<TableHead>Submitted</TableHead>
								<TableHead>Checked</TableHead>
								<TableHead class="text-right">Actions</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{#each submissions as submission}
								{@const status = getStatusIcon(submission.submission)}
								<TableRow class="hover:bg-gray-50">
									<TableCell>
										<div>
											<div class="font-medium">
												{submission.question.title || `Question ${submission.question.orderNum || submission.question.id}`}
											</div>
											{#if submission.question.description}
												<div class="text-sm text-gray-500 line-clamp-1">
													{submission.question.description}
												</div>
											{/if}
										</div>
									</TableCell>
									<TableCell>
										<Badge variant={getStatusVariant(submission.submission)} class="flex items-center gap-1 w-fit">
											<svelte:component this={status.icon} class="h-3 w-3 {status.color}" />
											{getStatusText(submission.submission)}
										</Badge>
									</TableCell>
									<TableCell>
										{#if submission.submission.answer}
											<div class="max-w-xs">
												<code class="text-xs bg-gray-100 px-2 py-1 rounded line-clamp-2">
													{submission.submission.answer.substring(0, 100)}{submission.submission.answer.length > 100 ? '...' : ''}
												</code>
											</div>
										{:else}
											<span class="text-gray-400 text-sm">No answer</span>
										{/if}
									</TableCell>
									<TableCell>
										{#if submission.submission.createdAt}
											<span class="text-sm">{formatDateTime(submission.submission.createdAt)}</span>
										{:else}
											<span class="text-gray-400 text-sm">Not submitted</span>
										{/if}
									</TableCell>
									<TableCell>
										{#if submission.submission.checkQueryAt || submission.submission.checkPromptAt}
											<div class="text-xs space-y-1">
												{#if submission.submission.checkPromptAt}
													<div>Prompt: {formatDateTime(submission.submission.checkPromptAt)}</div>
												{/if}
												{#if submission.submission.checkQueryAt}
													<div>Query: {formatDateTime(submission.submission.checkQueryAt)}</div>
												{/if}
											</div>
										{:else}
											<span class="text-gray-400 text-sm">Not checked</span>
										{/if}
									</TableCell>
									<TableCell class="text-right">
										<Button 
											variant="ghost" 
											size="sm" 
											onclick={() => handleViewDetail(submission)}
											disabled={!submission.submission.answer}
										>
											<EyeIcon class="h-4 w-4 mr-2" />
											Details
										</Button>
									</TableCell>
								</TableRow>
							{/each}
						</TableBody>
					</Table>
				{/if}
			</CardContent>
		</Card>
	{/if}
</Container>

<SubmissionDetailDialog 
	bind:open={showDetailDialog} 
	submission={selectedSubmission}
/>