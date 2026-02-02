<script setup>
import { ref, defineProps, defineEmits } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['close', 'add'])

const url = ref('')

function handleSubmit() {
  if (!url.value.trim()) return
  emit('add', url.value.trim())
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="emit('close')">
    <div class="card w-full max-w-md p-6">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Add Source</h2>

      <form @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label for="url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            YouTube URL *
          </label>
          <input
            id="url"
            v-model="url"
            type="url"
            class="input"
            placeholder="https://www.youtube.com/watch?v=..."
            required
            autofocus
            :disabled="loading"
          />
          <p class="mt-2 text-sm text-gray-500">
            Supports videos, playlists, and channels
          </p>
        </div>

        <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Supported formats:</h4>
          <ul class="text-sm text-gray-500 space-y-1">
            <li>• Video: youtube.com/watch?v=...</li>
            <li>• Playlist: youtube.com/playlist?list=...</li>
            <li>• Channel: youtube.com/@channel or /c/channel</li>
          </ul>
        </div>

        <div class="flex justify-end gap-3">
          <button type="button" @click="emit('close')" class="btn btn-secondary" :disabled="loading">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" :disabled="!url.trim() || loading">
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Adding...
            </span>
            <span v-else>Add Source</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
