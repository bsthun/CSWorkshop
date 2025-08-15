<script lang="ts">
	import { onMount } from 'svelte'
	import { Link } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Avatar, AvatarImage, AvatarFallback } from '$/lib/shadcn/components/ui/avatar'
	import {
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow
	} from '$/lib/shadcn/components/ui/table'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import { UsersIcon, Loader2Icon, EyeIcon, CheckCircleIcon, XCircleIcon, ClockIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadExamJoineeListItem } from '$/util/backend/backend.ts'
	import { formatDateTime } from '$/util/format.ts'

	export let examId: number

	let joinees: PayloadExamJoineeListItem[] = []
	let loading = true

	const loadJoinees = () => {
		loading = true
		backend.admin
			.examJoineeList({ examId })
			.then((response) => {
				if (response.success && response.data) {
					joinees = response.data.joinees || []
				}
			})
			.catch(catcher)
			.finally(() => {
				loading = false
			})
	}

	const getStatusBadge = (joinee: PayloadExamJoineeListItem) => {
		if (joinee.finishedAt) {
			return { text: 'Finished', variant: 'default', icon: CheckCircleIcon, color: 'text-green-600' }
		} else if (joinee.startedAt) {
			return { text: 'In Progress', variant: 'secondary', icon: ClockIcon, color: 'text-yellow-600' }
		} else if (joinee.openedAt) {
			return { text: 'Opened', variant: 'outline', icon: EyeIcon, color: 'text-blue-600' }
		} else {
			return { text: 'Not Started', variant: 'outline', icon: XCircleIcon, color: 'text-gray-600' }
		}
	}

	const getScoreColor = (score: number, total: number) => {
		const percentage = total > 0 ? (score / total) * 100 : 0
		if (percentage >= 80) return 'text-green-600'
		if (percentage >= 60) return 'text-yellow-600'
		return 'text-red-600'
	}

	const getTotalScore = (scores: any) => {
		return scores.passed + scores.rejected + scores.invalid + scores.unsubmitted
	}

	const getInitials = (firstname: string, lastname: string) => {
		return `${firstname.charAt(0)}${lastname.charAt(0)}`.toUpperCase()
	}

	onMount(() => {
		loadJoinees()
	})
</script>

<Card class="h-auto">
	<CardHeader class="pb-3">
		<CardTitle class="flex items-center gap-2">
			<UsersIcon class="h-5 w-5" />
			Students & Attempts ({joinees.length})
		</CardTitle>
	</CardHeader>
	<CardContent>
		{#if loading}
			<div class="flex items-center justify-center py-8">
				<Loader2Icon class="text-muted-foreground h-8 w-8 animate-spin" />
			</div>
		{:else if joinees.length === 0}
			<div class="flex flex-col items-center justify-center py-12">
				<UsersIcon class="mb-4 h-16 w-16 text-gray-400" />
				<h3 class="mb-2 text-lg font-semibold">No Student Attempts</h3>
				<p class="text-muted-foreground text-center">
					No students have started this exam yet.
				</p>
			</div>
		{:else}
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Student</TableHead>
						<TableHead>Status</TableHead>
						<TableHead>Score</TableHead>
						<TableHead>Started</TableHead>
						<TableHead>Finished</TableHead>
						<TableHead class="text-right">Actions</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#each joinees as joinee}
						{@const status = getStatusBadge(joinee)}
						{@const totalScore = getTotalScore(joinee.score)}
						<TableRow class="hover:bg-gray-50">
							<TableCell>
								<div class="flex items-center gap-3">
									<Avatar class="h-8 w-8">
										<AvatarImage src={joinee.joinee.user.pictureUrl} alt="" />
										<AvatarFallback class="text-xs">
											{getInitials(joinee.joinee.user.firstname, joinee.joinee.user.lastname)}
										</AvatarFallback>
									</Avatar>
									<div>
										<div class="font-medium">
											{joinee.joinee.user.firstname} {joinee.joinee.user.lastname}
										</div>
										<div class="text-sm text-gray-500">{joinee.joinee.user.email}</div>
									</div>
								</div>
							</TableCell>
							<TableCell>
								<Badge variant={status.variant} class="flex items-center gap-1 w-fit">
									<svelte:component this={status.icon} class="h-3 w-3 {status.color}" />
									{status.text}
								</Badge>
							</TableCell>
							<TableCell>
								{#if totalScore > 0}
									<div class="flex flex-col">
										<span class="font-medium {getScoreColor(joinee.score.passed, totalScore)}">
											{joinee.score.passed}/{totalScore}
										</span>
										<div class="flex gap-1 text-xs text-gray-500">
											<span class="text-green-600">✓{joinee.score.passed}</span>
											<span class="text-red-600">✗{joinee.score.rejected}</span>
											<span class="text-yellow-600">!{joinee.score.invalid}</span>
											<span class="text-gray-400">-{joinee.score.unsubmitted}</span>
										</div>
									</div>
								{:else}
									<span class="text-gray-400">No submissions</span>
								{/if}
							</TableCell>
							<TableCell>
								{#if joinee.startedAt}
									<span class="text-sm">{formatDateTime(joinee.startedAt)}</span>
								{:else}
									<span class="text-gray-400 text-sm">Not started</span>
								{/if}
							</TableCell>
							<TableCell>
								{#if joinee.finishedAt}
									<span class="text-sm">{formatDateTime(joinee.finishedAt)}</span>
								{:else}
									<span class="text-gray-400 text-sm">In progress</span>
								{/if}
							</TableCell>
							<TableCell class="text-right">
								{#if joinee.startedAt}
									<Link to="/admin/attempt?examAttemptId={joinee.id}">
										<Button variant="ghost" size="sm">
											<EyeIcon class="h-4 w-4 mr-2" />
											View Details
										</Button>
									</Link>
								{:else}
									<Button variant="ghost" size="sm" disabled>
										<EyeIcon class="h-4 w-4 mr-2" />
										No Attempt
									</Button>
								{/if}
							</TableCell>
						</TableRow>
					{/each}
				</TableBody>
			</Table>
		{/if}
	</CardContent>
</Card>