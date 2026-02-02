<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  feed: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['delete'])

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <router-link
    :to="`/feed/${feed.id}`"
    class="card p-4 hover:shadow-md transition-shadow flex gap-4 group"
  >
    <img
      v-if="feed.image_url"
      :src="feed.image_url"
      :alt="feed.title"
      class="w-20 h-20 rounded-lg object-cover flex-shrink-0"
    />
    <div v-else class="w-20 h-20 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center flex-shrink-0">
      <svg class="w-10 h-10 text-primary-500" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
      </svg>
    </div>

    <div class="flex-1 min-w-0">
      <div class="flex items-start justify-between gap-2">
        <h3 class="font-semibold text-gray-900 dark:text-white truncate group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
          {{ feed.title }}
          <span v-if="feed.is_default" class="ml-1 text-xs font-normal text-primary-600 dark:text-primary-400">(Default)</span>
        </h3>
        <button
          v-if="!feed.is_default"
          @click.prevent="emit('delete', feed.id)"
          class="p-1 text-gray-400 hover:text-red-500 transition-colors opacity-0 group-hover:opacity-100"
          title="Delete feed"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
      <p v-if="feed.description" class="mt-1 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
        {{ feed.description }}
      </p>
      <p class="mt-2 text-xs text-gray-500">
        Updated {{ formatDate(feed.updated_at) }}
      </p>
    </div>
  </router-link>
</template>
