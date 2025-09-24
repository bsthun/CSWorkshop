<script lang="ts">
	import { Avatar, AvatarImage, AvatarFallback } from '$/lib/shadcn/components/ui/avatar'

	export let user: {
		firstname: string
		lastname: string
		email: string
		pictureUrl?: string
	}
	export let size: 'sm' | 'md' | 'lg' = 'md'
	export let showEmail: boolean = true

	const getInitials = (firstname: string, lastname: string) => {
		return `${firstname.charAt(0)}${lastname.charAt(0)}`.toUpperCase()
	}

	const avatarSize = size === 'sm' ? 'h-6 w-6' : size === 'md' ? 'h-8 w-8' : 'h-12 w-12'
	const textSize = size === 'sm' ? 'text-xs' : size === 'md' ? 'text-xs' : 'text-sm'
</script>

<div class="flex items-center gap-3">
	<Avatar class={avatarSize}>
		<AvatarImage src={user.pictureUrl} alt="" />
		<AvatarFallback class={textSize}>
			{getInitials(user.firstname, user.lastname)}
		</AvatarFallback>
	</Avatar>
	<div>
		<div class="font-medium">
			{user.firstname}
			{user.lastname}
		</div>
		{#if showEmail}
			<div class="text-sm text-gray-500">{user.email}</div>
		{/if}
	</div>
</div>
