<script lang="ts">
	import { onMount } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import { GraduationCapIcon, PlusIcon, Loader2Icon, CalendarIcon, CheckCircleIcon, ClockIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend'
	import type { PayloadStudentClassListItem } from '$/util/backend/backend.ts'
	import JoinClassDialog from './dialog/JoinClassDialog.svelte'
	import ExamListDrawer from './components/ExamListDrawer.svelte'
	import Container from '$/component/layout/Container.svelte'
	import PageTitle from '$/component/ui/PageTitle.svelte'

	let classes: PayloadStudentClassListItem[] = []
	let loading = false
	let joinClassDialogOpen = false
	let examDrawerOpen = false
	let selectedClass: PayloadStudentClassListItem | null = null

	const loadClasses = () => {
		loading = true
		backend.student
			.classList({})
			.then((response) => {
				classes = response.data.classes
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleClassJoined = () => {
		loadClasses()
		joinClassDialogOpen = false
	}

	const handleClassClick = (classItem: PayloadStudentClassListItem) => {
		selectedClass = classItem
		examDrawerOpen = true
	}

	onMount(() => {
		loadClasses()
	})
</script>

<Container>
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<PageTitle
				title="My Classes"
				description="View and manage your enrolled classes"
				icon={GraduationCapIcon}
			/>
			<Button onclick={() => (joinClassDialogOpen = true)} class="gap-2">
				<PlusIcon class="h-4 w-4" />
				Join Class
			</Button>
		</div>

		{#if loading}
			<div class="flex min-h-[400px] items-center justify-center">
				<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
			</div>
		{:else if classes.length === 0}
			<div class="flex min-h-[400px] flex-col items-center justify-center">
				<GraduationCapIcon class="mb-4 h-16 w-16 text-gray-400" />
				<h3 class="mb-2 text-lg font-semibold">No classes enrolled</h3>
				<p class="text-muted-foreground mb-4">Join a class to get started</p>
				<Button onclick={() => (joinClassDialogOpen = true)} class="gap-2">
					<PlusIcon class="h-4 w-4" />
					Join Class
				</Button>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each classes as classItem}
					<button onclick={() => handleClassClick(classItem)}>
						<Card class="cursor-pointer transition-shadow hover:shadow-lg">
							<CardHeader class="flex flex-col items-start">
								<div class="mb-2 flex items-center justify-between">
									<Badge variant="outline" class="gap-1">
										<CalendarIcon size={12} />
										{classItem.semester.name}
									</Badge>
								</div>
								<CardTitle class="text-xl">{classItem.class.code}</CardTitle>
								<p class="text-muted-foreground text-sm">{classItem.class.name}</p>
							</CardHeader>
							<CardContent>
								<div class="space-y-3">
									<div class="flex items-center justify-between">
										<div class="flex items-center gap-2 text-sm">
											<ClockIcon class="text-muted-foreground h-4 w-4" />
											<span class="text-muted-foreground">Total Exams</span>
										</div>
										<span class="font-medium">{classItem.examTotalCount}</span>
									</div>
									<div class="flex items-center justify-between">
										<div class="flex items-center gap-2 text-sm">
											<CheckCircleIcon class="h-4 w-4 text-green-600" />
											<span class="text-muted-foreground">Completed</span>
										</div>
										<span class="font-medium text-green-600">{classItem.examFinishedCount}</span>
									</div>
									{#if classItem.examTotalCount > 0}
										<div class="mt-3">
											<div
												class="text-muted-foreground mb-1 flex items-center justify-between text-xs"
											>
												<span>Progress</span>
												<span
													>{Math.round(
														(classItem.examFinishedCount / classItem.examTotalCount) * 100
													)}%</span
												>
											</div>
											<div class="h-2 w-full rounded-full bg-gray-200">
												<div
													class="bg-primary h-2 rounded-full transition-all"
													style="width: {(classItem.examFinishedCount /
														classItem.examTotalCount) *
														100}%"
												></div>
											</div>
										</div>
									{/if}
								</div>
							</CardContent>
						</Card>
					</button>
				{/each}
			</div>
		{/if}
	</div>
</Container>

<JoinClassDialog bind:open={joinClassDialogOpen} on:joined={handleClassJoined} />

<ExamListDrawer bind:open={examDrawerOpen} {selectedClass} />
