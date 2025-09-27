<script lang="ts">
	import { CheckCircle2Icon, CircleSlashIcon, CircleXIcon, CircleIcon } from 'lucide-svelte'
	import { formatDateTime } from '$/util/format'

	export let questionDetail: any
</script>

{#if questionDetail?.submissions && questionDetail.submissions.length > 0}
	<div class="mb-6">
		<h3 class="mb-3 text-lg font-semibold">Submissions</h3>
		<div class="space-y-2">
			{#each questionDetail.submissions as submission, idx}
				<div
					class="rounded-lg border p-4 {submission.status === 'passed'
						? 'border-green-200 bg-green-50'
						: submission.status === 'rejected'
							? 'border-orange-200 bg-orange-50'
							: submission.status === 'invalid'
								? 'border-red-200 bg-red-50'
								: 'border-gray-200 bg-gray-50'}"
				>
					<!-- First row: Status icon and text on left, #number and timestamp on right -->
					<div class="mb-3 flex items-center justify-between">
						<div class="flex items-center gap-2">
							{#if submission.status === 'passed'}
								<CheckCircle2Icon class="h-6 w-6 text-green-500" />
							{:else if submission.status === 'rejected'}
								<CircleSlashIcon class="h-6 w-6 text-orange-500" />
							{:else if submission.status === 'invalid'}
								<CircleXIcon class="h-6 w-6 text-red-500" />
							{:else}
								<CircleIcon class="h-6 w-6 text-gray-400" />
							{/if}
							<span
								class="text-sm font-medium uppercase {submission.status === 'passed'
									? 'text-green-600'
									: submission.status === 'rejected'
										? 'text-orange-600'
										: submission.status === 'invalid'
											? 'text-red-600'
											: 'text-gray-600'}"
							>
								{submission.status}
							</span>
						</div>
						<div class="text-right">
							<span class="text-xs font-semibold text-gray-500">
								#{questionDetail.submissions.length - idx}
							</span>
							{#if submission.examSubmission?.createdAt}
								<span class="text-xs text-gray-500">
									• {formatDateTime(submission.examSubmission.createdAt)}
								</span>
							{/if}
						</div>
					</div>

					<div class="flex gap-4">
						<div class="flex-1 space-y-3">
							{#if submission.examSubmission?.checkQueryPassed !== undefined}
								<div>
									<div class="ml-1 flex items-center gap-3">
										{#if submission.examSubmission.checkQueryPassed}
											<CheckCircle2Icon class="h-4 w-4 text-green-500" />
										{:else}
											<CircleXIcon class="h-4 w-4 text-red-500" />
										{/if}
										<span class="text-sm text-gray-700">Query Result Checking</span>
									</div>
									{#if submission.examSubmission.checkQueryAt}
										<div class="ml-8 text-xs text-gray-500">
											Checked {formatDateTime(submission.examSubmission.checkQueryAt)}
										</div>
									{/if}
									{#if submission.examSubmission.result?.executionError && !submission.examSubmission.checkQueryPassed}
										<div class="mt-1 ml-8 rounded border border-red-200 bg-red-50 p-2 overflow-x-auto">
											<pre class="font-mono text-xs whitespace-pre-wrap text-red-700 break-words">{submission
													.examSubmission.result.executionError}</pre>
										</div>
									{/if}
								</div>
							{/if}
							{#if submission.examSubmission?.checkPromptPassed !== undefined}
								<div>
									<div class="ml-1 flex items-center gap-3">
										{#if submission.examSubmission.checkPromptPassed}
											<CheckCircle2Icon class="h-4 w-4 text-green-500" />
										{:else}
											<CircleXIcon class="h-4 w-4 text-red-500" />
										{/if}
										<span class="text-sm text-gray-700">Prompt Checking</span>
									</div>
									{#if submission.examSubmission.checkPromptAt}
										<div class="ml-8 text-xs text-gray-500">
											Checked {formatDateTime(submission.examSubmission.checkPromptAt)}
										</div>
									{/if}
									{#if submission.examSubmission.result?.promptDescription}
										<div
											class="mt-1 ml-8 rounded border p-2 overflow-x-auto {submission.examSubmission
												.checkPromptPassed
												? 'border-green-200 bg-green-50'
												: 'border-red-200 bg-red-50'}"
										>
											<p
												class="text-xs break-words {submission.examSubmission.checkPromptPassed
													? 'text-green-700'
													: 'text-red-700'}"
											>
												{submission.examSubmission.result.promptDescription}
											</p>
										</div>
									{/if}
									{#if submission.examSubmission.result?.promptError}
										<div class="mt-1 ml-8 rounded bg-gray-50 p-1.5 overflow-x-auto">
											<pre
												class="font-mono text-xs whitespace-pre-wrap text-red-700 break-words">{submission
													.examSubmission.result.promptError}</pre>
										</div>
									{/if}
								</div>
							{/if}
						</div>

						<!-- Right side: Code block -->
						<div class="flex-1">
							<div class="rounded border bg-gray-50 p-3 overflow-x-auto">
								<pre class="font-mono text-xs text-gray-700 whitespace-pre-wrap break-words"><code
										>{submission.examSubmission?.answer || ''}</code
									></pre>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}
