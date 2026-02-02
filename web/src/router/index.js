import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import FeedView from '../views/FeedView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/feed/:id',
    name: 'feed',
    component: FeedView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
