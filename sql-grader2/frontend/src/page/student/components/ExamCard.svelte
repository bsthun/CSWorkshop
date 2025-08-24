<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import { 
		CalendarIcon, 
		ClockIcon, 
		FileTextIcon,
		LockIcon,
		UnlockIcon
	} from 'lucide-svelte'
	import type { PayloadClassExamListItem } from '$/util/backend/backend.ts'

	export let exam: PayloadClassExamListItem

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		})
	}

	const getExamStatus = () => {
		const now = new Date()
		const openedAt = new Date(exam.exam.openedAt)
		const closedAt = new Date(exam.exam.closedAt)

		if (now < openedAt) {
			return { status: 'upcoming', label: 'Upcoming', variant: 'secondary' as const }
		} else if (now > closedAt) {
			return { status: 'closed', label: 'Closed', variant: 'destructive' as const }
		} else {
			return { status: 'open', label: 'Open', variant: 'success' as const }
		}
	}

	$: examStatus = getExamStatus()
</script>

<Card class="cursor-pointer transition-all hover:shadow-md hover:scale-[1.02]">
	<CardHeader>
		<div class="flex items-start justify-between">
			<CardTitle class="text-lg">{exam.exam.name}</CardTitle>
			<Badge variant={examStatus.variant}>
				{#if examStatus.status === 'open'}
					<UnlockIcon size={12} class="mr-1" />
				{:else}
					<LockIcon size={12} class="mr-1" />
				{/if}
				{examStatus.label}
			</Badge>
		</div>
	</CardHeader>
	<CardContent class="space-y-3">
		<div class="flex items-center gap-2 text-sm text-muted-foreground">
			<FileTextIcon class="h-4 w-4" />
			<span>{exam.questionCount} questions</span>
		</div>
		
		<div class="space-y-2">
			<div class="flex items-center gap-2 text-sm text-muted-foreground">
				<CalendarIcon class="h-4 w-4" />
				<span>Opens: {formatDate(exam.exam.openedAt)}</span>
			</div>
			<div class="flex items-center gap-2 text-sm text-muted-foreground">
				<ClockIcon class="h-4 w-4" />
				<span>Closes: {formatDate(exam.exam.closedAt)}</span>
			</div>
		</div>

		{#if exam.exam.accessCode}
			<div class="pt-2 border-t">
				<p class="text-xs text-muted-foreground">Access Code Required</p>
			</div>
		{/if}
	</CardContent>
</Card>