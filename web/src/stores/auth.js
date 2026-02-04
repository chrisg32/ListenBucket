import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const loading = ref(true)
  const setupRequired = ref(false)
  const error = ref(null)

  async function checkSetupStatus() {
    try {
      const response = await fetch('/api/auth/setup-status')
      if (!response.ok) throw new Error('Failed to check setup status')
      const data = await response.json()
      setupRequired.value = data.setup_required
      return data.setup_required
    } catch (e) {
      error.value = e.message
      return false
    }
  }

  async function checkAuth() {
    loading.value = true
    try {
      const response = await fetch('/api/auth/me')
      if (response.ok) {
        const data = await response.json()
        user.value = data.user
      } else {
        user.value = null
      }
    } catch (e) {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function setup(username, password) {
    error.value = null
    try {
      const response = await fetch('/api/auth/setup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })
      const data = await response.json()
      if (!response.ok) {
        throw new Error(data.message || 'Setup failed')
      }
      user.value = data.user
      setupRequired.value = false
      return data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function login(username, password) {
    error.value = null
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })
      const data = await response.json()
      if (!response.ok) {
        throw new Error(data.message || 'Login failed')
      }
      user.value = data.user
      return data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  async function logout() {
    try {
      await fetch('/api/auth/logout', { method: 'POST' })
    } catch (e) {
      // Ignore errors
    }
    user.value = null
  }

  return {
    user,
    loading,
    setupRequired,
    error,
    checkSetupStatus,
    checkAuth,
    setup,
    login,
    logout,
  }
})
