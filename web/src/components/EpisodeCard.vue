<script setup>
import { ref, computed, defineProps, defineEmits } from 'vue'

const props = defineProps({
  episode: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['delete'])

// Construct audio URL from episode ID to ensure it works from any device
const audioUrl = computed(() => {
  if (props.episode.status !== 'ready') return null
  return `/api/media/${props.episode.id}.mp3`
})

const audioRef = ref(null)
const isPlaying = ref(false)

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function formatDuration(seconds) {
  if (!seconds) return ''
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${minutes}:${secs.toString().padStart(2, '0')}`
}

function togglePlay() {
  if (audioRef.value) {
    if (isPlaying.value) {
      audioRef.value.pause()
    } else {
      audioRef.value.play()
    }
    isPlaying.value = !isPlaying.value
  }
}

function onEnded() {
  isPlaying.value = false
}

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  downloading: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  ready: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  error: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
}
</script>

<template>
  <div class="card p-4 overflow-hidden">
    <div class="flex gap-4 min-w-0">
      <div class="relative flex-shrink-0">
        <img
          v-if="episode.image_url"
          :src="episode.image_url"
          :alt="episode.title"
          class="w-24 h-24 sm:w-32 sm:h-32 rounded-lg object-cover"
        />
        <div v-else class="w-24 h-24 sm:w-32 sm:h-32 rounded-lg bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        </div>

        <!-- Play button overlay -->
        <button
          v-if="episode.status === 'ready' && audioUrl"
          @click="togglePlay"
          class="absolute inset-0 flex items-center justify-center bg-black/30 rounded-lg opacity-0 hover:opacity-100 transition-opacity"
        >
          <div class="w-12 h-12 rounded-full bg-white/90 flex items-center justify-center">
            <svg v-if="!isPlaying" class="w-6 h-6 text-gray-900 ml-1" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z"/>
            </svg>
            <svg v-else class="w-6 h-6 text-gray-900" fill="currentColor" viewBox="0 0 24 24">
              <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
            </svg>
          </div>
        </button>
      </div>

      <div class="flex-1 min-w-0">
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <h3 class="font-semibold text-gray-900 dark:text-white line-clamp-2">
              {{ episode.title }}
            </h3>
            <div class="flex flex-wrap items-center gap-2 mt-1">
              <span :class="['px-2 py-0.5 text-xs rounded-full', statusColors[episode.status]]">
                {{ episode.status }}
              </span>
              <span v-if="episode.duration && episode.status === 'ready'" class="text-sm text-gray-500">
                {{ formatDuration(episode.duration) }}
              </span>
              <span class="text-sm text-gray-500">
                {{ formatDate(episode.created_at) }}
              </span>
            </div>
          </div>
          <button
            @click="emit('delete', episode.id)"
            class="p-1 text-gray-400 hover:text-red-500 transition-colors flex-shrink-0"
            title="Delete episode"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>

        <p v-if="episode.description" class="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-2 break-words">
          {{ episode.description }}
        </p>

        <p v-if="episode.error_msg" class="mt-2 text-sm text-red-500 break-words line-clamp-2">
          {{ episode.error_msg }}
        </p>

        <a
          v-if="episode.source_url"
          :href="episode.source_url"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1 mt-2 text-sm text-primary-600 dark:text-primary-400 hover:underline"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
          View source
        </a>
      </div>
    </div>

    <!-- Audio player -->
    <audio
      v-if="episode.status === 'ready' && audioUrl"
      ref="audioRef"
      :src="audioUrl"
      class="w-full mt-4"
      controls
      @ended="onEnded"
      @play="isPlaying = true"
      @pause="isPlaying = false"
    ></audio>
  </div>
</template>
