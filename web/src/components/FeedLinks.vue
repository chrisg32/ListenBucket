<script setup>
import { ref, computed, defineProps } from 'vue'

const props = defineProps({
  feedId: {
    type: String,
    required: true,
  },
})

const copied = ref(false)

const rssUrl = computed(() => {
  const baseUrl = window.location.origin
  return `${baseUrl}/api/feeds/${props.feedId}/rss`
})

const podcastsUrl = computed(() => {
  return `podcast://${rssUrl.value.replace(/^https?:\/\//, '')}`
})

const overcastUrl = computed(() => {
  return `overcast://x-callback-url/add?url=${encodeURIComponent(rssUrl.value)}`
})

async function copyRssUrl() {
  try {
    await navigator.clipboard.writeText(rssUrl.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    alert('Failed to copy URL')
  }
}
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      @click="copyRssUrl"
      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
      </svg>
      {{ copied ? 'Copied!' : 'Copy RSS' }}
    </button>

    <a
      :href="rssUrl"
      target="_blank"
      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm bg-orange-100 dark:bg-orange-900 text-orange-700 dark:text-orange-300 hover:bg-orange-200 dark:hover:bg-orange-800 rounded-lg transition-colors"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
        <path d="M6.18 15.64a2.18 2.18 0 0 1 2.18 2.18C8.36 19 7.38 20 6.18 20C5 20 4 19 4 17.82a2.18 2.18 0 0 1 2.18-2.18M4 4.44A15.56 15.56 0 0 1 19.56 20h-2.83A12.73 12.73 0 0 0 4 7.27V4.44m0 5.66a9.9 9.9 0 0 1 9.9 9.9h-2.83A7.07 7.07 0 0 0 4 12.93V10.1Z"/>
      </svg>
      RSS Feed
    </a>

    <a
      :href="podcastsUrl"
      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-800 rounded-lg transition-colors"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
      </svg>
      Podcasts (iOS)
    </a>

    <a
      :href="overcastUrl"
      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm bg-orange-100 dark:bg-orange-900 text-orange-600 dark:text-orange-300 hover:bg-orange-200 dark:hover:bg-orange-800 rounded-lg transition-colors"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10"/>
      </svg>
      Overcast
    </a>
  </div>
</template>
