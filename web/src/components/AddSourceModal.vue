<script setup>
import { ref, computed, defineProps, defineEmits } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['close', 'add'])

const url = ref('')
const includeBackCatalog = ref(true)

// Detect if URL looks like a playlist or channel
const isPlaylistOrChannel = computed(() => {
  const urlValue = url.value.toLowerCase()
  return urlValue.includes('list=') ||
         urlValue.includes('/playlist') ||
         urlValue.includes('/@') ||
         urlValue.includes('/channel/') ||
         urlValue.includes('/c/') ||
         urlValue.includes('/user/')
})

function handleSubmit() {
  if (!url.value.trim()) return
  emit('add', url.value.trim(), includeBackCatalog.value)
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

        <!-- Back catalog toggle - shown when URL looks like playlist/channel -->
        <div v-if="isPlaylistOrChannel" class="mb-4 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <label for="backCatalog" class="text-sm font-medium text-gray-900 dark:text-white">
                Include existing videos
              </label>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Add all current videos from this playlist/channel
              </p>
            </div>
            <button
              type="button"
              role="switch"
              :aria-checked="includeBackCatalog"
              @click="includeBackCatalog = !includeBackCatalog"
              :class="[
                'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
                includeBackCatalog ? 'bg-primary-600' : 'bg-gray-300 dark:bg-gray-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  includeBackCatalog ? 'translate-x-5' : 'translate-x-0'
                ]"
              />
            </button>
          </div>
          <p v-if="!includeBackCatalog" class="mt-2 text-xs text-amber-600 dark:text-amber-400">
            Only new videos added after today will be included
          </p>
        </div>

        <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Supported formats:</h4>
          <ul class="text-sm text-gray-500 space-y-1">
            <li>* Video: youtube.com/watch?v=...</li>
            <li>* Playlist: youtube.com/playlist?list=...</li>
            <li>* Channel: youtube.com/@channel or /c/channel</li>
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
