<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import {
		CalendarIcon,
		ClockIcon,
		FileTextIcon,
		LockIcon,
		UnlockIcon,
		CheckCircle2Icon,
		PlayIcon,
		XCircleIcon,
	} from 'lucide-svelte'
	import type { PayloadClassExamListItem } from '$/util/backend/backend.ts'
	import { PayloadClassExamListItemStatusEnum } from '$/util/backend/backend.ts'

	export let exam: PayloadClassExamListItem

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
		})
	}

	const getStatusDisplay = () => {
		switch (exam.status) {
			case PayloadClassExamListItemStatusEnum.Upcoming:
				return { label: 'Upcoming', variant: 'secondary' as const, icon: ClockIcon }
			case PayloadClassExamListItemStatusEnum.Opened:
				return { label: 'Open', variant: 'default' as const, icon: UnlockIcon }
			case PayloadClassExamListItemStatusEnum.Attempted:
				return { label: 'In Progress', variant: 'outline' as const, icon: PlayIcon }
			case PayloadClassExamListItemStatusEnum.Finished:
				return { label: 'Completed', variant: 'success' as const, icon: CheckCircle2Icon }
			case PayloadClassExamListItemStatusEnum.Closed:
				return { label: 'Closed', variant: 'destructive' as const, icon: XCircleIcon }
			default:
				return { label: 'Unknown', variant: 'secondary' as const, icon: LockIcon }
		}
	}

	$: statusDisplay = getStatusDisplay()
	$: isClickable =
		exam.status === PayloadClassExamListItemStatusEnum.Opened ||
		exam.status === PayloadClassExamListItemStatusEnum.Attempted
</script>

<Card class="transition-all {isClickable ? 'cursor-pointer hover:shadow-md' : 'opacity-75'}">
	<CardHeader>
		<div class="flex items-start justify-between">
			<CardTitle class="text-lg">{exam.exam.name}</CardTitle>
			<Badge variant={statusDisplay.variant}>
				<svelte:component this={statusDisplay.icon} size={12} class="mr-1" />
				{statusDisplay.label}
			</Badge>
		</div>
	</CardHeader>
	<CardContent class="space-y-3">
		<div class="text-muted-foreground flex items-center gap-2 text-sm">
			<FileTextIcon class="h-4 w-4" />
			<span>{exam.questionCount} questions</span>
		</div>

		<div class="space-y-2">
			<div class="text-muted-foreground flex items-center gap-2 text-sm">
				<CalendarIcon class="h-4 w-4" />
				<span>Opens: {formatDate(exam.exam.openedAt)}</span>
			</div>
			<div class="text-muted-foreground flex items-center gap-2 text-sm">
				<ClockIcon class="h-4 w-4" />
				<span>Closes: {formatDate(exam.exam.closedAt)}</span>
			</div>
		</div>

		{#if exam.exam.accessCode && exam.status === PayloadClassExamListItemStatusEnum.Opened}
			<div class="border-t pt-2">
				<p class="text-muted-foreground text-xs">Access Code Required</p>
			</div>
		{/if}

		{#if exam.examAttempt && exam.status === PayloadClassExamListItemStatusEnum.Attempted}
			<div class="border-t pt-2">
				<p class="text-xs text-blue-600">Started: {formatDate(exam.examAttempt.openedAt)}</p>
			</div>
		{/if}

		{#if exam.examAttempt && exam.status === PayloadClassExamListItemStatusEnum.Finished}
			<div class="border-t pt-2">
				<p class="text-xs text-green-600">Finished: {formatDate(exam.examAttempt.finishedAt)}</p>
			</div>
		{/if}
	</CardContent>
</Card>
