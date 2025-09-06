<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import {
		ArrowLeftIcon,
		BookOpenIcon,
		DatabaseIcon,
		CalendarIcon,
		ClockIcon,
		InfoIcon,
		Loader2Icon,
		UserIcon,
		TerminalIcon,
	} from 'lucide-svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import ExamAttemptDetailCredentialField from '../../components/ExamAttemptDetailCredentialField.svelte'
	import { backend, catcher } from '$/util/backend'
	import { formatDateTime } from '$/util/format'
	import type {
		PayloadClassExamDetailResponse,
		PayloadExam,
		PayloadClass,
		PayloadExamCredential,
	} from '$/util/backend/backend.ts'

	export let examAttemptId: number

	let examDetailData: PayloadClassExamDetailResponse
	let examData: PayloadExam
	let classData: PayloadClass
	let credentialData: PayloadExamCredential
	let loading = true
	let copiedField: string | null = null

	const loadExamAttemptDetail = () => {
		loading = true
		backend.student
			.classExamAttemptDetail({ examAttemptId })
			.then((response) => {
				if (response.success && response.data) {
					examDetailData = response.data
					examData = response.data.exam
					classData = response.data.class
					credentialData = response.data.credential
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const copyToClipboard = async (text: string, field: string) => {
		await navigator.clipboard.writeText(text)
		copiedField = field
		setTimeout(() => {
			copiedField = null
		}, 2000)
	}

	onMount(() => {
		loadExamAttemptDetail()
	})
</script>

<Container>
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else if !examDetailData}
		<div class="flex min-h-[400px] flex-col items-center justify-center">
			<InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
			<h3 class="mb-2 text-lg font-semibold">Exam attempt not found</h3>
			<p class="text-muted-foreground mb-4">The exam attempt you're looking for doesn't exist</p>
			<Button onclick={() => navigate('/student')}>
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back to Classes
			</Button>
		</div>
	{:else}
		<div class="mb-6 flex flex-col gap-4">
			<button
				class="text-muted-foreground hover:text-primary flex items-center gap-2 hover:cursor-pointer"
				onclick={() => navigate('/student')}
			>
				<ArrowLeftIcon size={16} />
				<span class="text-xs font-medium tracking-wide uppercase">EXAMS</span>
			</button>
			<div class="flex items-center justify-between">
				<PageTitle title={examData.name} description={`${classData.code} - ${classData.name}`} />
			</div>
		</div>

		<!-- Status Badge -->
		<div class="mb-6">
			<Badge variant="default" class="px-3 py-1 text-sm">
				<TerminalIcon size={14} class="mr-2" />
				Database Connection Available
			</Badge>
		</div>

		<!-- Info Cards Row -->
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
			<div class="flex flex-col gap-6">
				<Card class="transition-shadow hover:shadow-lg">
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<BookOpenIcon class="h-5 w-5 text-blue-600" />
							Exam Information
						</CardTitle>
						<CardDescription>Exam details and schedule</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="space-y-4">
							<div>
								<h4 class="font-semibold">{examData.name}</h4>
								<p class="text-sm text-gray-600">{examDetailData.examQuestionCount} questions</p>
							</div>
							<div class="space-y-2 border-t pt-3 text-sm">
								<div class="flex items-center justify-between">
									<span class="flex items-center gap-1 text-gray-500">
										<CalendarIcon class="h-3 w-3" />
										Opens
									</span>
									<span>{formatDateTime(examData.openedAt)}</span>
								</div>
								<div class="flex items-center justify-between">
									<span class="flex items-center gap-1 text-gray-500">
										<ClockIcon class="h-3 w-3" />
										Closes
									</span>
									<span>{formatDateTime(examData.closedAt)}</span>
								</div>
							</div>
						</div>
					</CardContent>
				</Card>
				<Card class="transition-shadow hover:shadow-lg">
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<UserIcon class="h-5 w-5 text-purple-600" />
							Class Information
						</CardTitle>
						<CardDescription>Course details</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="flex items-start justify-between">
							<div>
								<h4 class="font-semibold">{classData.code}</h4>
								<p class="text-sm text-gray-600">{classData.name}</p>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<!-- Database Connection Card -->
			<Card class="transition-shadow hover:shadow-lg">
				<CardHeader class="pb-3">
					<CardTitle class="flex items-center gap-2">
						<DatabaseIcon class="h-5 w-5 text-green-600" />
						Database Connection
					</CardTitle>
					<CardDescription>Connection details for this exam</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						<div class="rounded-lg bg-gray-50 p-4">
							<div class="space-y-3 font-mono text-sm">
								<ExamAttemptDetailCredentialField
									name="Dialect"
									value={credentialData.dialect}
									copy={false}
									{copiedField}
									onCopy={copyToClipboard}
								/>
								<ExamAttemptDetailCredentialField
									name="Host"
									value={credentialData.host}
									copy={true}
									{copiedField}
									onCopy={copyToClipboard}
								/>
								<ExamAttemptDetailCredentialField
									name="Port"
									value={credentialData.port}
									copy={true}
									{copiedField}
									onCopy={copyToClipboard}
								/>
								<ExamAttemptDetailCredentialField
									name="Database"
									value={credentialData.databaseName}
									copy={true}
									{copiedField}
									onCopy={copyToClipboard}
								/>
								<ExamAttemptDetailCredentialField
									name="Username"
									value={credentialData.user}
									copy={true}
									{copiedField}
									onCopy={copyToClipboard}
								/>
								<ExamAttemptDetailCredentialField
									name="Password"
									value={credentialData.password}
									blur={true}
									copy={true}
									{copiedField}
									onCopy={copyToClipboard}
								/>
							</div>
						</div>

						<div class="border-t pt-3">
							<p class="text-xs text-gray-500">
								Use these credentials to connect to the database for your exam
							</p>
						</div>
					</div>
				</CardContent>
			</Card>
		</div>

		<!-- Instructions Card -->
		<Card>
			<CardHeader>
				<CardTitle class="flex items-center gap-2">
					<InfoIcon class="h-5 w-5 text-blue-600" />
					Instructions
				</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-4">
					<div>
						<h5 class="mb-2 font-medium">Getting Started</h5>
						<ul class="list-inside list-disc space-y-1 text-sm text-gray-600">
							<li>Use the database connection details above to connect to your exam database</li>
							<li>Answer all {examDetailData.examQuestionCount} questions in the time allotted</li>
							<li>Submit your SQL queries for each question</li>
							<li>Make sure to test your queries before submitting</li>
						</ul>
					</div>
					<div class="border-t pt-4">
						<h5 class="mb-2 font-medium">Important Notes</h5>
						<ul class="list-inside list-disc space-y-1 text-sm text-gray-600">
							<li>Your database connection is active until the exam closes</li>
							<li>Save your work frequently</li>
							<li>Contact your instructor if you encounter technical difficulties</li>
						</ul>
					</div>
				</div>
			</CardContent>
		</Card>
	{/if}
</Container>
