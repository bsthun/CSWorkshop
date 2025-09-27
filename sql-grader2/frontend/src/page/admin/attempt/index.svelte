<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$/lib/shadcn/components/ui/card'
	import { Avatar, AvatarImage, AvatarFallback } from '$/lib/shadcn/components/ui/avatar'
	import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '$/lib/shadcn/components/ui/table'
	import * as HoverCard from '$/lib/shadcn/components/ui/hover-card/index.js'
	import UserProfile from '$/component/share/UserProfile.svelte'
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
		Loader2Icon,
		CheckCircle2Icon,
		CircleSlashIcon,
		CircleXIcon,
		CircleIcon,
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { formatDateTime } from '$/util/format.ts'
	import type { PayloadSubmissionListItem } from '$/util/backend/backend.ts'
	import SubmissionDetailDialog from './dialog/SubmissionDetailDialog.svelte'

	import { useLocation } from 'svelte-navigator'

	const location = useLocation()

	const urlParams = new URLSearchParams($location.search)
	const examAttemptId = urlParams.get('examAttemptId')
	const examQuestionId = urlParams.get('examQuestionId')

	const group: 'examAttempt' | 'examQuestion' = examQuestionId ? 'examQuestion' : 'examAttempt'

	let submissions: PayloadSubmissionListItem[] = []
	let loading = true
	let selectedSubmission: PayloadSubmissionListItem | null = null
	let showDetailDialog = false

	let studentInfo: any = null
	let examInfo: any = null

	const loadSubmissions! = () => {
		loading = true
		const params: any = { examAttemptId: examAttemptId }
		if (examQuestionId) {
			params.examQuestionId = examQuestionId
		}

		backend.admin
			.submissionList(params)
			.then((response) => {
				submissions = response.data.submissions || []
				if (submissions.length > 0) {
					studentInfo = submissions[0].student
					examInfo = submissions[0].exam
				}
			})
			.catch(catcher)
			.finally(() => {
				loading = false
			})
	}

	const getStatusIcon = (submission: any) => {
		if (submission.checkQueryPassed && submission.checkPromptPassed) {
			return { icon: CheckCircle2Icon, color: 'text-green-600', bg: 'bg-green-100' }
		} else if (submission.checkQueryPassed && !submission.checkPromptPassed) {
			return { icon: CircleSlashIcon, color: 'text-yellow-600', bg: 'bg-yellow-100' }
		} else if (!submission.checkQueryPassed && !submission.checkPromptPassed) {
			return { icon: CircleXIcon, color: 'text-red-600', bg: 'bg-red-100' }
		} else {
			return { icon: CircleIcon, color: 'text-gray-400', bg: 'bg-gray-100' }
		}
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
					<PageTitle title="Submission Details" description="" />
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
						<div class="mb-4 flex items-center gap-4">
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
									<span
										>{submissions[0].attempt.startedAt
											? formatDateTime(submissions[0].attempt.startedAt)
											: 'Not started'}</span
									>
								</div>
								<div class="flex items-center justify-between">
									<span class="text-gray-500">Finished</span>
									<span
										>{submissions[0].attempt.finishedAt
											? formatDateTime(submissions[0].attempt.finishedAt)
											: 'In progress'}</span
									>
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
					{examQuestionId ? 'All submissions for this question' : 'All submissions for this attempt'}
				</CardDescription>
			</CardHeader>
			<CardContent>
				{#if submissions.length === 0}
					<div class="flex flex-col items-center justify-center py-12">
						<FileTextIcon class="mb-4 h-16 w-16 text-gray-400" />
						<h3 class="mb-2 text-lg font-semibold">No Submissions</h3>
						<p class="text-muted-foreground text-center">No submissions found for this attempt.</p>
					</div>
				{:else}
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>{group === 'examAttempt' ? 'Question' : 'Student'}</TableHead>
								<TableHead>Status</TableHead>
								<TableHead>Answer</TableHead>
								<TableHead>Message</TableHead>
								<TableHead>Submitted</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{#each submissions as submission}
								{@const status = getStatusIcon(submission.submission)}
								<TableRow class="hover:bg-gray-50">
									<TableCell>
										{#if group === 'examAttempt'}
											<div>
												<div class="font-medium">
													{submission.question.title}
												</div>
												<div class="line-clamp-1 max-w-xs truncate text-sm text-gray-500">
													{submission.question.description}
												</div>
											</div>
										{:else}
											<UserProfile user={submission.student} />
										{/if}
									</TableCell>
									<TableCell>
										<div class="flex justify-center">
											<svelte:component this={status.icon} size={16} class={status.color} />
										</div>
									</TableCell>
									<TableCell>
										{#if submission.submission.answer}
											<div class="max-w-xs">
												<code class="line-clamp-2 rounded bg-gray-100 px-2 py-1 text-xs">
													{submission.submission.answer.substring(0, 100)}{submission
														.submission.answer.length > 100
														? '...'
														: ''}
												</code>
											</div>
										{:else}
											<span class="text-sm text-gray-400">No answer</span>
										{/if}
									</TableCell>
									<TableCell>
										<div class="max-w-48">
											{#if submission.submission.result?.promptDescription}
												<HoverCard.Root>
													<HoverCard.Trigger>
														<p
															class="cursor-pointer truncate font-mono text-xs {submission
																.submission.checkPromptPassed
																? 'text-green-700'
																: 'text-red-700'}"
														>
															{submission.submission.result.promptDescription}
														</p>
													</HoverCard.Trigger>
													<HoverCard.Content>
														<p class="font-mono text-sm">
															{submission.submission.result.promptDescription}
														</p>
													</HoverCard.Content>
												</HoverCard.Root>
											{:else if submission.submission.result?.executionError}
												<HoverCard.Root>
													<HoverCard.Trigger>
														<p
															class="cursor-pointer truncate font-mono text-xs text-red-700"
														>
															{submission.submission.result.executionError}
														</p>
													</HoverCard.Trigger>
													<HoverCard.Content>
														<pre class="font-mono text-sm whitespace-pre-wrap">{submission
																.submission.result.executionError}</pre>
													</HoverCard.Content>
												</HoverCard.Root>
											{:else if submission.submission.result?.promptError}
												<HoverCard.Root>
													<HoverCard.Trigger>
														<p
															class="cursor-pointer truncate font-mono text-xs text-red-700"
														>
															{submission.submission.result.promptError}
														</p>
													</HoverCard.Trigger>
													<HoverCard.Content>
														<pre class="font-mono text-sm whitespace-pre-wrap">{submission
																.submission.result.promptError}</pre>
													</HoverCard.Content>
												</HoverCard.Root>
											{/if}
										</div>
									</TableCell>
									<TableCell>
										{#if submission.submission.createdAt}
											<span class="text-sm"
												>{formatDateTime(submission.submission.createdAt)}</span
											>
										{:else}
											<span class="text-sm text-gray-400">Not submitted</span>
										{/if}
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

<SubmissionDetailDialog bind:open={showDetailDialog} submission={selectedSubmission} />
