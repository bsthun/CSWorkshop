<script lang="ts">
	import { navigate } from 'svelte-navigator'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import {
		ClockIcon,
		CheckCircle2Icon,
		CircleSlashIcon,
		CircleXIcon,
		CircleIcon,
	} from 'lucide-svelte'
	import type { PayloadExamQuestionWithStatus } from '$/util/backend/backend.ts'
	import { PayloadExamQuestionWithStatusStatusEnum } from '$/util/backend/backend.ts'

	export let open = false
	export let questions: PayloadExamQuestionWithStatus[] = []

	const handleHome = () => {
		navigate('/student')
	}

	$: score = {
		passed: questions.filter((q) => q.status === PayloadExamQuestionWithStatusStatusEnum.Passed)
			.length,
		rejected: questions.filter(
			(q) => q.status === PayloadExamQuestionWithStatusStatusEnum.Rejected
		).length,
		invalid: questions.filter(
			(q) => q.status === PayloadExamQuestionWithStatusStatusEnum.Invalid
		).length,
		unsubmitted: questions.filter(
			(q) => q.status === PayloadExamQuestionWithStatusStatusEnum.Unsubmitted
		).length,
	}

	$: if (open) {
		setTimeout(() => {
			handleHome()
		}, 12000)
	}
</script>

<Dialog bind:open>
	<DialogContent class="sm:max-w-md">
		<DialogHeader>
			<div class="mb-1 flex items-center justify-center">
				<ClockIcon class="h-12 w-12 text-red-500" />
			</div>
			<DialogTitle class="text-center text-xl">Time's Up!</DialogTitle>
			<DialogDescription class="text-center">
				The examination time has ended. <br />You will be redirected to the home page in shortly.
			</DialogDescription>
		</DialogHeader>
		<div class="my-4 flex items-center justify-center gap-3">
			<div class="flex items-center gap-1">
				<CheckCircle2Icon class="h-4 w-4 text-green-500" />
				<span class="text-sm font-medium text-green-600">{score.passed}</span>
			</div>
			<div class="flex items-center gap-1">
				<CircleSlashIcon class="h-4 w-4 text-orange-500" />
				<span class="text-sm font-medium text-orange-600">{score.rejected}</span>
			</div>
			<div class="flex items-center gap-1">
				<CircleXIcon class="h-4 w-4 text-red-500" />
				<span class="text-sm font-medium text-red-600">{score.invalid}</span>
			</div>
			<div class="flex items-center gap-1">
				<CircleIcon class="h-4 w-4 text-gray-400" />
				<span class="text-sm font-medium text-gray-500">{score.unsubmitted}</span>
			</div>
		</div>
		<DialogFooter class="flex justify-center">
			<Button onclick={handleHome}>Home</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
