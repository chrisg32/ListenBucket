<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useFeedsStore } from '../stores/feeds'
import EpisodeCard from '../components/EpisodeCard.vue'
import SourceCard from '../components/SourceCard.vue'
import AddSourceModal from '../components/AddSourceModal.vue'
import FeedLinks from '../components/FeedLinks.vue'

const route = useRoute()
const store = useFeedsStore()
const activeTab = ref('episodes')
const showAddSource = ref(false)
const addingSource = ref(false)

const feedId = computed(() => route.params.id)

onMounted(() => {
  loadData()
})

watch(feedId, () => {
  loadData()
})

function loadData() {
  store.fetchFeed(feedId.value)
  store.fetchEpisodes(feedId.value)
  store.fetchSources(feedId.value)
}

const readyEpisodes = computed(() => {
  return store.episodes.filter(e => e.status === 'ready')
})

const hasEpisodes = computed(() => {
  return store.episodes.length > 0
})

async function handleAddSource(url) {
  addingSource.value = true
  try {
    await store.addSource(feedId.value, url)
    showAddSource.value = false
  } catch (e) {
    alert(e.message)
  } finally {
    addingSource.value = false
  }
}

async function handleDeleteSource(id) {
  if (!confirm('Are you sure you want to delete this source?')) return
  try {
    await store.deleteSource(id)
  } catch (e) {
    alert(e.message)
  }
}

async function handleDeleteEpisode(id) {
  if (!confirm('Are you sure you want to delete this episode?')) return
  try {
    await store.deleteEpisode(id)
  } catch (e) {
    alert(e.message)
  }
}

function refreshData() {
  store.fetchEpisodes(feedId.value)
}
</script>

<template>
  <div v-if="store.currentFeed">
    <!-- Feed Header -->
    <div class="card p-6 mb-6">
      <div class="flex flex-col sm:flex-row gap-4">
        <img
          v-if="store.currentFeed.image_url"
          :src="store.currentFeed.image_url"
          :alt="store.currentFeed.title"
          class="w-24 h-24 rounded-lg object-cover flex-shrink-0"
        />
        <div v-else class="w-24 h-24 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center flex-shrink-0">
          <svg class="w-12 h-12 text-primary-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-start justify-between gap-4">
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white truncate">
                {{ store.currentFeed.title }}
                <span v-if="store.currentFeed.is_default" class="ml-2 text-sm font-normal text-primary-600 dark:text-primary-400">Default</span>
              </h1>
              <p v-if="store.currentFeed.description" class="mt-1 text-gray-600 dark:text-gray-400">
                {{ store.currentFeed.description }}
              </p>
            </div>
          </div>
          <FeedLinks :feed-id="feedId" class="mt-4" />
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-gray-200 dark:border-gray-700 mb-6">
      <button
        @click="activeTab = 'episodes'"
        :class="[
          'px-4 py-2 font-medium border-b-2 -mb-px transition-colors',
          activeTab === 'episodes'
            ? 'border-primary-500 text-primary-600 dark:text-primary-400'
            : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
        ]"
      >
        Episodes ({{ store.episodes.length }})
      </button>
      <button
        @click="activeTab = 'sources'"
        :class="[
          'px-4 py-2 font-medium border-b-2 -mb-px transition-colors',
          activeTab === 'sources'
            ? 'border-primary-500 text-primary-600 dark:text-primary-400'
            : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
        ]"
      >
        Sources ({{ store.sources.length }})
      </button>
    </div>

    <!-- Episodes Tab -->
    <div v-if="activeTab === 'episodes'">
      <div class="flex justify-end mb-4">
        <button @click="refreshData" class="btn btn-secondary text-sm">
          <svg class="w-4 h-4 mr-1 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>

      <div v-if="store.episodes.length === 0" class="card p-12 text-center">
        <svg class="w-16 h-16 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">No episodes yet</h3>
        <p class="mt-2 text-gray-500">Add a source to start generating episodes.</p>
        <button @click="activeTab = 'sources'" class="btn btn-primary mt-4">
          Go to Sources
        </button>
      </div>

      <div v-else class="grid gap-4">
        <EpisodeCard
          v-for="episode in store.episodes"
          :key="episode.id"
          :episode="episode"
          @delete="handleDeleteEpisode"
        />
      </div>
    </div>

    <!-- Sources Tab -->
    <div v-if="activeTab === 'sources'">
      <div class="flex justify-end mb-4">
        <button @click="showAddSource = true" class="btn btn-primary">
          <svg class="w-5 h-5 mr-1 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Source
        </button>
      </div>

      <div v-if="store.sources.length === 0" class="card p-12 text-center">
        <svg class="w-16 h-16 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">No sources yet</h3>
        <p class="mt-2 text-gray-500">Add a YouTube video, playlist, or channel URL to get started.</p>
        <button @click="showAddSource = true" class="btn btn-primary mt-4">
          Add Your First Source
        </button>
      </div>

      <div v-else class="grid gap-4">
        <SourceCard
          v-for="source in store.sources"
          :key="source.id"
          :source="source"
          @delete="handleDeleteSource"
        />
      </div>
    </div>

    <AddSourceModal
      v-if="showAddSource"
      :loading="addingSource"
      @close="showAddSource = false"
      @add="handleAddSource"
    />
  </div>

  <div v-else-if="store.loading" class="text-center py-12">
    <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
    <p class="mt-4 text-gray-500">Loading feed...</p>
  </div>

  <div v-else class="card p-12 text-center">
    <p class="text-gray-500">Feed not found</p>
    <router-link to="/" class="btn btn-primary mt-4 inline-block">
      Back to Home
    </router-link>
  </div>
</template>
