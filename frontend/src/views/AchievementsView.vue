<template>
  <div class="achievements-view">
    <h2>Достижения</h2>
    <ul v-if="achievements.length">
      <li v-for="a in achievements" :key="a.id">
        <b>{{ a.title }}</b> — {{ a.description }} ({{ a.points }} баллов)
      </li>
    </ul>
    <p v-else>Нет достижений</p>
    <form @submit.prevent="createAchievement">
      <input v-model="title" placeholder="Название" required />
      <input v-model="description" placeholder="Описание" required />
      <input v-model.number="points" type="number" placeholder="Баллы" required />
      <button type="submit">Добавить</button>
    </form>
    <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'
import { apiUrl } from '../api.js'

export default {
  name: 'AchievementsView',
  components: { Notification },
  data() {
    return {
      achievements: [],
      title: '',
      description: '',
      points: 0,
      notification: '',
      notificationType: 'info',
      user_id: ''
    }
  },
  async mounted() {
    await this.fetchUserId()
    await this.fetchAchievements()
  },
  methods: {
    async fetchUserId() {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch(apiUrl('/auth/api/me'), {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (!res.ok) throw new Error('Ошибка получения пользователя')
        const user = await res.json()
        this.user_id = user.id || user.ID
        console.log('user_id получен из /auth/api/me:', this.user_id)
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    },
    async fetchAchievements() {
      try {
        const token = localStorage.getItem('token')
        let url = '/auth/api/get_user_achievements'
        if (this.user_id) {
          url += `?user_id=${this.user_id}`
        }
        const res = await fetch(apiUrl(url), {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (!res.ok) throw new Error('Ошибка получения достижений')
        this.achievements = await res.json()
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    },
    async createAchievement() {
      try {
        if (!this.user_id) {
          this.notification = 'Не удалось определить пользователя (user_id)';
          this.notificationType = 'error';
          return;
        }
        console.log('user_id перед отправкой:', this.user_id)
        const token = localStorage.getItem('token')
        const body = {
          title: this.title,
          description: this.description || "",
          points: Number(this.points),
          user_id: this.user_id
        }
        const res = await fetch(apiUrl('/auth/api/create_achievement'), {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(body)
        })
        if (!res.ok) {
          let errMsg = 'Ошибка создания достижения'
          try {
            const err = await res.json()
            if (err && err.error) errMsg += ': ' + err.error
          } catch {}
          throw new Error(errMsg)
        }
        this.notification = 'Достижение добавлено!'
        this.notificationType = 'success'
        this.title = this.description = ''
        this.points = 0
        await this.fetchAchievements()
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
ul { padding-left: 20px; }
li { margin-bottom: 8px; }
form {
  display: flex;
  gap: 8px;
  input {
    flex: 1;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    background: #2196f3;
    color: #fff;
    cursor: pointer;
    &:hover { background: #1976d2; }
  }
}
</style>