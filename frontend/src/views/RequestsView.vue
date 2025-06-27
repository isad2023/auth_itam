<template>
  <div class="requests-view">
    <h2>Запросы</h2>
    <ul v-if="requests.length">
      <li v-for="r in requests" :key="r.id">
        <b>{{ r.type }}</b>: {{ r.description }} ({{ r.status || '—' }})
      </li>
    </ul>
    <p v-else>Нет запросов</p>
    <form @submit.prevent="createRequest">
      <input v-model="type" placeholder="Тип" required />
      <input v-model="description" placeholder="Описание" required />
      <button type="submit">Создать</button>
    </form>
    <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'

export default {
  name: 'RequestsView',
  components: { Notification },
  data() {
    return {
      requests: [],
      type: '',
      description: '',
      notification: '',
      notificationType: 'info'
    }
  },
  async mounted() {
    await this.fetchRequests()
  },
  methods: {
    async fetchRequests() {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch('/auth/api/get_all_user_requests', {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (!res.ok) throw new Error('Ошибка получения запросов')
        this.requests = await res.json()
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    },
    async createRequest() {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch('/auth/api/create_user_request', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            type: this.type,
            description: this.description
          })
        })
        if (!res.ok) throw new Error('Ошибка создания запроса')
        this.notification = 'Запрос создан!'
        this.notificationType = 'success'
        this.type = this.description = ''
        await this.fetchRequests()
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