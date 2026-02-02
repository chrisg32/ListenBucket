<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  source: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['delete'])

function formatDate(dateString) {
  if (!dateString) return 'Never'
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const typeLabels = {
  video: 'Video',
  playlist: 'Playlist',
  channel: 'Channel',
}

const typeColors = {
  video: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  playlist: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
  channel: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
}
</script>

<template>
  <div class="card p-4 overflow-hidden">
    <div class="flex gap-4 min-w-0">
      <img
        v-if="source.image_url"
        :src="source.image_url"
        :alt="source.title"
        class="w-16 h-16 rounded-lg object-cover flex-shrink-0"
      />
      <div v-else class="w-16 h-16 rounded-lg bg-gray-200 dark:bg-gray-700 flex items-center justify-center flex-shrink-0">
        <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
      </div>

      <div class="flex-1 min-w-0">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <h3 class="font-medium text-gray-900 dark:text-white truncate">
              {{ source.title || 'Untitled Source' }}
            </h3>
            <div class="flex items-center gap-2 mt-1">
              <span :class="['px-2 py-0.5 text-xs rounded-full', typeColors[source.type]]">
                {{ typeLabels[source.type] }}
              </span>
            </div>
          </div>
          <button
            @click="emit('delete', source.id)"
            class="p-1 text-gray-400 hover:text-red-500 transition-colors flex-shrink-0"
            title="Delete source"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>

        <a
          :href="source.url"
          target="_blank"
          rel="noopener noreferrer"
          class="block mt-1 text-sm text-gray-500 truncate hover:text-primary-600 dark:hover:text-primary-400"
        >
          {{ source.url }}
        </a>

        <p v-if="source.type !== 'video'" class="mt-2 text-xs text-gray-500">
          Last checked: {{ formatDate(source.last_checked) }}
        </p>
      </div>
    </div>
  </div>
</template>
