<script lang="ts">
	import { Loader2Icon, TableIcon, CheckIcon, ColumnsIcon, SquareIcon, PlusIcon, MinusIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend'
	import type {
		PayloadExamQuestionSubmissionCompare,
		PayloadExamQuestionSubmissionCompareTable,
	} from '$/util/backend/backend.ts'
	import { Drawer, DrawerContent, DrawerHeader, DrawerTitle } from '$/lib/shadcn/components/ui/drawer'
	import { ScrollArea } from '$/lib/shadcn/components/ui/scroll-area'
	import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '$/lib/shadcn/components/ui/table'
	import { Badge } from '$/lib/shadcn/components/ui/badge'

	export let open = false
	export let submissionId: number | null = null
	export let questionName: string = ''

	let comparison: PayloadExamQuestionSubmissionCompare | null = null
	let loading = false

	$: if (open && submissionId) {
		loadComparison()
	}

	const loadComparison = () => {
		if (!submissionId) return

		loading = true
		comparison = null

		backend.student
			.examSubmissionCompare({ submissionId })
			.then((response) => {
				comparison = response.data.comparison
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const areRowsEqual = (row1: any[], row2: any[]): boolean => {
		if (row1.length !== row2.length) return false
		return row1.every((val, idx) => val === row2[idx])
	}

	const findMatchingRow = (row: any[], table: PayloadExamQuestionSubmissionCompareTable): boolean => {
		return table.rows.some((r) => areRowsEqual(r, row))
	}

	const isCellCorrect = (
		cellValue: any,
		cellIndex: number,
		currentHeader: string[],
		otherTable: PayloadExamQuestionSubmissionCompareTable
	): boolean => {
		// Get the column name for this cell
		const columnName = currentHeader[cellIndex]

		// Find the corresponding column index in the other table
		const otherColumnIndex = otherTable.header.findIndex((header) => header === columnName)

		// If column doesn't exist in other table, cell is incorrect
		if (otherColumnIndex === -1) return false

		// Check if this cell value exists in the same column of any row in the other table
		return otherTable.rows.some((otherRow) => otherRow[otherColumnIndex] === cellValue)
	}
</script>

<Drawer bind:open direction="right">
	<DrawerContent class="min-w-5xl">
		<DrawerHeader>
			<DrawerTitle>
				<div class="flex items-center gap-2">
					<TableIcon class="h-5 w-5" />
					Result Comparison
				</div>
			</DrawerTitle>
		</DrawerHeader>

		<div class="flex-1 overflow-hidden px-4 pb-4">
			{#if loading}
				<div class="flex min-h-96 items-center justify-center">
					<Loader2Icon class="text-primary h-8 w-8 animate-spin" />
				</div>
			{:else if comparison}
				{@const submittedCount = comparison.submitted.rows.length}
				{@const expectedCount = comparison.expected.rows.length}
				{@const matchingRows = comparison.submitted.rows.filter((row) =>
					findMatchingRow(row, comparison.expected)
				).length}
				{@const extraRows = submittedCount - matchingRows}
				{@const missingRows = expectedCount - matchingRows}

				<ScrollArea class="h-full">
					<div class="space-y-6">
						<!-- Submitted Result -->
						<div>
							<div class="mb-3 flex items-center gap-2">
								<h3 class="text-lg font-semibold">Your Result</h3>
								<Badge variant="secondary">
									{comparison.submitted.rows.length} row{comparison.submitted.rows.length !== 1
										? 's'
										: ''}
								</Badge>
							</div>
							<div class="rounded-lg border">
								<Table>
									<TableHeader>
										<TableRow>
											{#each comparison.submitted.header as header}
												{@const headerExists = comparison.expected.header.includes(header)}
												<TableHead
													class="font-mono text-sm font-bold {!headerExists
														? 'text-red-600'
														: ''}">{header}</TableHead
												>
											{/each}
										</TableRow>
									</TableHeader>
									<TableBody>
										{#if comparison.submitted.rows.length === 0}
											<TableRow>
												<TableCell
													colspan={comparison.submitted.header.length}
													class="text-muted-foreground text-center"
												>
													No rows returned
												</TableCell>
											</TableRow>
										{:else}
											{#each comparison.submitted.rows as row}
												<TableRow>
													{#each row as cell, cellIndex}
														{@const cellCorrect = isCellCorrect(
															cell,
															cellIndex,
															comparison.submitted.header,
															comparison.expected
														)}
														<TableCell
															class="font-mono text-sm {!cellCorrect
																? 'text-red-600'
																: ''}"
														>
															{cell ?? 'NULL'}
														</TableCell>
													{/each}
												</TableRow>
											{/each}
										{/if}
									</TableBody>
								</Table>
							</div>
						</div>

						<!-- Expected Result -->
						<div>
							<div class="mb-3 flex items-center gap-2">
								<h3 class="text-lg font-semibold">Expected Result</h3>
								<Badge variant="secondary">
									{comparison.expected.rows.length} row{comparison.expected.rows.length !== 1
										? 's'
										: ''}
								</Badge>
							</div>
							<div class="rounded-lg border">
								<Table>
									<TableHeader>
										<TableRow>
											{#each comparison.expected.header as header}
												{@const headerExists = comparison.submitted.header.includes(header)}
												<TableHead
													class="font-mono text-sm font-bold {!headerExists
														? 'text-red-600'
														: ''}">{header}</TableHead
												>
											{/each}
										</TableRow>
									</TableHeader>
									<TableBody>
										{#if comparison.expected.rows.length === 0}
											<TableRow>
												<TableCell
													colspan={comparison.expected.header.length}
													class="text-muted-foreground text-center"
												>
													No rows expected
												</TableCell>
											</TableRow>
										{:else}
											{#each comparison.expected.rows as row}
												<TableRow>
													{#each row as cell, cellIndex}
														{@const cellCorrect = isCellCorrect(
															cell,
															cellIndex,
															comparison.expected.header,
															comparison.submitted
														)}
														<TableCell
															class="font-mono text-sm {!cellCorrect
																? 'text-red-600'
																: ''}"
														>
															{cell ?? 'NULL'}
														</TableCell>
													{/each}
												</TableRow>
											{/each}
										{/if}
									</TableBody>
								</Table>
							</div>
						</div>

						<div class="rounded-lg border border-blue-200 bg-blue-50 p-4">
							<h4 class="mb-3 font-semibold text-blue-900">Comparison Summary</h4>

							<div class="space-y-3">
								{#if true}
									{@const matchingHeaders = comparison.submitted.header.filter((h) =>
										comparison.expected.header.includes(h)
									).length}
									{@const extraHeaders = comparison.submitted.header.length - matchingHeaders}
									{@const missingHeaders = comparison.expected.header.length - matchingHeaders}
									{@const isHeaderPerfect =
										matchingHeaders === comparison.submitted.header.length &&
										matchingHeaders === comparison.expected.header.length}

									<div class="flex items-center gap-3">
										<ColumnsIcon class="h-4 w-4 text-blue-600" />
										<span class="text-sm font-medium">Matching Headers</span>
										<Badge variant="outline" class="{isHeaderPerfect ? 'border-green-500 bg-green-50 text-green-700' : 'border-blue-500 bg-blue-50 text-blue-700'}">
											<div class="flex items-center gap-0.5">
												<CheckIcon class="h-3 w-3" />
												{matchingHeaders}
												{#if extraHeaders > 0}
													<PlusIcon class="ml-2 h-3 w-3" />
													{extraHeaders}
												{/if}
												{#if missingHeaders > 0}
													<MinusIcon class="ml-2 h-3 w-3" />
													{missingHeaders}
												{/if}
											</div>
										</Badge>
									</div>
								{/if}

								<!-- Row Matching -->
								{#if true}
									{@const isRowPerfect =
										matchingRows === submittedCount && matchingRows === expectedCount}

									<div class="flex items-center gap-3">
										<TableIcon class="h-4 w-4 text-blue-600" />
										<span class="text-sm font-medium">Matching Rows</span>
										<Badge variant="outline" class="{isRowPerfect ? 'border-green-500 bg-green-50 text-green-700' : 'border-blue-500 bg-blue-50 text-blue-700'}">
											<div class="flex items-center gap-0.5">
												<CheckIcon class="h-3 w-3" />
												{matchingRows}
												{#if extraRows > 0}
													<PlusIcon class="ml-2 h-3 w-3" />
													{extraRows}
												{/if}
												{#if missingRows > 0}
													<MinusIcon class="ml-2 h-3 w-3" />
													{missingRows}
												{/if}
											</div>
										</Badge>
									</div>
								{/if}

								{#if true}
									{@const totalSubmittedCells =
										submittedCount * (comparison.submitted.header.length || 0)}
									{@const totalExpectedCells =
										expectedCount * (comparison.expected.header.length || 0)}
									{@const matchingCells = comparison.submitted.rows.reduce((acc, row) => {
										return (
											acc +
											row.filter((cell, idx) =>
												isCellCorrect(
													cell,
													idx,
													comparison.submitted.header,
													comparison.expected
												)
											).length
										)
									}, 0)}
									{@const extraCells = totalSubmittedCells - matchingCells}
									{@const missingCells =
										totalExpectedCells -
										comparison.expected.rows.reduce((acc, row) => {
											return (
												acc +
												row.filter((cell, idx) =>
													isCellCorrect(
														cell,
														idx,
														comparison.expected.header,
														comparison.submitted
													)
												).length
											)
										}, 0)}
									{@const isCellPerfect =
										matchingCells === totalSubmittedCells && matchingCells === totalExpectedCells}

									<div class="flex items-center gap-3">
										<SquareIcon class="h-4 w-4 text-blue-600" />
										<span class="text-sm font-medium">Matching Cells</span>
										<Badge variant="outline" class="{isCellPerfect ? 'border-green-500 bg-green-50 text-green-700' : 'border-blue-500 bg-blue-50 text-blue-700'}">
											<div class="flex items-center gap-0.5">
												<CheckIcon class="h-3 w-3" />
												{matchingCells}
												{#if extraCells > 0}
													<PlusIcon class="ml-2 h-3 w-3" />
													{extraCells}
												{/if}
												{#if missingCells > 0}
													<MinusIcon class="ml-2 h-3 w-3" />
													{missingCells}
												{/if}
											</div>
										</Badge>
									</div>
								{/if}
							</div>
						</div>
					</div>
				</ScrollArea>
			{:else if !loading}
				<div class="flex min-h-96 flex-col items-center justify-center">
					<TableIcon class="mb-4 h-16 w-16 text-gray-400" />
					<h3 class="mb-2 text-lg font-semibold">No comparison data available</h3>
					<p class="text-muted-foreground text-center text-sm">Unable to load comparison data.</p>
				</div>
			{/if}
		</div>
	</DrawerContent>
</Drawer>
