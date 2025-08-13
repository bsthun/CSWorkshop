<script lang="ts">
	import { Dialog, DialogContent, DialogHeader, DialogTitle } from '$/lib/shadcn/components/ui/dialog'
	import { DatabaseIcon, TableIcon } from 'lucide-svelte'

	export let open = false
	export let structure: any[] = []

	const formatNumber = (num: number) => {
		return new Intl.NumberFormat().format(num)
	}
</script>

<Dialog bind:open>
	<DialogContent class="max-w-2xl">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-2">
				<DatabaseIcon class="h-5 w-5" />
				Database Structure
			</DialogTitle>
		</DialogHeader>

		<div class="max-h-96 overflow-y-auto">
			<div class="relative w-full overflow-auto">
				<table class="w-full caption-bottom text-sm">
					<thead class="[&_tr]:border-b">
						<tr class="hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors">
							<th
								class="text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
							>
								Table
							</th>
							<th
								class="text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0"
							>
								Row
							</th>
						</tr>
					</thead>
					<tbody class="[&_tr:last-child]:border-0">
						{#each structure as table}
							<tr class="hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors">
								<td class="p-4 align-middle [&:has([role=checkbox])]:pr-0">
									<div class="flex items-center gap-2">
										<TableIcon class="text-muted-foreground h-4 w-4" />
										<span class="font-medium">{table.tableName}</span>
									</div>
								</td>
								<td class="text-muted-foreground p-4 align-middle [&:has([role=checkbox])]:pr-0">
									{formatNumber(table.rowCount)}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</DialogContent>
</Dialog>
