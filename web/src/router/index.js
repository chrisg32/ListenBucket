import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import FeedView from '../views/FeedView.vue'
import LoginView from '../views/LoginView.vue'
import SetupView from '../views/SetupView.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { requiresAuth: true },
  },
  {
    path: '/feed/:id',
    name: 'feed',
    component: FeedView,
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { guest: true },
  },
  {
    path: '/setup',
    name: 'setup',
    component: SetupView,
    meta: { setup: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  // Check setup status on first navigation
  if (auth.loading) {
    await auth.checkSetupStatus()
    await auth.checkAuth()
  }

  // If setup is required, redirect to setup page
  if (auth.setupRequired && to.name !== 'setup') {
    return next({ name: 'setup' })
  }

  // If setup is complete but trying to access setup page
  if (!auth.setupRequired && to.meta.setup) {
    return next({ name: 'home' })
  }

  // Check if route requires auth
  if (to.meta.requiresAuth && !auth.user) {
    return next({ name: 'login' })
  }

  // If authenticated and trying to access guest-only route
  if (to.meta.guest && auth.user) {
    return next({ name: 'home' })
  }

  next()
})

export default router
