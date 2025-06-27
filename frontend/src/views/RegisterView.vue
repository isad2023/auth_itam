<template>
  <div class="auth-dark-bg full-center">
    <div class="auth-card">
      <h2>Регистрация</h2>
      <form @submit.prevent="handleRegister">
        <input v-model="name" type="text" placeholder="Имя" required />
        <input v-model="email" type="email" placeholder="Gmail" required />
        <input v-model="password" type="password" placeholder="Пароль" required />
        <button type="submit" :disabled="loading">Зарегистрироваться</button>
      </form>
      <div class="auth-switch">
        <span>Есть аккаунт?</span>
        <button @click="$router.push('/auth')">Войти</button>
      </div>
      <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
    </div>
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'

export default {
  name: 'RegisterView',
  components: { Notification },
  data() {
    return {
      name: '',
      email: '',
      password: '',
      notification: '',
      notificationType: 'info',
      loading: false
    }
  },
  methods: {
    async handleRegister() {
      this.loading = true
      this.notification = ''
      try {
        const res = await fetch('/auth/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name: this.name, email: this.email, password: this.password })
        })
        const data = await res.json()
        if (!res.ok) throw new Error(data.message || 'Ошибка')
        this.notification = 'Регистрация успешна!'
        this.notificationType = 'success'
        setTimeout(() => this.$router.push('/auth'), 1200)
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