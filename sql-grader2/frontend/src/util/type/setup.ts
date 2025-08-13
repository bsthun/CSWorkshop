export type Setup = {
	profile: {
		id?: string
		name?: string
		email?: string
		userId?: string
	}
	initialized: boolean
	reload: () => Promise<void>
}
