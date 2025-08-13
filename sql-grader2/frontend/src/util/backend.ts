import { toast } from 'svelte-sonner'
import type { AxiosError } from 'axios'
import { Backend } from '$/util/backend/backend'

export const SERVER_HOST = import.meta.env.VITE_SERVER_HOST
const SERVER_URL = `${SERVER_HOST}`
const controller = new AbortController()

export const backend = new Backend({
	baseURL: SERVER_URL,
	withCredentials: true,
	signal: controller.signal,
})

export const catcher = (err: AxiosError<{ message: string; error: string }>) => {
	if (err.response?.headers['content-type'].split(';')[0] === 'application/json' && err.response.data) {
		toast.error(err.response.data!.message!, {
			description: err.response.data!.error || undefined,
		})
	} else {
		toast.error(err.message)
	}
}

export const cat = (err: AxiosError<{ message: string; error: string }>) => {
	if (err.response?.headers['content-type'].split(';')[0] === 'application/json' && err.response.data) {
		return {
			message: err.response.data!.message!,
			description: err.response.data!.error || undefined,
		}
	} else {
		return {
			message: err.message,
		}
	}
}
