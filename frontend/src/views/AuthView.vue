<template>
  <div class="auth-dark-bg full-center">
    <div class="auth-card">
      <h2>Вход</h2>
      <form @submit.prevent="handleSubmit">
        <input v-model="email" type="email" placeholder="Gmail" required />
        <input v-model="password" type="password" placeholder="Пароль" required />
        <button type="submit" :disabled="loading">Войти</button>
      </form>
      <div class="auth-switch">
        <span>Нет аккаунта?</span>
        <button @click="$router.push('/register')">Зарегистрироваться</button>
      </div>
      <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
    </div>
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'

export default {
  name: 'AuthView',
  components: { Notification },
  data() {
    return {
      email: '',
      password: '',
      notification: '',
      notificationType: 'info',
      loading: false
    }
  },
  methods: {
    async handleSubmit() {
      this.loading = true
      this.notification = ''
      try {
        const res = await fetch('/auth/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: this.email, password: this.password })
        })
        const data = await res.json()
        if (!res.ok) throw new Error(data.message || 'Ошибка')
        localStorage.setItem('token', data.token || data.access_token)
        this.$router.push('/profile')
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.auth-dark-bg.full-center {
  position: fixed;
  inset: 0;
  min-height: 100vh;
  min-width: 100vw;
  background: #181b1f;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}
.auth-card {
  background: #23262b;
  border-radius: 16px;
  box-shadow: 0 4px 32px rgba(0,0,0,0.25);
  padding: 40px 32px 32px 32px;
  min-width: 340px;
  max-width: 400px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  h2 {
    color: #fff;
    margin-bottom: 24px;
    font-weight: 600;
    font-size: 1.5em;
  }
  form {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 18px;
    input {
      background: #181b1f;
      color: #fff;
      border: 1.5px solid #35373b;
      border-radius: 8px;
      padding: 12px 14px;
      font-size: 1em;
      outline: none;
      transition: border 0.2s;
      &:focus {
        border-color: #1976d2;
      }
    }
    button[type="submit"] {
      background: #1976d2;
      color: #fff;
      border: none;
      border-radius: 8px;
      padding: 12px 0;
      font-size: 1.1em;
      font-weight: 500;
      cursor: pointer;
      margin-top: 8px;
      transition: background 0.2s;
      &:hover:not(:disabled) {
        background: #1565c0;
      }
      &:disabled {
        opacity: 0.7;
        cursor: not-allowed;
      }
    }
  }
  .auth-switch {
    margin-top: 18px;
    color: #b0b3b8;
    font-size: 0.98em;
    display: flex;
    gap: 8px;
    align-items: center;
    button {
      background: none;
      border: none;
      color: #1976d2;
      cursor: pointer;
      text-decoration: underline;
      font-size: 1em;
      padding: 0;
      &:hover { color: #42a5f5; }
    }
  }
}
</style> 