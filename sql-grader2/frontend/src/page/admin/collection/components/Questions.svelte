<script lang="ts">
	import { onMount } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import DeleteQuestionDialog from '../dialog/DeleteQuestionDialog.svelte'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Textarea } from '$/lib/shadcn/components/ui/textarea'
	import { BookOpenIcon, Loader2Icon, SaveIcon, PlusIcon, Trash2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadCollectionQuestionListItem, PayloadCollectionQuestionDetail } from '$/util/backend/backend.ts'
	import { toast } from 'svelte-sonner'
	import { Carta, MarkdownEditor } from 'carta-md'

	export let collectionId: number

	let questions: PayloadCollectionQuestionListItem[] = []
	let selectedQuestion: PayloadCollectionQuestionDetail | null = null
	let questionsLoading = true
	let questionDetailLoading = false
	let saving = false
	let deleteDialogOpen = false
	let deleting = false

	// Carta instance
	const carta = new Carta({
		sanitizer: false,
	})

	// Form fields
	let editedTitle = ''
	let editedDescription = ''
	let editedCheckPrompt = ''
	let editedCheckQuery = ''
	let hasChanges = false

	const loadQuestions = () => {
		questionsLoading = true
		backend.admin
			.collectionQuestionList({ collectionId })
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

	const loadQuestionDetail = (questionId: number) => {
		questionDetailLoading = true
		backend.admin
			.collectionQuestionDetail({ questionId })
			.then((response) => {
				if (response.success && response.data) {
					selectedQuestion = response.data
					// Initialize form fields
					editedTitle = response.data.title || ''
					editedDescription = response.data.description || ''
					editedCheckPrompt = response.data.checkPrompt || ''
					editedCheckQuery = response.data.checkQuery || ''
					hasChanges = false
				}
			})
			.catch(catcher)
			.finally(() => {
				questionDetailLoading = false
			})
	}

	const handleQuestionSelect = (question: PayloadCollectionQuestionListItem) => {
		if (selectedQuestion?.id !== question.id) {
			loadQuestionDetail(question.id)
		}
	}

	const handleSave = () => {
		if (!selectedQuestion || !hasChanges) return

		saving = true
		backend.admin
			.collectionQuestionEdit({
				id: selectedQuestion.id,
				title: editedTitle,
				description: editedDescription,
				checkPrompt: editedCheckPrompt,
				checkQuery: editedCheckQuery,
			})
			.then((response) => {
				if (response.success) {
					hasChanges = false
					toast.success('Question updated successfully')
					// Reload the questions list to update titles
					loadQuestions()
					// Update the selected question data
					if (selectedQuestion) {
						selectedQuestion.title = editedTitle
						selectedQuestion.description = editedDescription
						selectedQuestion.checkPrompt = editedCheckPrompt
						selectedQuestion.checkQuery = editedCheckQuery
					}
				}
			})
			.catch(catcher)
			.finally(() => {
				saving = false
			})
	}

	const createNewQuestion = () => {
		backend.admin
			.collectionQuestionCreate({
				collectionId,
				description: 'New question description',
			})
			.then((response) => {
				if (response.success) {
					toast.success('Question created successfully')
					loadQuestions()
				}
			})
			.catch(catcher)
	}

	const handleDelete = () => {
		if (!selectedQuestion) return

		deleting = true
		backend.admin
			.collectionQuestionDelete({ id: selectedQuestion.id })
			.then((response) => {
				if (response.success) {
					toast.success('Question deleted successfully')
					deleteDialogOpen = false
					selectedQuestion = null
					loadQuestions()
				}
			})
			.catch(catcher)
			.finally(() => {
				deleting = false
			})
	}

	const handleDialogClose = () => {
		deleteDialogOpen = false
	}

	$: {
		if (selectedQuestion) {
			hasChanges =
				editedTitle !== (selectedQuestion.title || '') ||
				editedDescription !== (selectedQuestion.description || '') ||
				editedCheckPrompt !== (selectedQuestion.checkPrompt || '') ||
				editedCheckQuery !== (selectedQuestion.checkQuery || '')
		}
	}

	onMount(() => {
		loadQuestions()
	})
</script>

<Card class="h-auto">
	<CardHeader class="pb-3">
		<CardTitle class="flex items-center gap-2">
			<BookOpenIcon class="h-5 w-5" />
			Questions
		</CardTitle>
	</CardHeader>
	<CardContent class="pb-0">
		<div class="flex">
			<div class="flex h-full max-h-[640px] w-1/3 flex-col">
				<div class="mb-4 flex items-center justify-between">
					<h3 class="text-sm font-medium">Questions</h3>
					<Button onclick={createNewQuestion} size="sm">
						<PlusIcon class="h-4 w-4" />
					</Button>
				</div>
				<div class="overflow-y-auto">
					{#if questionsLoading}
						<div class="flex items-center justify-center py-8">
							<Loader2Icon class="text-muted-foreground h-6 w-6 animate-spin" />
						</div>
					{:else if questions.length === 0}
						<div class="flex flex-col items-center justify-center py-8">
							<BookOpenIcon class="mb-4 h-16 w-16 text-gray-400" />
							<p class="text-muted-foreground text-center">No questions yet</p>
							<Button size="sm" class="mt-4" onclick={createNewQuestion}>
								<PlusIcon class="mr-2 h-4 w-4" />
								Create Question
							</Button>
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
										class="mb-1 text-sm font-medium {selectedQuestion?.id === question.id
											? 'text-blue-600'
											: ''}"
									>
										{question.title || `Question ${question.orderNum || question.id}`}
									</div>
									<div
										class="truncate text-xs {selectedQuestion?.id === question.id
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
				{#if selectedQuestion}
					<div class="mb-4 flex items-center justify-between">
						<h3 class="text-sm font-medium">Question Details</h3>
						<div class="flex items-center gap-2">
							{#if hasChanges}
								<Button onclick={handleSave} disabled={saving}>
									<SaveIcon size={20} />
									Save
								</Button>
							{/if}
							<Button variant="destructive" onclick={() => (deleteDialogOpen = true)}>
								<Trash2Icon size={20} />
							</Button>
						</div>
					</div>
				{/if}
				<div>
					{#if questionDetailLoading}
						<div class="flex min-h-96 items-center justify-center py-8">
							<Loader2Icon class="text-muted-foreground h-6 w-6 animate-spin" />
						</div>
					{:else if !selectedQuestion}
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

							<!-- Description Field -->
							<div class="space-y-2">
								<Label for="description">Description</Label>
								<div class="border rounded-md">
									<MarkdownEditor
										{carta}
										bind:value={editedDescription}
										placeholder="Question description..."
										mode="tabs"
									/>
								</div>
							</div>

							<!-- Check Prompt Field -->
							<div class="space-y-2">
								<Label for="checkPrompt">Check Prompt</Label>
								<Textarea
									id="checkPrompt"
									bind:value={editedCheckPrompt}
									placeholder="Check prompt for validation..."
									rows="4"
									class="font-mono text-sm"
								/>
							</div>

							<!-- Check Query Field -->
							<div class="space-y-2">
								<Label for="checkQuery">Check Query</Label>
								<Textarea
									id="checkQuery"
									bind:value={editedCheckQuery}
									placeholder="SQL query for validation..."
									rows="4"
									class="font-mono text-sm"
								/>
							</div>

							<!-- Metadata -->
							<div class="grid grid-cols-2 gap-4 border-t pt-4">
								<div>
									<Label class="text-muted-foreground text-xs">Question ID</Label>
									<p class="text-sm">{selectedQuestion.id}</p>
								</div>
								<div>
									<Label class="text-muted-foreground text-xs">Order Number</Label>
									<p class="text-sm">{selectedQuestion.orderNum || 'N/A'}</p>
								</div>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</CardContent>
</Card>

<DeleteQuestionDialog
	bind:open={deleteDialogOpen}
	{deleting}
	on:close={handleDialogClose}
	on:delete={handleDelete}
	question={selectedQuestion}
/>
