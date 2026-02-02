<script setup>
import { ref, onMounted, onUnmounted, watch, defineProps, defineEmits } from 'vue'

const props = defineProps({
  src: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['play', 'pause', 'ended'])

const audioRef = ref(null)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const progress = ref(0)
const isLoaded = ref(false)
const isSeeking = ref(false)

function formatTime(seconds) {
  if (!seconds || isNaN(seconds)) return '0:00'
  const hrs = Math.floor(seconds / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  if (hrs > 0) {
    return `${hrs}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function togglePlay() {
  if (!audioRef.value) return

  if (isPlaying.value) {
    audioRef.value.pause()
  } else {
    audioRef.value.play()
  }
}

function onTimeUpdate() {
  if (!audioRef.value || isSeeking.value) return
  currentTime.value = audioRef.value.currentTime
  if (duration.value > 0) {
    progress.value = (currentTime.value / duration.value) * 100
  }
}

function onLoadedMetadata() {
  if (!audioRef.value) return
  duration.value = audioRef.value.duration
  isLoaded.value = true
}

function onPlay() {
  isPlaying.value = true
  emit('play')
}

function onPause() {
  isPlaying.value = false
  emit('pause')
}

function onEnded() {
  isPlaying.value = false
  progress.value = 0
  currentTime.value = 0
  emit('ended')
}

function seek(event) {
  if (!audioRef.value || !isLoaded.value) return

  const progressBar = event.currentTarget
  const rect = progressBar.getBoundingClientRect()
  const clickX = event.clientX - rect.left
  const percentage = Math.max(0, Math.min(100, (clickX / rect.width) * 100))

  progress.value = percentage
  currentTime.value = (percentage / 100) * duration.value
  audioRef.value.currentTime = currentTime.value
}

function startSeek(event) {
  isSeeking.value = true
  seek(event)
  document.addEventListener('mousemove', handleSeekDrag)
  document.addEventListener('mouseup', endSeek)
  document.addEventListener('touchmove', handleTouchDrag)
  document.addEventListener('touchend', endSeek)
}

function handleSeekDrag(event) {
  if (!isSeeking.value) return
  const progressBar = document.querySelector(`[data-player-id="${playerId.value}"]`)
  if (!progressBar) return

  const rect = progressBar.getBoundingClientRect()
  const clickX = event.clientX - rect.left
  const percentage = Math.max(0, Math.min(100, (clickX / rect.width) * 100))

  progress.value = percentage
  currentTime.value = (percentage / 100) * duration.value
}

function handleTouchDrag(event) {
  if (!isSeeking.value || !event.touches[0]) return
  const progressBar = document.querySelector(`[data-player-id="${playerId.value}"]`)
  if (!progressBar) return

  const rect = progressBar.getBoundingClientRect()
  const touchX = event.touches[0].clientX - rect.left
  const percentage = Math.max(0, Math.min(100, (touchX / rect.width) * 100))

  progress.value = percentage
  currentTime.value = (percentage / 100) * duration.value
}

function endSeek() {
  if (audioRef.value && isSeeking.value) {
    audioRef.value.currentTime = currentTime.value
  }
  isSeeking.value = false
  document.removeEventListener('mousemove', handleSeekDrag)
  document.removeEventListener('mouseup', endSeek)
  document.removeEventListener('touchmove', handleTouchDrag)
  document.removeEventListener('touchend', endSeek)
}

const playerId = ref(`player-${Math.random().toString(36).substr(2, 9)}`)

onUnmounted(() => {
  document.removeEventListener('mousemove', handleSeekDrag)
  document.removeEventListener('mouseup', endSeek)
  document.removeEventListener('touchmove', handleTouchDrag)
  document.removeEventListener('touchend', endSeek)
})

watch(() => props.src, () => {
  isPlaying.value = false
  progress.value = 0
  currentTime.value = 0
  duration.value = 0
  isLoaded.value = false
})
</script>

<template>
  <div class="audio-player w-full">
    <!-- Hidden audio element -->
    <audio
      ref="audioRef"
      :src="src"
      preload="metadata"
      @timeupdate="onTimeUpdate"
      @loadedmetadata="onLoadedMetadata"
      @play="onPlay"
      @pause="onPause"
      @ended="onEnded"
    ></audio>

    <div class="flex items-center gap-3">
      <!-- Play/Pause Button -->
      <button
        @click="togglePlay"
        class="flex-shrink-0 w-10 h-10 rounded-full bg-primary-600 hover:bg-primary-700 text-white flex items-center justify-center transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
        :aria-label="isPlaying ? 'Pause' : 'Play'"
      >
        <svg v-if="!isPlaying" class="w-5 h-5 ml-0.5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M8 5v14l11-7z"/>
        </svg>
        <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
        </svg>
      </button>

      <!-- Progress Section -->
      <div class="flex-1 min-w-0">
        <!-- Progress Bar -->
        <div
          :data-player-id="playerId"
          @mousedown="startSeek"
          @touchstart.prevent="startSeek"
          class="relative h-2 bg-gray-200 dark:bg-gray-700 rounded-full cursor-pointer group touch-none"
        >
          <!-- Buffered/Background -->
          <div class="absolute inset-0 rounded-full overflow-hidden">
            <!-- Progress Fill -->
            <div
              class="h-full bg-primary-500 rounded-full transition-[width] duration-75"
              :style="{ width: `${progress}%` }"
            ></div>
          </div>
          <!-- Seek Handle -->
          <div
            class="absolute top-1/2 -translate-y-1/2 w-4 h-4 bg-primary-600 rounded-full shadow-md opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
            :class="{ 'opacity-100': isSeeking }"
            :style="{ left: `calc(${progress}% - 8px)` }"
          ></div>
        </div>

        <!-- Time Display -->
        <div class="flex justify-between mt-1.5 text-xs text-gray-500 dark:text-gray-400 tabular-nums">
          <span>{{ formatTime(currentTime) }}</span>
          <span>{{ formatTime(duration) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.audio-player {
  user-select: none;
}
</style>
