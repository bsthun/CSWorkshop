import { onMount } from 'svelte'
import { useLocation } from 'svelte-navigator'

export const scrollTop = () => {
	const location = useLocation()

	onMount(() => {
		const unsubscribe = location.subscribe(() => {
			window.scrollTo({ top: 0, behavior: 'smooth' })
		})

		return () => {
			unsubscribe()
		}
	})
}
