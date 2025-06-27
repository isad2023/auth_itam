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
      notificationType: 'info'
    }
  },
  async mounted() {
    await this.fetchAchievements()
  },
  methods: {
    async fetchAchievements() {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch('/auth/api/get_user_achievements', {
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
        const token = localStorage.getItem('token')
        const res = await fetch('/auth/api/create_achievement', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            title: this.title,
            description: this.description,
            points: this.points
          })
        })
        if (!res.ok) throw new Error('Ошибка создания достижения')
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