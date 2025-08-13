<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { ChevronLeftIcon, ChevronRightIcon, SearchIcon } from 'lucide-svelte'

	export let query: string = ''
	export let currentPage: number = 1
	export let totalItems: number = 0
	export let itemsPerPage: number = 24
	export let placeholder: string = 'Search...'

	const dispatch = createEventDispatcher<{
		search: { query: string; page: number }
		paginate: { page: number }
	}>()

	$: totalPages = Math.ceil(totalItems / itemsPerPage)
	$: hasNext = currentPage < totalPages
	$: hasPrevious = currentPage > 1

	let pageInputValue = currentPage.toString()
	$: pageInputValue = currentPage.toString()

	const handleSearch = () => {
		dispatch('search', { query: query, page: 1 })
	}

	const handleSearchKeydown = (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			handleSearch()
		}
	}

	const goToPage = (page: number) => {
		if (page >= 1 && page <= totalPages) {
			dispatch('paginate', { page })
		}
	}

	const handlePageInputKeydown = (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			const page = parseInt(pageInputValue, 10)
			if (!isNaN(page) && page >= 1 && page <= totalPages) {
				goToPage(page)
			} else {
				pageInputValue = currentPage.toString()
			}
		}
	}

	const handlePageInputBlur = () => {
		const page = parseInt(pageInputValue, 10)
		if (isNaN(page) || page < 1 || page > totalPages) {
			pageInputValue = currentPage.toString()
		}
	}

	const goToPrevious = () => {
		if (hasPrevious) {
			goToPage(currentPage - 1)
		}
	}

	const goToNext = () => {
		if (hasNext) {
			goToPage(currentPage + 1)
		}
	}
</script>

<div class="flex items-center gap-4">
	<!-- Search Input -->
	<div class="relative">
		<SearchIcon class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-gray-400" />
		<Input bind:value={query} class="w-64 pl-10" onkeydown={handleSearchKeydown} {placeholder} />
		<Button
			class="absolute top-1/2 right-1 h-7 w-7 -translate-y-1/2 p-0"
			onclick={handleSearch}
			size="sm"
			variant="ghost"
		>
			<SearchIcon class="h-3 w-3" />
		</Button>
	</div>

	<div class="flex items-center gap-1">
		<Button disabled={!hasPrevious} onclick={goToPrevious} variant="ghost">
			<ChevronLeftIcon class="h-4 w-4" />
		</Button>

		<div class="flex items-center gap-2 text-sm">
			<Input
				bind:value={pageInputValue}
				class="w-16 text-center"
				onblur={handlePageInputBlur}
				onkeydown={handlePageInputKeydown}
			/>
			<span> of {totalPages}</span>
		</div>

		<Button disabled={!hasNext} onclick={goToNext} variant="ghost">
			<ChevronRightIcon class="h-4 w-4" />
		</Button>
	</div>
</div>
