<script setup>
import { onMounted } from 'vue'
import { useThemeStore } from './stores/theme'
import Header from './components/Header.vue'

const theme = useThemeStore()

onMounted(() => {
  theme.init()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <Header />
    <main class="container mx-auto px-4 py-6 max-w-4xl">
      <router-view />
    </main>
  </div>
</template>
