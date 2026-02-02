<script setup>
import { ref, onMounted } from 'vue'
import { useFeedsStore } from '../stores/feeds'
import FeedCard from '../components/FeedCard.vue'
import CreateFeedModal from '../components/CreateFeedModal.vue'

const store = useFeedsStore()
const showCreateModal = ref(false)

onMounted(() => {
  store.fetchFeeds()
})

async function handleCreateFeed(title, description) {
  try {
    await store.createFeed(title, description)
    showCreateModal.value = false
  } catch (e) {
    // Error handled in store
  }
}

async function handleDeleteFeed(id) {
  if (!confirm('Are you sure you want to delete this feed?')) return
  try {
    await store.deleteFeed(id)
  } catch (e) {
    alert(e.message)
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">My Feeds</h1>
      <button @click="showCreateModal = true" class="btn btn-primary">
        <span class="flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          New Feed
        </span>
      </button>
    </div>

    <div v-if="store.loading" class="text-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
      <p class="mt-4 text-gray-500">Loading feeds...</p>
    </div>

    <div v-else-if="store.error" class="card p-6 text-center text-red-500">
      {{ store.error }}
    </div>

    <div v-else-if="store.feeds.length === 0" class="card p-12 text-center">
      <svg class="w-16 h-16 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">No feeds yet</h3>
      <p class="mt-2 text-gray-500">Create your first feed to start adding content.</p>
    </div>

    <div v-else class="grid gap-4">
      <FeedCard
        v-for="feed in store.feeds"
        :key="feed.id"
        :feed="feed"
        @delete="handleDeleteFeed"
      />
    </div>

    <CreateFeedModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @create="handleCreateFeed"
    />
  </div>
</template>
