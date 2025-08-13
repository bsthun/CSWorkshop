<script lang="ts">
	import { createEventDispatcher } from 'svelte'
	import { Dialog, DialogContent, DialogHeader, DialogTitle } from '$/lib/shadcn/components/ui/dialog'
	import { Button } from '$/lib/shadcn/components/ui/button'
	import { Input } from '$/lib/shadcn/components/ui/input'
	import { AlertCircleIcon, CheckCircleIcon, FileIcon, UploadIcon } from 'lucide-svelte'
	import { backend, catcher } from '$/util/backend.ts'

	export let open = false
	export let collection: number
	export let schemaFilename: string | null = null

	const dispatch = createEventDispatcher()

	let fileInput: HTMLInputElement
	let uploadProgress = 0
	let uploadStatus: 'idle' | 'uploading' | 'success' | 'error' = 'idle'
	let errorMessage = ''
	let dragOver = false
	let selectedFile: File | null = null

	$: statusIcon = uploadStatus === 'success' ? CheckCircleIcon : uploadStatus === 'error' ? AlertCircleIcon : FileIcon
	$: statusText =
		uploadStatus === 'uploading'
			? `Uploading... ${uploadProgress}%`
			: uploadStatus === 'success'
				? 'Upload successful!'
				: uploadStatus === 'error'
					? errorMessage || 'Upload failed'
					: selectedFile
						? selectedFile.name
						: schemaFilename
							? schemaFilename
							: 'Select or drop a .sql file'
	$: color =
		uploadStatus === 'uploading'
			? 'text-blue-700'
			: uploadStatus === 'success'
				? 'text-green-700'
				: uploadStatus === 'error'
					? 'text-red-700'
					: 'text-gray-700'
	$: borderClasses = dragOver
		? 'border-blue-500 bg-blue-50'
		: uploadStatus === 'error'
			? 'border-red-300 bg-red-50'
			: uploadStatus === 'success'
				? 'border-green-300 bg-green-50'
				: 'border-gray-300'

	const handleDragOver = (event: DragEvent) => {
		event.preventDefault()
		dragOver = true
	}

	const handleDragLeave = () => {
		dragOver = false
	}

	const handleDrop = (event: DragEvent) => {
		event.preventDefault()
		dragOver = false

		if (event.dataTransfer?.files) {
			handleFiles(event.dataTransfer.files)
		}
	}

	const handleFileInputChange = (event: Event) => {
		const input = event.target as HTMLInputElement
		if (input.files) {
			handleFiles(input.files)
		}
	}

	const handleFiles = (files: FileList) => {
		if (files.length > 1) {
			errorMessage = 'Please upload only one file at a time'
			uploadStatus = 'error'
			setTimeout(() => {
				uploadStatus = 'idle'
				errorMessage = ''
			}, 3000)
			return
		}

		const file = files[0]
		if (!file.name.toLowerCase().endsWith('.sql')) {
			errorMessage = 'Please upload only .sql files'
			uploadStatus = 'error'
			setTimeout(() => {
				uploadStatus = 'idle'
				errorMessage = ''
			}, 3000)
			return
		}

		selectedFile = file
		uploadStatus = 'idle'
	}

	const uploadFile = () => {
		if (!selectedFile) return

		uploadStatus = 'uploading'
		uploadProgress = 0

		const formData = new FormData()
		formData.append('collectionId', collection.toString())
		formData.append('file', selectedFile)

		// Using native fetch with progress tracking
		const xhr = new XMLHttpRequest()

		xhr.upload.onprogress = (e) => {
			if (e.lengthComputable) {
				uploadProgress = Math.round((e.loaded * 100) / e.total)
			}
		}

		xhr.onload = () => {
			if (xhr.status === 200) {
				uploadStatus = 'success'
				dispatch('uploaded')
				setTimeout(() => {
					resetForm()
					open = false
				}, 2000)
			} else {
				uploadStatus = 'error'
				try {
					const response = JSON.parse(xhr.responseText)
					errorMessage = response.message || 'Upload failed'
				} catch {
					errorMessage = 'Upload failed. Please try again.'
				}
				setTimeout(() => {
					uploadStatus = 'idle'
					errorMessage = ''
				}, 3000)
			}
		}

		xhr.onerror = () => {
			uploadStatus = 'error'
			errorMessage = 'Upload failed. Please try again.'
			setTimeout(() => {
				uploadStatus = 'idle'
				errorMessage = ''
			}, 3000)
		}

		xhr.open('POST', '/api/admin/collection/schema/upload')
		xhr.send(formData)
	}

	const resetForm = () => {
		selectedFile = null
		uploadStatus = 'idle'
		uploadProgress = 0
		errorMessage = ''
		if (fileInput) {
			fileInput.value = ''
		}
	}

	const handleClose = () => {
		if (uploadStatus !== 'uploading') {
			resetForm()
		}
	}

	$: if (!open) {
		handleClose()
	}
</script>

<Dialog bind:open>
	<DialogContent class="max-w-md">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-2">
				<UploadIcon class="h-5 w-5" />
				Upload Schema
			</DialogTitle>
		</DialogHeader>

		<div class="space-y-4">
			<!-- File Upload Area -->
			<div
				class="flex items-center justify-between rounded-lg border-2 border-dashed p-4 transition-colors {borderClasses}"
				on:dragleave={handleDragLeave}
				on:dragover={handleDragOver}
				on:drop={handleDrop}
				role="application"
			>
				<div class="flex-shrink-0">
					<svelte:component this={statusIcon} class="h-6 w-6 {color}" />
				</div>

				<div class="flex-1 px-4">
					<p class="text-sm font-medium {color}">
						{statusText}
					</p>
					{#if uploadStatus === 'uploading'}
						<div class="mt-2 w-full rounded-full bg-gray-200">
							<div
								class="h-1.5 rounded-full bg-blue-500 transition-all duration-300"
								style="width: {uploadProgress}%"
							></div>
						</div>
					{/if}
				</div>

				<div class="flex-shrink-0">
					{#if uploadStatus === 'idle'}
						<Button size="sm" onclick={() => fileInput.click()} disabled={uploadStatus === 'uploading'}>
							Browse
						</Button>
					{/if}
				</div>
			</div>

			<input bind:this={fileInput} class="hidden" on:change={handleFileInputChange} type="file" accept=".sql" />

			{#if selectedFile && uploadStatus !== 'uploading' && uploadStatus !== 'success'}
				<Button class="w-full" onclick={uploadFile} disabled={uploadStatus === 'uploading'}>
					<UploadIcon class="mr-2 h-4 w-4" />
					Upload Schema
				</Button>
			{/if}
		</div>
	</DialogContent>
</Dialog>
