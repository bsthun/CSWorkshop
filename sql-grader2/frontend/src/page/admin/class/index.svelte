<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { GraduationCapIcon, PlusIcon, Loader2Icon, UsersIcon, HashIcon, EditIcon, PencilIcon } from 'lucide-svelte'
	import { Link } from 'svelte-navigator'
	import type { PayloadSemester } from '$/util/backend/backend.ts'

	export let semesters: PayloadSemester[] = []
	export let loading = false

	const dispatch = createEventDispatcher<{
		createClass: void
		createSemester: void
		editSemester: { semester: PayloadSemester }
	}>()

	const handleSemesterEdit = (semester: PayloadSemester) => {
		dispatch('editSemester', { semester })
	}
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
			<Button variant="outline" class="gap-2" onclick={() => dispatch('createSemester')}>
				<PlusIcon class="h-4 w-4" />
				Create Semester
			</Button>
			<Button class="gap-2" onclick={() => dispatch('createClass')}>
				<PlusIcon class="h-4 w-4" />
				Create Class
			</Button>
		</div>
	</div>
{:else}
	<div class="space-y-8">
		{#each semesters as semester}
			<div class="space-y-4">
				<div class="flex items-center justify-start gap-2">
					<h2 class="text-xl font-semibold">{semester.name}</h2>
					<Button variant="ghost" size="sm" onclick={() => handleSemesterEdit(semester)}>
						<PencilIcon />
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
											<span
												class="text-muted-foreground text-xs font-medium tracking-wide uppercase"
											>
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
