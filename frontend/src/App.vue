<script>
export default {
  name: 'App',
  computed: {
    isAuth() {
      return !!localStorage.getItem('token')
    },
    isProfile() {
      return this.$route.path === '/profile'
    }
  },
  methods: {
    logout() {
      localStorage.removeItem('token')
      this.$router.push('/auth')
    }
  }
}
</script>

<template>
  <div :id="'app'" :class="{ 'profile-bg': isProfile }">
    <template v-if="isAuth && !isProfile">
      <nav class="main-nav">
        <router-link to="/achievements">Достижения</router-link>
        <router-link to="/profile">Профиль</router-link>
        <router-link to="/requests">Запросы</router-link>
        <router-link to="/notifications">Уведомления</router-link>
        <button class="logout-btn" @click="logout">Выйти</button>
      </nav>
    </template>
    <router-view />
  </div>
</template>

<style lang="scss">
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}
#app {
  font-family: 'Segoe UI', Arial, sans-serif;
  background: #f5f7fa;
  min-height: 100vh;
}
#app.profile-bg {
  background: #181b1f !important;
}
.main-nav {
  display: flex;
  gap: 24px;
  justify-content: center;
  padding: 24px 0 12px 0;
  background: #fff;
  border-bottom: 1px solid #eee;
  margin-bottom: 24px;
  a {
    color: #2196f3;
    text-decoration: none;
    font-weight: 500;
    &:hover { text-decoration: underline; }
  }
  .logout-btn {
    background: none;
    border: none;
    color: #c62828;
    font-weight: 500;
    cursor: pointer;
    margin-left: 16px;
    &:hover { text-decoration: underline; }
  }
}
.main-content {
  padding: 0 16px;
  min-height: 70vh;
}
/* Ограничиваем ширину карточек */
.profile-view, .achievements-view, .requests-view, .notifications-view, .auth-view {
  max-width: 500px;
  width: 100%;
  margin: 40px auto;
  padding: 32px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
