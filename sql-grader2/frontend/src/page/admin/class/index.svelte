<script lang="ts">
	import { onMount } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { GraduationCapIcon, PlusIcon, Loader2Icon, UsersIcon, HashIcon, EditIcon } from 'lucide-svelte'
	import { Link } from 'svelte-navigator'
	import CreateClassDialog from './dialog/CreateClassDialog.svelte'
	import EditSemesterDialog from './dialog/EditSemesterDialog.svelte'
	import CreateSemesterDialog from './dialog/CreateSemesterDialog.svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import type { PayloadSemester } from '$/util/backend/backend.ts'

	let semesters: PayloadSemester[] = []
	let loading = true
	let showCreateDialog = false
	let showEditSemesterDialog = false
	let showCreateSemesterDialog = false
	let editingSemester: PayloadSemester | null = null

	const loadSemesters = () => {
		loading = true
		backend.admin
			.semesterList({})
			.then((response) => {
				semesters = response.data.semesters!
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleClassCreated = () => {
		loadSemesters()
	}

	const handleSemesterEdit = (semester: PayloadSemester) => {
		editingSemester = semester
		showEditSemesterDialog = true
	}

	const handleSemesterEdited = () => {
		loadSemesters()
		editingSemester = null
	}

	const handleSemesterCreated = () => {
		loadSemesters()
	}

	onMount(loadSemesters)
</script>

{#if loading}
	<div class="flex min-h-[400px] items-center justify-center">
		<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
	</div>
{:else if semesters.length === 0}
	<div class="flex min-h-[400px] flex-col items-center justify-center">
		<GraduationCapIcon class="mb-4 h-16 w-16 text-gray-400" />
		<h3 class="mb-2 text-lg font-semibold">No classes yet</h3>
		<p class="text-muted-foreground mb-4">Create your first semester and class to get started</p>
		<div class="flex gap-2">
			<Button variant="outline" class="gap-2" onclick={() => (showCreateSemesterDialog = true)}>
				<PlusIcon class="h-4 w-4" />
				Create Semester
			</Button>
			<Button class="gap-2" onclick={() => (showCreateDialog = true)}>
				<PlusIcon class="h-4 w-4" />
				Create Class
			</Button>
		</div>
	</div>
{:else}
	<div class="mb-6 flex justify-end gap-2">
		<Button variant="outline" class="gap-2" onclick={() => (showCreateSemesterDialog = true)}>
			<PlusIcon class="h-4 w-4" />
			Create Semester
		</Button>
		<Button class="gap-2" onclick={() => (showCreateDialog = true)}>
			<PlusIcon class="h-4 w-4" />
			Create Class
		</Button>
	</div>

	<div class="space-y-8">
		{#each semesters as semester}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h2 class="text-xl font-semibold">{semester.name}</h2>
					<Button
						variant="ghost"
						size="sm"
						onclick={() => handleSemesterEdit(semester)}
						class="gap-2"
					>
						<EditIcon class="h-4 w-4" />
						Edit
					</Button>
				</div>

				{#if semester.classes.length === 0}
					<p class="text-muted-foreground text-sm">No classes in this semester</p>
				{:else}
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each semester.classes as classItem}
							<Link to="/admin/class/{classItem.id}">
								<Card class="transition-shadow hover:shadow-lg">
									<CardHeader class="gap-2">
										<CardTitle class="flex items-center gap-1">
											<GraduationCapIcon size={12} class="text-muted-foreground" />
											<span class="text-muted-foreground text-xs font-medium tracking-wide uppercase">
												CLASS
											</span>
										</CardTitle>
										<h3 class="text-lg font-semibold">{classItem.code}</h3>
										<h4 class="text-sm text-gray-600">{classItem.name}</h4>
									</CardHeader>
									<CardContent>
										<div class="space-y-2">
											<div class="text-muted-foreground flex items-center gap-2 text-sm">
												<UsersIcon class="h-4 w-4" />
												<span>{classItem.joineeCount} students</span>
											</div>
											<div class="text-muted-foreground flex items-center gap-2 text-sm">
												<HashIcon class="h-4 w-4" />
												<span>{classItem.id}</span>
											</div>
										</div>
									</CardContent>
								</Card>
							</Link>
						{/each}
					</div>
				{/if}
			</div>
		{/each}
	</div>
{/if}

<CreateClassDialog bind:open={showCreateDialog} on:created={handleClassCreated} />
<CreateSemesterDialog bind:open={showCreateSemesterDialog} on:created={handleSemesterCreated} />
{#if editingSemester}
	<EditSemesterDialog
		bind:open={showEditSemesterDialog}
		semester={editingSemester}
		on:edited={handleSemesterEdited}
	/>
{/if}