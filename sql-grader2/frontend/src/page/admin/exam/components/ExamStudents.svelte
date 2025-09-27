<script lang="ts">
	import { onMount } from 'svelte'
	import { Link } from 'svelte-navigator'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '$/lib/shadcn/components/ui/table'
	import UserProfile from '$/component/share/UserProfile.svelte'
	import {
		UsersIcon,
		Loader2Icon,
		CheckCircle2Icon,
		CircleSlashIcon,
		CircleXIcon,
		CircleIcon,
		EyeIcon,
	} from 'lucide-svelte'
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

	const getTotalScore = (scores: any) => {
		return scores.passed + scores.rejected + scores.invalid + scores.unsubmitted
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
				<p class="text-muted-foreground text-center">No students have started this exam yet.</p>
			</div>
		{:else}
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Student</TableHead>
						<TableHead>Score</TableHead>
						<TableHead>Started</TableHead>
						<TableHead>Updated</TableHead>
						<TableHead class="text-right">Actions</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#each joinees as joinee}
						{@const totalScore = getTotalScore(joinee.score)}
						<TableRow class="hover:bg-gray-50">
							<TableCell>
								<UserProfile user={joinee.joinee.user} />
							</TableCell>
							<TableCell>
								{#if totalScore > 0}
									<div class="flex items-center gap-3">
										<div class="flex items-center gap-1">
											<CheckCircle2Icon class="h-4 w-4 text-green-500" />
											<span class="text-sm font-medium text-green-600">{joinee.score.passed}</span
											>
										</div>
										<div class="flex items-center gap-1">
											<CircleSlashIcon class="h-4 w-4 text-orange-500" />
											<span class="text-sm font-medium text-orange-600"
												>{joinee.score.rejected}</span
											>
										</div>
										<div class="flex items-center gap-1">
											<CircleXIcon class="h-4 w-4 text-red-500" />
											<span class="text-sm font-medium text-red-600">{joinee.score.invalid}</span>
										</div>
										<div class="flex items-center gap-1">
											<CircleIcon class="h-4 w-4 text-gray-400" />
											<span class="text-sm font-medium text-gray-500"
												>{joinee.score.unsubmitted}</span
											>
										</div>
									</div>
								{:else}
									<span class="text-gray-400">No submissions</span>
								{/if}
							</TableCell>
							<TableCell>
								<div class="flex flex-col">
									<span class="text-xs text-gray-500">Started at</span>
									<span class="text-sm">{formatDateTime(joinee.startedAt)}</span>
								</div>
							</TableCell>
							<TableCell>
								<div class="flex flex-col">
									<span class="text-xs text-gray-500">
										{#if joinee.finishedAt}
											Finished at
										{:else}
											Last action at
										{/if}
									</span>
									<span class="text-sm">
										{#if joinee.finishedAt}
											{formatDateTime(joinee.finishedAt)}
										{:else}
											{formatDateTime(joinee.updatedAt)}
										{/if}
									</span>
								</div>
							</TableCell>
							<TableCell class="text-right">
								<Link to="/admin/attempt?examAttemptId={joinee.id}">
									<Button variant="ghost" size="sm">
										<EyeIcon class="mr-2 h-4 w-4" />
										View Details
									</Button>
								</Link>
							</TableCell>
						</TableRow>
					{/each}
				</TableBody>
			</Table>
		{/if}
	</CardContent>
</Card>
