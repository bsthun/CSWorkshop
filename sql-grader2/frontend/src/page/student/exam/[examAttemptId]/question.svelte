<script lang="ts">
	import { onMount } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import {
		ArrowLeftIcon,
		CheckCircle2Icon,
		CircleSlashIcon,
		CircleXIcon,
		CircleIcon,
		Loader2Icon,
	} from 'lucide-svelte'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Textarea } from '$/lib/shadcn/components/ui/textarea'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import PageTitle from '$/component/ui/PageTitle.svelte'
	import { backend, catcher } from '$/util/backend'
	import {
		type PayloadClass,
		type PayloadExam,
		type PayloadExamQuestionWithStatus,
		PayloadExamQuestionWithStatusStatusEnum,
	} from '$/util/backend/backend.ts'
	import { toast } from 'svelte-sonner'
	import SubmissionList from './component/SubmissionList.svelte'
	import ExamTimer from './component/ExamTimer.svelte'
	import { marked } from 'marked'

	export let examAttemptId: number

	let loading = true
	let questions: PayloadExamQuestionWithStatus[] = []
	let selectedQuestion: PayloadExamQuestionWithStatus | null = null
	let submissionAnswer = ''
	let submitting = false
	let examData: PayloadExam
	let classData: PayloadClass
	let loadingDetail = false
	let questionDetail: any = null

	const loadExamData = () => {
		loading = true
		Promise.all([
			backend.student.classExamAttemptDetail({ examAttemptId }),
			backend.student.studentExamQuestionList({ examAttemptId }),
		])
			.then(([detailResponse, questionResponse]) => {
				if (detailResponse.success && detailResponse.data) {
					examData = detailResponse.data.exam
					classData = detailResponse.data.class
				}
				questions = questionResponse.data.questions
				if (questions.length > 0 && !selectedQuestion) {
					selectQuestion(questions[0])
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const selectQuestion = (question: PayloadExamQuestionWithStatus) => {
		selectedQuestion = question
		submissionAnswer = ''
		loadingDetail = true

		backend.student
			.studentExamQuestionDetail({
				examAttemptId,
				examQuestionId: question.examQuestion.id,
			})
			.then((response) => {
				if (response.success && response.data) {
					questionDetail = response.data
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loadingDetail = false
			})
	}

	const loadQuestionList = () => {
		backend.student
			.studentExamQuestionList({ examAttemptId })
			.then((response) => {
				if (response.success && response.data) {
					questions = response.data.questions
				}
			})
			.catch((err) => {
				catcher(err)
			})
	}

	const submitAnswer = () => {
		if (!selectedQuestion || !submissionAnswer.trim()) {
			toast.error('Please enter an answer')
			return
		}

		submitting = true
		backend.student
			.examSubmit({
				examAttemptId,
				examQuestionId: selectedQuestion.examQuestion.id,
				answer: submissionAnswer.trim(),
			})
			.then((response) => {
				if (response.success) {
					toast.success('Answer submitted successfully')
					submissionAnswer = ''
					loadQuestionList()
					selectQuestion(selectedQuestion!)
				} else {
					toast.error(response.message || 'Submission failed')
				}
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				submitting = false
			})
	}

	onMount(() => {
		loadExamData()
	})
</script>

<div class="m-auto max-w-screen-xl pt-16 pr-6 pl-2 max-lg:px-4 max-lg:pt-12 max-md:px-2">
	{#if loading}
		<div class="flex min-h-[400px] items-center justify-center">
			<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
		</div>
	{:else}
		<div class="flex h-[calc(100vh-4rem)]">
			<div class="relative flex w-96 flex-col border-r">
				<div class="flex-shrink-0 p-4">
					{#if examData && classData}
						<div class="my-4 flex flex-col gap-4">
							<button
								class="text-muted-foreground hover:text-primary flex items-center gap-2 hover:cursor-pointer"
								onclick={() => navigate(`/student/exam/${examAttemptId}/detail`)}
							>
								<ArrowLeftIcon size={16} />
								<span class="text-xs font-medium tracking-wide uppercase">EXAM DETAIL</span>
							</button>
							<div class="flex items-center justify-between">
								<PageTitle
									title={examData.name}
									description={`${classData.code} - ${classData.name}`}
								/>
							</div>
						</div>
					{/if}
					<h2 class="text-lg font-semibold">Questions</h2>
				</div>
				<div class="flex-1 overflow-y-auto px-4 pb-4">
					<div class="space-y-2">
						{#each questions as question, index (question.examQuestion.id)}
							<button
								class="w-full rounded-lg border p-3 text-left transition-colors hover:bg-gray-50 {selectedQuestion
									?.examQuestion.id === question.examQuestion.id
									? 'border-primary/50 bg-primary/5'
									: 'border-gray-200 bg-white'}"
								onclick={() => selectQuestion(question)}
							>
								<div class="flex items-center justify-between">
									<div class="flex-1">
										<div class="text-muted-foreground text-xs">#{index + 1}</div>
										<div class="line-clamp-1 text-sm">
											{question.examQuestion.title}
										</div>
									</div>
									{#if question.status === PayloadExamQuestionWithStatusStatusEnum.Passed}
										<CheckCircle2Icon class="ml-2 h-5 w-5 text-green-500" />
									{:else if question.status === PayloadExamQuestionWithStatusStatusEnum.Rejected}
										<CircleSlashIcon class="ml-2 h-5 w-5 text-orange-500" />
									{:else if question.status === PayloadExamQuestionWithStatusStatusEnum.Invalid}
										<CircleXIcon class="ml-2 h-5 w-5 text-red-500" />
									{:else if question.status === PayloadExamQuestionWithStatusStatusEnum.Unsubmitted}
										<CircleIcon class="ml-2 h-5 w-5 text-gray-400" />
									{/if}
								</div>
							</button>
						{/each}
					</div>
				</div>
				<ExamTimer closedAt={examData?.closedAt} />
			</div>

			<div class="relative flex-1 overflow-y-auto">
				<div class="mx-auto max-w-4xl p-8">
					{#if loadingDetail}
						<div class="flex min-h-[400px] items-center justify-center">
							<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
						</div>
					{:else if selectedQuestion}
						<div class="mb-6">
							<h2 class="mb-2 text-xl font-bold">{selectedQuestion.examQuestion.title}</h2>
							<div class="mb-4 text-sm text-gray-600">
								Question {questions.findIndex(
									(q) => q.examQuestion.id === selectedQuestion?.examQuestion.id
								) + 1} of {questions.length}
							</div>
						</div>

						<div class="mb-8 rounded-lg border bg-gray-50 p-6">
							<div class="prose max-w-none">
								{@html marked.parse(selectedQuestion.examQuestion.description || '')}
							</div>
						</div>

						<div class="space-y-4">
							<div>
								<Label for="answer">Your Answer</Label>
								<Textarea
									id="answer"
									bind:value={submissionAnswer}
									placeholder="Enter your SQL query here..."
									class="mt-2 min-h-[200px] font-mono"
									disabled={submitting}
								/>
							</div>

							<div class="flex justify-end">
								<Button
									onclick={submitAnswer}
									disabled={submitting || !submissionAnswer.trim()}
									size="lg"
								>
									{#if submitting}
										<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
										Submitting...
									{:else}
										Submit
									{/if}
								</Button>
							</div>
						</div>
						<SubmissionList {questionDetail} />
					{:else if !selectedQuestion}
						<div class="flex h-full items-center justify-center text-gray-500">
							Select a question to begin
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style lang="scss">
	.prose {
		:global(h1), :global(h2), :global(h3), :global(h4) {
			font-weight: bold;
			margin-top: 1em;
			margin-bottom: .5em;
		}

		:global(h1) {
			font-size: 1.5em;
		}

		:global(h2) {
			font-size: 1.3em;
		}

		:global(h3) {
			font-size: 1.1em;
		}

		:global(p) {
			margin-bottom: .5em;
		}

		:global(ul), :global(ol) {
			list-style: initial;
			padding-left: 2em;
			margin-bottom: 1em;
		}

		:global(li) {
			margin-bottom: .5em;
		}

		:global(pre) {
			background-color: #f5f5f5;
			padding: 0.75em;
			border-radius: 0.25em;
			overflow-x: auto;
			margin-bottom: 1em;
		}

		:global(code) {
			font-family: monospace;
			background-color: #f5f5f5;
			padding: 0.2em 0.4em;
			border-radius: 0.2em;
		}

		:global(table) {
			width: 100%;
			border-collapse: collapse;
			margin-bottom: 1em;
			overflow-x: auto;
			display: block;
		}

		:global(th) {
			background-color: #f5f5f5;
			font-weight: bold;
			text-align: left;
			padding: 0.5em 0.75em;
			border: 1px solid #e2e8f0;
		}

		:global(td) {
			padding: 0.5em 0.75em;
			border: 1px solid #e2e8f0;
			vertical-align: top;
		}

		:global(tr:nth-child(even)) {
			background-color: #fdfdfd;
		}

		:global(tbody tr:hover) {
			background-color: #f1f5f9;
		}

		:global(caption) {
			font-style: italic;
			padding: 0.5em 0;
			caption-side: bottom;
		}

		:global(hr) {
			margin-bottom: 1em;
		}
	}
</style>
