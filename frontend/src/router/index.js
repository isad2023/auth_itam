import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'
import AchievementsView from '../views/AchievementsView.vue'
import RequestsView from '../views/RequestsView.vue'
import NotificationsView from '../views/NotificationsView.vue'

const routes = [
  { path: '/', redirect: '/profile' },
  { path: '/auth', component: AuthView },
  { path: '/register', component: RegisterView },
  { path: '/profile', component: ProfileView },
  { path: '/achievements', component: AchievementsView },
  { path: '/requests', component: RequestsView },
  { path: '/notifications', component: NotificationsView }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Глобальный guard: если нет токена, редирект на /auth (кроме /auth и /register)
router.beforeEach((to, from, next) => {
  const publicPages = ['/auth', '/register']
  const authRequired = !publicPages.includes(to.path)
  const token = localStorage.getItem('token')
  if (authRequired && !token) {
    return next('/auth')
  }
  if ((to.path === '/auth' || to.path === '/register') && token) {
    return next('/profile')
  }
  next()
})

export default router
