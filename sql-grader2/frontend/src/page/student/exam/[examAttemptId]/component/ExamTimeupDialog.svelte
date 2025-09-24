<script lang="ts">
	import { navigate } from 'svelte-navigator'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { ClockIcon } from 'lucide-svelte'

	export let open = false

	const handleHome = () => {
		navigate('/student')
	}

	$: if (open) {
		setTimeout(() => {
			handleHome()
		}, 12000)
	}
</script>

<Dialog bind:open>
	<DialogContent class="sm:max-w-md">
		<DialogHeader>
			<div class="mb-4 flex items-center justify-center">
				<ClockIcon class="h-12 w-12 text-red-500" />
			</div>
			<DialogTitle class="text-center text-xl">Time's Up!</DialogTitle>
			<DialogDescription class="text-center">
				The examination time has ended. You will be redirected to the home page in shortly.
			</DialogDescription>
		</DialogHeader>
		<DialogFooter class="flex justify-center">
			<Button onclick={handleHome}>Home</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
