export const formatDate = (dateString: string | undefined): string => {
	const date = new Date(dateString!)
	return date.toLocaleDateString('en-GB', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
	})
}

export const formatDateTime = (dateString: string): string => {
	const date = new Date(dateString)
	return date.toLocaleDateString('en-GB', {
		day: 'numeric',
		month: 'short',
		year: 'numeric',
		hour: 'numeric',
		minute: '2-digit',
		hour12: true,
	})
}
