<template>
  <div class="notifications-view">
    <h2>Уведомления</h2>
    <ul v-if="notifications.length">
      <li v-for="n in notifications" :key="n.id">
        <Notification :message="n.content" :type="n.isRead ? 'info' : 'success'" :duration="0" />
      </li>
    </ul>
    <p v-else>Нет уведомлений</p>
    <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'

export default {
  name: 'NotificationsView',
  components: { Notification },
  data() {
    return {
      notifications: [],
      notification: '',
      notificationType: 'info'
    }
  },
  async mounted() {
    await this.fetchNotifications()
  },
  methods: {
    async fetchNotifications() {
      try {
        const token = localStorage.getItem('token')
        // Получаем user_id из профиля (если есть в localStorage)
        let userId = null;
        const meRes = await fetch('/auth/api/me', {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        if (meRes.ok) {
          const me = await meRes.json();
          userId = me.ID || me.id;
        }
        if (!userId) throw new Error('Не удалось получить user_id');
        const url = `/auth/api/get_all_notifications?user_id=${userId}`;
        const res = await fetch(url, {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!res.ok) throw new Error('Ошибка получения уведомлений');
        this.notifications = await res.json();
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
ul { padding-left: 0; list-style: none; }
li { margin-bottom: 8px; }
</style> 