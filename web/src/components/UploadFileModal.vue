<script setup>
import { ref, defineProps, defineEmits } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  progress: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['close', 'upload'])

const title = ref('')
const description = ref('')
const selectedFile = ref(null)
const isDragging = ref(false)

const supportedFormats = [
  '.mp3', '.mp4', '.m4a', '.m4v', '.mov', '.avi', '.mkv',
  '.webm', '.ogg', '.oga', '.opus', '.flac', '.wav', '.wma', '.aac', '.3gp', '.flv'
]

function handleDragOver(event) {
  event.preventDefault()
  isDragging.value = true
}

function handleDragLeave() {
  isDragging.value = false
}

function handleDrop(event) {
  event.preventDefault()
  isDragging.value = false

  const files = event.dataTransfer.files
  if (files.length > 0) {
    selectFile(files[0])
  }
}

function handleFileSelect(event) {
  const files = event.target.files
  if (files.length > 0) {
    selectFile(files[0])
  }
}

function selectFile(file) {
  const ext = '.' + file.name.split('.').pop().toLowerCase()
  if (!supportedFormats.includes(ext)) {
    alert('Unsupported file format: ' + ext)
    return
  }
  selectedFile.value = file
  if (!title.value) {
    title.value = file.name.replace(/\.[^/.]+$/, '')
  }
}

function clearFile() {
  selectedFile.value = null
}

function formatFileSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function handleSubmit() {
  if (!selectedFile.value) return
  emit('upload', selectedFile.value, title.value, description.value)
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="emit('close')">
    <div class="card w-full max-w-lg p-6">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Upload Audio/Video File</h2>

      <form @submit.prevent="handleSubmit">
        <!-- Drop Zone -->
        <div
          @dragover="handleDragOver"
          @dragleave="handleDragLeave"
          @drop="handleDrop"
          :class="[
            'relative border-2 border-dashed rounded-lg p-8 text-center transition-colors mb-4',
            isDragging
              ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
              : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500',
            loading ? 'pointer-events-none opacity-50' : ''
          ]"
        >
          <input
            type="file"
            id="fileInput"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            :accept="supportedFormats.join(',')"
            @change="handleFileSelect"
            :disabled="loading"
          />

          <div v-if="!selectedFile">
            <svg class="w-12 h-12 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
              <span class="font-medium text-primary-600 dark:text-primary-400">Click to upload</span>
              or drag and drop
            </p>
            <p class="mt-1 text-xs text-gray-500">
              Audio and video files up to 500MB
            </p>
          </div>

          <div v-else class="flex items-center justify-center gap-3">
            <svg class="w-10 h-10 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
            <div class="text-left">
              <p class="font-medium text-gray-900 dark:text-white truncate max-w-xs">
                {{ selectedFile.name }}
              </p>
              <p class="text-sm text-gray-500">
                {{ formatFileSize(selectedFile.size) }}
              </p>
            </div>
            <button
              type="button"
              @click.stop="clearFile"
              class="p-1 text-gray-400 hover:text-red-500 transition-colors"
              :disabled="loading"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Progress Bar -->
        <div v-if="loading && progress > 0" class="mb-4">
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-1">
            <span>Uploading...</span>
            <span>{{ progress }}%</span>
          </div>
          <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full bg-primary-500 transition-all duration-300"
              :style="{ width: progress + '%' }"
            ></div>
          </div>
        </div>

        <!-- Title -->
        <div class="mb-4">
          <label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Title
          </label>
          <input
            id="title"
            v-model="title"
            type="text"
            class="input"
            placeholder="Episode title"
            :disabled="loading"
          />
        </div>

        <!-- Description -->
        <div class="mb-6">
          <label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Description (optional)
          </label>
          <textarea
            id="description"
            v-model="description"
            rows="2"
            class="input"
            placeholder="Episode description"
            :disabled="loading"
          ></textarea>
        </div>

        <!-- Supported formats info -->
        <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Supported formats:</h4>
          <p class="text-xs text-gray-500">
            {{ supportedFormats.join(', ') }}
          </p>
        </div>

        <div class="flex justify-end gap-3">
          <button type="button" @click="emit('close')" class="btn btn-secondary" :disabled="loading">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" :disabled="!selectedFile || loading">
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Converting...
            </span>
            <span v-else>Upload</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
