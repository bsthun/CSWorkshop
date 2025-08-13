<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Loader2Icon, Trash2Icon } from 'lucide-svelte'
	import type { PayloadCollectionQuestionDetail } from '$/util/backend/backend.ts'

	export let open = false
	export let question: PayloadCollectionQuestionDetail | null = null
	export let deleting = false

	const dispatch = createEventDispatcher<{
		delete: void
		close: void
	}>()

	const handleDelete = () => {
		dispatch('delete')
	}

	const handleClose = () => {
		dispatch('close')
	}
</script>

<Dialog bind:open>
	<DialogContent>
		<DialogHeader>
			<DialogTitle class="flex gap-2">
				<Trash2Icon size={16} />
				Delete Question
			</DialogTitle>
			<DialogDescription>
				Delete this question? This action cannot be undone.
				{#if question}
					<div class="mt-2 rounded bg-gray-50 p-2 text-sm">
						<strong>{question.title || `Question ${question.id}`}</strong>
					</div>
				{/if}
			</DialogDescription>
		</DialogHeader>
		<DialogFooter>
			<Button variant="outline" onclick={handleClose} disabled={deleting}>Cancel</Button>
			<Button variant="destructive" onclick={handleDelete} disabled={deleting}>Delete Question</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
