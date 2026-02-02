import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useFeedsStore = defineStore('feeds', () => {
  const feeds = ref([])
  const currentFeed = ref(null)
  const episodes = ref([])
  const sources = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchFeeds() {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('/api/feeds')
      if (!response.ok) throw new Error('Failed to fetch feeds')
      feeds.value = await response.json()
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function fetchFeed(id) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`/api/feeds/${id}`)
      if (!response.ok) throw new Error('Failed to fetch feed')
      currentFeed.value = await response.json()
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function createFeed(title, description) {
    try {
      const response = await fetch('/api/feeds', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description }),
      })
      if (!response.ok) throw new Error('Failed to create feed')
      const newFeed = await response.json()
      feeds.value.push(newFeed)
      return newFeed
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function deleteFeed(id) {
    try {
      const response = await fetch(`/api/feeds/${id}`, {
        method: 'DELETE',
      })
      if (!response.ok) {
        const msg = await response.text()
        throw new Error(msg || 'Failed to delete feed')
      }
      feeds.value = feeds.value.filter(f => f.id !== id)
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function fetchEpisodes(feedId) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`/api/feeds/${feedId}/episodes`)
      if (!response.ok) throw new Error('Failed to fetch episodes')
      episodes.value = await response.json()
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function deleteEpisode(id) {
    try {
      const response = await fetch(`/api/episodes/${id}`, {
        method: 'DELETE',
      })
      if (!response.ok) throw new Error('Failed to delete episode')
      episodes.value = episodes.value.filter(e => e.id !== id)
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function fetchSources(feedId) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`/api/feeds/${feedId}/sources`)
      if (!response.ok) throw new Error('Failed to fetch sources')
      sources.value = await response.json()
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function addSource(feedId, url) {
    try {
      const response = await fetch(`/api/feeds/${feedId}/sources`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      if (!response.ok) {
        const msg = await response.text()
        throw new Error(msg || 'Failed to add source')
      }
      const newSource = await response.json()
      sources.value.push(newSource)
      // Refresh episodes after adding source
      await fetchEpisodes(feedId)
      return newSource
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function deleteSource(id) {
    try {
      const response = await fetch(`/api/sources/${id}`, {
        method: 'DELETE',
      })
      if (!response.ok) throw new Error('Failed to delete source')
      sources.value = sources.value.filter(s => s.id !== id)
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  return {
    feeds,
    currentFeed,
    episodes,
    sources,
    loading,
    error,
    fetchFeeds,
    fetchFeed,
    createFeed,
    deleteFeed,
    fetchEpisodes,
    deleteEpisode,
    fetchSources,
    addSource,
    deleteSource,
  }
})
