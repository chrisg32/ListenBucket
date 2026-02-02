<script setup>
import { useThemeStore } from '../stores/theme'
import { useRouter } from 'vue-router'

const theme = useThemeStore()
const router = useRouter()
</script>

<template>
  <header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
    <div class="container mx-auto px-4 max-w-4xl">
      <div class="flex items-center justify-between h-16">
        <router-link to="/" class="flex items-center gap-2 text-xl font-bold text-primary-600 dark:text-primary-400">
          <svg class="w-8 h-8" viewBox="0 0 100 100" fill="currentColor">
            <circle cx="50" cy="50" r="45" fill="currentColor"/>
            <path d="M35 30 L35 70 L70 50 Z" fill="white"/>
          </svg>
          ListenBucket
        </router-link>

        <button
          @click="theme.toggle()"
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          :title="theme.isDark ? 'Switch to light mode' : 'Switch to dark mode'"
        >
          <!-- Sun icon -->
          <svg v-if="theme.isDark" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <!-- Moon icon -->
          <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </button>
      </div>
    </div>
  </header>
</template>
