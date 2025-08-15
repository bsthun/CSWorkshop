<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle
	} from '$/lib/shadcn/components/ui/dialog'
	import { Card, CardContent, CardHeader, CardTitle } from '$/lib/shadcn/components/ui/card'
	import { Badge } from '$/lib/shadcn/components/ui/badge'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Textarea } from '$/lib/shadcn/components/ui/textarea'
	import { 
		CheckCircleIcon, 
		XCircleIcon, 
		AlertTriangleIcon, 
		MinusCircleIcon,
		ClockIcon,
		FileTextIcon,
		HelpCircleIcon,
		DatabaseIcon,
		Loader2Icon
	} from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'
	import { formatDateTime } from '$/util/format.ts'
	import type { 
		PayloadSubmissionListItem,
		PayloadSubmissionDetailResponse 
	} from '$/util/backend/backend.ts'

	export let open = false
	export let submission: PayloadSubmissionListItem | null = null

	let detailData: PayloadSubmissionDetailResponse | null = null
	let loading = false

	const loadSubmissionDetail = (submissionId: number) => {
		loading = true
		backend.admin
			.submissionDetail({ submissionId })
			.then((response) => {
				if (response.success && response.data) {
					detailData = response.data
				}
			})
			.catch(catcher)
			.finally(() => {
				loading = false
			})
	}

	const getStatusInfo = (sub: any) => {
		if (sub.checkPromptPassed && sub.checkQueryPassed) {
			return { 
				icon: CheckCircleIcon, 
				color: 'text-green-600', 
				bg: 'bg-green-100',
				text: 'Passed',
				variant: 'default' as const
			}
		} else if (sub.checkPromptPassed === false || sub.checkQueryPassed === false) {
			return { 
				icon: XCircleIcon, 
				color: 'text-red-600', 
				bg: 'bg-red-100',
				text: 'Failed',
				variant: 'destructive' as const
			}
		} else if (sub.answer) {
			return { 
				icon: AlertTriangleIcon, 
				color: 'text-yellow-600', 
				bg: 'bg-yellow-100',
				text: 'Submitted',
				variant: 'secondary' as const
			}
		} else {
			return { 
				icon: MinusCircleIcon, 
				color: 'text-gray-400', 
				bg: 'bg-gray-100',
				text: 'Not Submitted',
				variant: 'outline' as const
			}
		}
	}

	// Watch for submission changes
	$: if (open && submission) {
		loadSubmissionDetail(submission.submission.id)
	}

	// Reset when dialog closes
	$: if (!open) {
		detailData = null
	}
</script>

