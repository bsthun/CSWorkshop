<script lang="ts">
	import { onDestroy } from 'svelte'
	import { navigate } from 'svelte-navigator'
	import ExamTimeupDialog from './ExamTimeupDialog.svelte'
	import type { PayloadExamQuestionWithStatus } from '$/util/backend/backend.ts'

	export let closedAt: string | null = null
	export let questions: PayloadExamQuestionWithStatus[] = []

	let timeRemaining = { hours: 0, minutes: 0, seconds: 0 }
	let timerInterval: ReturnType<typeof setInterval> | null = null
	let timeUpDialogOpen = false
	let timeUp = false

	const updateTimer = () => {
		if (!closedAt) return

		const now = new Date().getTime()
		const closeTime = new Date(closedAt).getTime()
		const diff = closeTime - now

		if (diff <= 0) {
			timeRemaining = { hours: 0, minutes: 0, seconds: 0 }
			if (timerInterval) {
				clearInterval(timerInterval)
				timerInterval = null
			}
			if (!timeUpDialogOpen) {
				timeUp = true
				timeUpDialogOpen = true
			}
			return
		}

		const hours = Math.floor(diff / (1000 * 60 * 60))
		const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
		const seconds = Math.floor((diff % (1000 * 60)) / 1000)

		timeRemaining = { hours, minutes, seconds }
	}

	const startTimer = () => {
		if (timerInterval) clearInterval(timerInterval)
		updateTimer()
		timerInterval = setInterval(updateTimer, 1000)
	}

	const formatDate = (dateString: string) => {
		const date = new Date(dateString)
		return date.toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
		})
	}

	$: if (closedAt) {
		startTimer()
	}

	$: if (timeUp && !timeUpDialogOpen) {
		navigate('/student')
	}

	onDestroy(() => {
		if (timerInterval) {
			clearInterval(timerInterval)
		}
	})
</script>

{#if closedAt}
	<div class="absolute bottom-4 left-1/2 z-10 -translate-x-1/2">
		<div class="rounded-lg border border-gray-200 bg-white px-4 py-2 shadow-lg">
			<div class="flex items-center gap-2">
				<div class="text-center">
					<div class="font-mono text-lg font-medium text-gray-700">
						{String(timeRemaining.hours).padStart(2, '0')}:{String(timeRemaining.minutes).padStart(
							2,
							'0'
						)}:<span class="text-sm">{String(timeRemaining.seconds).padStart(2, '0')}</span>
					</div>
					<div class="text-xs text-gray-500">
						Closes {formatDate(closedAt)}
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<ExamTimeupDialog bind:open={timeUpDialogOpen} {questions} />
