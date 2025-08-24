<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
	} from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { Label } from '$/lib/shadcn/components/ui/label'
	import { Loader2Icon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend'
	import { toast } from 'svelte-sonner'

	export let open = false

	let registerCode = ''
	let loading = false

	const dispatch = createEventDispatcher<{
		joined: void
	}>()

	const handleSubmit = () => {
		if (!registerCode.trim()) {
			return
		}

		loading = true
		backend.student
			.classJoin({ registerCode: registerCode.trim() })
			.then((response) => {
				dispatch('joined')
				registerCode = ''
				toast.success(response.message)
				open = false
			})
			.catch((err) => {
				catcher(err)
			})
			.finally(() => {
				loading = false
			})
	}

	const handleClose = () => {
		if (!loading) {
			open = false
			registerCode = ''
		}
	}

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter' && !loading) {
			handleSubmit()
		}
	}
</script>

<Dialog bind:open onOpenChange={handleClose}>
	<DialogContent class="sm:max-w-[425px]">
		<DialogHeader>
			<DialogTitle>Join Class</DialogTitle>
			<DialogDescription>Enter the class registration code provided by your instructor</DialogDescription>
		</DialogHeader>

		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="register-code">Class Registration Code</Label>
				<Input
					placeholder="Enter code"
					bind:value={registerCode}
					disabled={loading}
					onkeydown={handleKeyDown}
				/>
			</div>
		</div>

		<DialogFooter>
			<Button variant="outline" onclick={handleClose} disabled={loading}>Cancel</Button>
			<Button onclick={handleSubmit} disabled={loading || !registerCode.trim()}>
				{#if loading}
					<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
					Joining...
				{:else}
					Join Class
				{/if}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
