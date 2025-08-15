<script lang="ts">
	import { onMount } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Textarea } from '$/lib/shadcn/components/ui/textarea'
	import { BookOpenIcon, Loader2Icon, SaveIcon, PlusIcon, Trash2Icon, RotateCcwIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type {
		PayloadCollectionQuestionListItem,
		PayloadExamQuestionDetailResponse,
		PayloadCollection,
	} from '$/util/backend/backend.ts'
	import { toast } from 'svelte-sonner'
	import AddQuestionDialog from '../dialog/AddQuestionDialog.svelte'

	export let examId: number
	export let collectionData: PayloadCollection

	let questions: PayloadCollectionQuestionListItem[] = []
	let selectedQuestionDetail: PayloadExamQuestionDetailResponse | null = null
	let questionsLoading = true
	let questionDetailLoading = false
	let saving = false
	let addDialogOpen = false
	let deleting = false

	// Form fields
	let editedTitle = ''
	let editedDescription = ''
	let editedCheckPrompt = ''
	let editedCheckQuery = ''
	let hasChanges = false

	const loadQuestions = () => {
		questionsLoading = true
		backend.admin
			.examQuestionList({ examId })
			.then((response) => {
				if (response.success && response.data) {
					questions = response.data.questions || []
				}
			})
			.catch(catcher)
			.finally(() => {
				questionsLoading = false
			})
	}

	const loadQuestionDetail = (examQuestionId: number) => {
		questionDetailLoading = true
		backend.admin
			.examQuestionDetail({ examQuestionId })
			.then((response) => {
				if (response.success && response.data) {
					selectedQuestionDetail = response.data
					// Initialize form fields from exam question
					editedTitle = response.data.examQuestion.title || ''
					editedDescription = response.data.examQuestion.description || ''
					editedCheckPrompt = response.data.examQuestion.checkPrompt || ''
					editedCheckQuery = response.data.examQuestion.checkQuery || ''
					hasChanges = false
				}
			})
			.catch(catcher)
			.finally(() => {
				questionDetailLoading = false
			})
	}

	const handleQuestionSelect = (question: PayloadCollectionQuestionListItem) => {
		if (selectedQuestionDetail?.examQuestion.id !== question.id) {
			loadQuestionDetail(question.id)
		}
	}

	const handleSave = () => {
		if (!selectedQuestionDetail || !hasChanges) return

		saving = true
		backend.admin
			.examQuestionEdit({
				examQuestionId: selectedQuestionDetail.examQuestion.id,
				title: editedTitle,
				description: editedDescription,
				checkPrompt: editedCheckPrompt,
				checkQuery: editedCheckQuery,
			})
			.then((response) => {
				if (response.success && response.data) {
					hasChanges = false
					toast.success('Question updated successfully')
					// Reload the questions list to update titles
					loadQuestions()
					// Update the selected question data
					if (selectedQuestionDetail) {
						selectedQuestionDetail.examQuestion = response.data
					}
				}
			})
			.catch(catcher)
			.finally(() => {
				saving = false
			})
	}

	const handleReset = () => {
		if (!selectedQuestionDetail) return

		// Reset to collection question values
		editedTitle = selectedQuestionDetail.collectionQuestion.title || ''
		editedDescription = selectedQuestionDetail.collectionQuestion.description || ''
		editedCheckPrompt = selectedQuestionDetail.collectionQuestion.checkPrompt || ''
		editedCheckQuery = selectedQuestionDetail.collectionQuestion.checkQuery || ''
	}

	const handleDelete = () => {
		if (!selectedQuestionDetail) return

		deleting = true
		backend.admin
			.examQuestionDelete({ examQuestionId: selectedQuestionDetail.examQuestion.id })
			.then((response) => {
				if (response.success) {
					toast.success('Question removed from exam')
					selectedQuestionDetail = null
					loadQuestions()
				}
			})
			.catch(catcher)
			.finally(() => {
				deleting = false
			})
	}

	const handleQuestionAdded = () => {
		loadQuestions()
	}

	// Track changes
	$: {
		if (selectedQuestionDetail) {
			hasChanges =
				editedTitle !== (selectedQuestionDetail.examQuestion.title || '') ||
				editedDescription !== (selectedQuestionDetail.examQuestion.description || '') ||
				editedCheckPrompt !== (selectedQuestionDetail.examQuestion.checkPrompt || '') ||
				editedCheckQuery !== (selectedQuestionDetail.examQuestion.checkQuery || '')
		}
	}

	onMount(() => {
		loadQuestions()
	})
</script>

<Card class="h-auto">
	<CardHeader class="pb-3">
		<div class="flex items-center justify-between">
			<CardTitle class="flex items-center gap-2">
				<BookOpenIcon class="h-5 w-5" />
				Exam Questions
			</CardTitle>
			<Button onclick={() => (addDialogOpen = true)} size="sm">
				<PlusIcon class="mr-2 h-4 w-4" />
				Add Question
			</Button>
		</div>
	</CardHeader>
	<CardContent class="pb-0">
		<div class="flex">
			<div class="flex w-1/3 flex-col">
				<div class="mb-4 flex items-center justify-between">
					<h3 class="text-sm font-medium">Questions ({questions.length})</h3>
				</div>
				<div class="overflow-y-auto">
					{#if questionsLoading}
						<div class="flex min-h-80 items-center justify-center py-8">
							<Loader2Icon class="text-muted-foreground h-6 w-6 animate-spin" />
						</div>
					{:else if questions.length === 0}
						<div class="flex min-h-80 flex-col items-center justify-center py-8">
							<BookOpenIcon class="mb-4 h-16 w-16 text-gray-400" />
							<p class="text-muted-foreground mb-4 text-center">No questions in this exam</p>
						</div>
					{:else}
						<div class="divide-y divide-gray-200">
							{#each questions as question, index}
								<button
									class="w-full p-3 text-left transition-colors hover:bg-gray-50 {index <
									questions.length - 1
										? 'border-b border-gray-200'
										: ''}"
									onclick={() => handleQuestionSelect(question)}
								>
									<div
										class="mb-1 text-sm font-medium {selectedQuestionDetail?.examQuestion.id ===
										question.id
											? 'text-blue-600'
											: ''}"
									>
										{question.title || `Question ${question.orderNum || question.id}`}
									</div>
									<div
										class="truncate text-xs {selectedQuestionDetail?.examQuestion.id === question.id
											? 'text-blue-500'
											: 'text-muted-foreground'}"
									>
										{question.description || 'No description'}
									</div>
								</button>
							{/each}
						</div>
					{/if}
				</div>
			</div>

			<div class="mx-4 w-px bg-gray-200"></div>

			<div class="flex flex-1 flex-col">
				{#if selectedQuestionDetail}
					<div class="mb-4 flex items-center justify-between">
						<h3 class="text-sm font-medium">Question Details</h3>
						<div class="flex items-center gap-2">
							<Button variant="outline" onclick={handleReset} size="sm">
								<RotateCcwIcon class="mr-2 h-4 w-4" />
								Reset
							</Button>
							{#if hasChanges}
								<Button onclick={handleSave} disabled={saving} size="sm">
									<SaveIcon class="mr-2 h-4 w-4" />
									Save
								</Button>
							{/if}
							<Button variant="destructive" onclick={handleDelete} disabled={deleting} size="sm">
								<Trash2Icon class="h-4 w-4" />
							</Button>
						</div>
					</div>
				{/if}
				<div>
					{#if questionDetailLoading}
						<div class="flex min-h-96 items-center justify-center py-8">
							<Loader2Icon class="text-muted-foreground h-6 w-6 animate-spin" />
						</div>
					{:else if !selectedQuestionDetail}
						<div class="flex min-h-96 flex-col items-center justify-center py-8">
							<BookOpenIcon class="mb-4 h-16 w-16 text-gray-400" />
							<p class="text-muted-foreground text-center">
								Select a question from the list to view and edit details
							</p>
						</div>
					{:else}
						<div class="space-y-6">
							<div class="space-y-2">
								<Label for="title">Title</Label>
								<Input id="title" bind:value={editedTitle} placeholder="Question title..." />
							</div>

							<div class="space-y-2">
								<Label for="description">Description</Label>
								<Textarea
									id="description"
									bind:value={editedDescription}
									placeholder="Question description..."
									rows={3}
								/>
							</div>

							<div class="space-y-2">
								<Label for="checkPrompt">Check Prompt</Label>
								<Textarea
									id="checkPrompt"
									bind:value={editedCheckPrompt}
									placeholder="Check prompt for validation..."
									rows={4}
								/>
							</div>

							<div class="space-y-2">
								<Label for="checkQuery">Check Query</Label>
								<Textarea
									id="checkQuery"
									bind:value={editedCheckQuery}
									placeholder="SQL query for validation..."
									rows={4}
									class="font-mono text-sm"
								/>
							</div>

							<!-- Metadata -->
							<div class="grid grid-cols-2 gap-4 border-t pt-4">
								<div>
									<Label class="text-muted-foreground text-xs">Exam Question ID</Label>
									<p class="text-sm">{selectedQuestionDetail.examQuestion.id}</p>
								</div>
								<div>
									<Label class="text-muted-foreground text-xs">Original Question ID</Label>
									<p class="text-sm">{selectedQuestionDetail.examQuestion.originalQuestionId}</p>
								</div>
								<div>
									<Label class="text-muted-foreground text-xs">Order Number</Label>
									<p class="text-sm">{selectedQuestionDetail.examQuestion.orderNum || 'N/A'}</p>
								</div>
								<div>
									<Label class="text-muted-foreground text-xs">Collection Source</Label>
									<p class="text-sm font-medium text-blue-600">{collectionData.name}</p>
								</div>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</CardContent>
</Card>

<AddQuestionDialog bind:open={addDialogOpen} {examId} {collectionData} on:added={handleQuestionAdded} />
