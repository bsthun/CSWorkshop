<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogHeader,
		DialogTitle
	} from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Card, CardContent } from '$/lib/shadcn/components/ui/card'
	import { BookOpenIcon, Loader2Icon, SearchIcon, PlusIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type {
		PayloadCollectionQuestionListItem,
		PayloadCollection
	} from '$/util/backend/backend.ts'
	import { toast } from 'svelte-sonner'

	export let open = false
	export let examId: number
	export let collectionData: PayloadCollection

	const dispatch = createEventDispatcher<{
		added: void
	}>()

	let questions: PayloadCollectionQuestionListItem[] = []
	let filteredQuestions: PayloadCollectionQuestionListItem[] = []
	let loading = true
	let adding = false
	let searchQuery = ''

	const loadQuestions = () => {
		loading = true
		backend.admin
			.collectionQuestionList({ collectionId: collectionData.id })
			.then((response) => {
				if (response.success && response.data) {
					questions = response.data.questions || []
					filteredQuestions = questions
				}
			})
			.catch(catcher)
			.finally(() => {
				loading = false
			})
	}

	const handleSearch = () => {
		if (!searchQuery.trim()) {
			filteredQuestions = questions
		} else {
			const query = searchQuery.toLowerCase()
			filteredQuestions = questions.filter(
				(question) =>
					question.title?.toLowerCase().includes(query) ||
					question.description?.toLowerCase().includes(query)
			)
		}
	}

	const handleAddQuestion = (questionId: number) => {
		adding = true
		backend.admin
			.examQuestionAdd({
				examId,
				collectionQuestionId: questionId
			})
			.then((response) => {
				if (response.success) {
					toast.success('Question added to exam')
					dispatch('added')
					open = false
				}
			})
			.catch(catcher)
			.finally(() => {
				adding = false
			})
	}

	const handleKeydown = (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			handleSearch()
		}
	}

	// Watch for search query changes
	$: {
		handleSearch()
	}

	// Load questions when dialog opens
	$: if (open) {
		loadQuestions()
		searchQuery = ''
	}
</script>

<Dialog bind:open>
	<DialogContent class="max-w-2xl">
		<DialogHeader>
			<DialogTitle>Add Question to Exam</DialogTitle>
			<DialogDescription>
				Select questions from the "{collectionData.name}" collection to add to your exam.
			</DialogDescription>
		</DialogHeader>

		<div class="space-y-4">
			<!-- Search -->
			<div class="space-y-2">
				<Label for="search">Search Questions</Label>
				<div class="relative">
					<SearchIcon class="absolute left-3 top-3 h-4 w-4 text-gray-400" />
					<Input
						id="search"
						bind:value={searchQuery}
						placeholder="Search by title or description..."
						class="pl-10"
						on:keydown={handleKeydown}
					/>
				</div>
			</div>

			<!-- Questions List -->
			<div class="max-h-96 space-y-2 overflow-y-auto">
				{#if loading}
					<div class="flex items-center justify-center py-8">
						<Loader2Icon class="text-muted-foreground h-6 w-6 animate-spin" />
					</div>
				{:else if filteredQuestions.length === 0}
					<div class="flex flex-col items-center justify-center py-8">
						<BookOpenIcon class="mb-4 h-12 w-12 text-gray-400" />
						<p class="text-muted-foreground text-center">
							{searchQuery ? 'No questions match your search' : 'No questions available in this collection'}
						</p>
					</div>
				{:else}
					{#each filteredQuestions as question}
						<Card class="transition-shadow hover:shadow-md">
							<CardContent class="p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<h4 class="font-medium">
											{question.title || `Question ${question.orderNum || question.id}`}
										</h4>
										{#if question.description}
											<p class="mt-1 text-sm text-gray-600 line-clamp-2">
												{question.description}
											</p>
										{/if}
										<div class="mt-2 text-xs text-gray-500">
											ID: {question.id} • Order: {question.orderNum || 'N/A'}
										</div>
									</div>
									<Button
										size="sm"
										onclick={() => handleAddQuestion(question.id)}
										disabled={adding}
									>
										{#if adding}
											<Loader2Icon class="h-4 w-4 animate-spin" />
										{:else}
											<PlusIcon class="h-4 w-4 mr-1" />
											Add
										{/if}
									</Button>
								</div>
							</CardContent>
						</Card>
					{/each}
				{/if}
			</div>
		</div>

		<div class="flex justify-end gap-2">
			<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
		</div>
	</DialogContent>
</Dialog>