<Dialog bind:open>
	<DialogContent class="max-w-4xl max-h-[90vh] overflow-y-auto">
		{#if loading}
			<div class="flex items-center justify-center py-8">
				<Loader2Icon class="text-muted-foreground h-8 w-8 animate-spin" />
			</div>
		{:else if submission && detailData}
			{@const statusInfo = getStatusInfo(detailData.submission)}
			<DialogHeader>
				<DialogTitle class="flex items-center gap-2">
					<FileTextIcon class="h-5 w-5" />
					Submission Details
				</DialogTitle>
			</DialogHeader>

			<div class="space-y-6">
				<!-- Status Overview -->
				<Card>
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center justify-between">
							<span class="flex items-center gap-2">
								<svelte:component this={statusInfo.icon} class="h-5 w-5 {statusInfo.color}" />
								Submission Status
							</span>
							<Badge variant={statusInfo.variant}>
								{statusInfo.text}
							</Badge>
						</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
							<div>
								<Label class="text-muted-foreground text-xs">Submitted</Label>
								<p>{detailData.submission.createdAt ? formatDateTime(detailData.submission.createdAt) : 'Not submitted'}</p>
							</div>
							<div>
								<Label class="text-muted-foreground text-xs">Prompt Check</Label>
								<p class="flex items-center gap-1">
									{#if detailData.submission.checkPromptPassed === true}
										<CheckCircleIcon class="h-3 w-3 text-green-600" /> Passed
									{:else if detailData.submission.checkPromptPassed === false}
										<XCircleIcon class="h-3 w-3 text-red-600" /> Failed
									{:else}
										<MinusCircleIcon class="h-3 w-3 text-gray-400" /> Not checked
									{/if}
								</p>
							</div>
							<div>
								<Label class="text-muted-foreground text-xs">Query Check</Label>
								<p class="flex items-center gap-1">
									{#if detailData.submission.checkQueryPassed === true}
										<CheckCircleIcon class="h-3 w-3 text-green-600" /> Passed
									{:else if detailData.submission.checkQueryPassed === false}
										<XCircleIcon class="h-3 w-3 text-red-600" /> Failed
									{:else}
										<MinusCircleIcon class="h-3 w-3 text-gray-400" /> Not checked
									{/if}
								</p>
							</div>
							<div>
								<Label class="text-muted-foreground text-xs">Last Updated</Label>
								<p>{formatDateTime(detailData.submission.updatedAt)}</p>
							</div>
						</div>
					</CardContent>
				</Card>

				<!-- Question Details -->
				<Card>
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<HelpCircleIcon class="h-5 w-5 text-blue-600" />
							Question
						</CardTitle>
					</CardHeader>
					<CardContent class="space-y-4">
						<div>
							<Label class="text-sm font-medium">Title</Label>
							<p class="text-sm mt-1">{detailData.question.title || 'Untitled Question'}</p>
						</div>
						
						{#if detailData.question.description}
							<div>
								<Label class="text-sm font-medium">Description</Label>
								<p class="text-sm mt-1 text-gray-600">{detailData.question.description}</p>
							</div>
						{/if}

						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div>
								<Label class="text-sm font-medium">Check Prompt</Label>
								<Textarea 
									value={detailData.question.checkPrompt || 'No check prompt'} 
									readonly
									rows="4"
									class="mt-1 text-xs font-mono bg-gray-50"
								/>
							</div>
							<div>
								<Label class="text-sm font-medium">Check Query</Label>
								<Textarea 
									value={detailData.question.checkQuery || 'No check query'} 
									readonly
									rows="4"
									class="mt-1 text-xs font-mono bg-gray-50"
								/>
							</div>
						</div>
					</CardContent>
				</Card>

				<!-- Student Answer -->
				<Card>
					<CardHeader class="pb-3">
						<CardTitle class="flex items-center gap-2">
							<DatabaseIcon class="h-5 w-5 text-purple-600" />
							Student Answer
						</CardTitle>
					</CardHeader>
					<CardContent>
						{#if detailData.submission.answer}
							<div>
								<Label class="text-sm font-medium">SQL Query</Label>
								<Textarea 
									value={detailData.submission.answer} 
									readonly
									rows="8"
									class="mt-1 font-mono text-sm bg-gray-50"
								/>
							</div>
						{:else}
							<div class="flex items-center justify-center py-8 text-gray-400">
								<div class="text-center">
									<MinusCircleIcon class="h-12 w-12 mx-auto mb-2" />
									<p>No answer submitted</p>
								</div>
							</div>
						{/if}
					</CardContent>
				</Card>

				<!-- Check Results Timeline -->
				{#if detailData.submission.checkPromptAt || detailData.submission.checkQueryAt}
					<Card>
						<CardHeader class="pb-3">
							<CardTitle class="flex items-center gap-2">
								<ClockIcon class="h-5 w-5 text-orange-600" />
								Check Timeline
							</CardTitle>
						</CardHeader>
						<CardContent>
							<div class="space-y-3">
								{#if detailData.submission.checkPromptAt}
									<div class="flex items-center gap-3 p-3 rounded-lg bg-gray-50">
										<div class="flex items-center justify-center w-8 h-8 rounded-full {detailData.submission.checkPromptPassed ? 'bg-green-100' : 'bg-red-100'}">
											{#if detailData.submission.checkPromptPassed}
												<CheckCircleIcon class="h-4 w-4 text-green-600" />
											{:else}
												<XCircleIcon class="h-4 w-4 text-red-600" />
											{/if}
										</div>
										<div class="flex-1">
											<p class="font-medium text-sm">Prompt Check</p>
											<p class="text-xs text-gray-500">{formatDateTime(detailData.submission.checkPromptAt)}</p>
										</div>
										<Badge variant={detailData.submission.checkPromptPassed ? 'default' : 'destructive'}>
											{detailData.submission.checkPromptPassed ? 'Passed' : 'Failed'}
										</Badge>
									</div>
								{/if}

								{#if detailData.submission.checkQueryAt}
									<div class="flex items-center gap-3 p-3 rounded-lg bg-gray-50">
										<div class="flex items-center justify-center w-8 h-8 rounded-full {detailData.submission.checkQueryPassed ? 'bg-green-100' : 'bg-red-100'}">
											{#if detailData.submission.checkQueryPassed}
												<CheckCircleIcon class="h-4 w-4 text-green-600" />
											{:else}
												<XCircleIcon class="h-4 w-4 text-red-600" />
											{/if}
										</div>
										<div class="flex-1">
											<p class="font-medium text-sm">Query Check</p>
											<p class="text-xs text-gray-500">{formatDateTime(detailData.submission.checkQueryAt)}</p>
										</div>
										<Badge variant={detailData.submission.checkQueryPassed ? 'default' : 'destructive'}>
											{detailData.submission.checkQueryPassed ? 'Passed' : 'Failed'}
										</Badge>
									</div>
								{/if}
							</div>
						</CardContent>
					</Card>
				{/if}
			</div>
		{:else}
			<div class="flex items-center justify-center py-8">
				<p class="text-gray-500">No submission data available</p>
			</div>
		{/if}
	</DialogContent>
</Dialog>