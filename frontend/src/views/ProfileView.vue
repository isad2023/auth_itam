<template>
  <div class="profile-full-bg">
    <div class="profile-page-new">
      <h1 class="profile-title">Личный кабинет</h1>
      
      <section class="profile-section info-section">
        <h2 class="section-title">ЛИЧНАЯ ИНФОРМАЦИЯ</h2>
        <div class="info-flex">
          <div class="profile-avatar-block">
            <img v-if="photoPreview || user?.photoURL" :src="photoPreview || user.photoURL" class="profile-avatar" alt="avatar" />
            <div v-else class="profile-avatar profile-avatar-placeholder">{{ initials }}</div>
            <input ref="fileInput" type="file" accept="image/*" @change="onPhotoChange" style="display:none" />
            <div class="profile-avatar-actions">
              <button class="add-photo-btn" @click="$refs.fileInput.click()">Добавить фото +</button>
              <button class="delete-photo-btn" v-if="user?.photoURL || photoPreview" @click="removePhoto">Удалить</button>
            </div>
          </div>
          <form class="profile-form" @submit.prevent="saveProfile">
            <label>ФИО <span class="required">*</span>
              <input v-model="form.name" type="text" required />
            </label>
            <label>Электронная почта <span class="required">*</span>
              <input v-model="form.email" type="email" required />
            </label>
            <label>Роль <span class="required">*</span>
              <select v-model="form.specification">
                <option value="" disabled selected hidden>Выберите роль</option>
                <option v-for="role in roles" :key="role" :value="role">{{ role }}</option>
              </select>
            </label>
            <label>О себе
              <input v-model="form.about" type="text" />
            </label>
            <label>Telegram
              <input v-model="form.telegram" type="text" />
            </label>
            <label>Резюме (URL)
              <input v-model="form.resumeURL" type="text" />
            </label>
            <div class="form-actions">
              <button type="submit" class="save-btn" :disabled="!isChanged || loading">Сохранить</button>
            </div>
          </form>
        </div>
      </section>
      <section class="profile-section achievements-section">
        <div class="ach-header">
          <h2 class="section-title">Достижения</h2>
          <button class="add-ach-btn" @click="showAchForm = true">Добавить +</button>
        </div>
        <div v-if="showAchForm" class="ach-form-modal">
          <form class="ach-form" @submit.prevent="createAchievement">
            <label>Достижение <span class="required">*</span>
              <input v-model="newAch.title" required />
            </label>
            <label>Дата <span class="required">*</span>
              <input v-model="newAch.date" type="date" required />
            </label>
            <label>Баллы <span class="required">*</span>
              <input v-model="newAch.points" type="number" required />
            </label>
            <div class="form-actions">
              <button type="submit" class="save-btn">Сохранить</button>
              <button type="button" class="cancel-btn" @click="showAchForm = false">Отмена</button>
            </div>
          </form>
        </div>
        <div class="ach-table-wrap">
          <table class="ach-table">
            <thead>
              <tr>
                <th>Достижение</th>
                <th>Дата</th>
                <th>Баллы</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(ach, i) in achievements || []" :key="i">
                <td>{{ ach.title || ach.Title || ach.achievement || '-' }}</td>
                <td>{{ formatDate(ach.createdAt || ach.CreatedAt || ach.created_at) }}</td>
                <td>{{ ach.points || ach.Points || ach.score || '-' }}</td>
              </tr>
              <tr v-if="achievements && achievements.length && !achievements[0].title && !achievements[0].date && !achievements[0].points">
                <td colspan="3">{{ achievements[0] }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
      <Notification v-if="notification" :message="notification" :type="notificationType" @close="notification = ''" />
      <div class="profile-header-actions-fixed">
        <div class="notify-dropdown-wrap">
          <button class="notify-btn no-bg" @click="toggleNotifyDropdown" title="Уведомления">
            <svg class="notify-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M18 16V11C18 7.68629 15.3137 5 12 5C8.68629 5 6 7.68629 6 11V16L4 18V19H20V18L18 16Z" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M9 21C9 22.1046 9.89543 23 11 23H13C14.1046 23 15 22.1046 15 21" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
          <div v-if="showNotifyDropdown" class="notify-dropdown">
            <div v-if="notifications.length" class="notify-list">
              <div v-for="(n, i) in notifications" :key="i" class="notify-item">
                <span>{{ n.content }}</span>
                <span v-if="!n.isRead" class="notify-dot"></span>
              </div>
            </div>
            <div v-else class="notify-empty">Нет уведомлений</div>
          </div>
        </div>
        <button class="logout-btn" @click="logout" title="Выйти">
          <svg class="logout-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M16 17L21 12M21 12L16 7M21 12H9" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M12 19C11.4477 19 11 18.5523 11 18V16C11 15.4477 10.5523 15 10 15H7C5.89543 15 5 14.1046 5 13V11C5 9.89543 5.89543 9 7 9H10C10.5523 9 11 8.55228 11 8V6C11 5.44772 11.4477 5 12 5" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import Notification from '../components/Notification.vue'

export default {
  name: 'ProfileView',
  components: { Notification },
  data() {
    return {
      user: null,
      form: {
        name: '',
        email: '',
        about: '',
        telegram: '',
        specification: '',
        resumeURL: ''
      },
      roles: ['Project manager', 'Developer', 'Designer', 'QA', 'DevOps'],
      notification: '',
      notificationType: 'info',
      loading: false,
      photoPreview: '',
      photoFile: null,
      achievements: [],
      showAchForm: false,
      newAch: { event: '', title: '', date: '', points: '' },
      notifications: [
        { content: 'Новое достижение добавлено!', isRead: false },
        { content: 'Ваш профиль обновлён', isRead: true }
      ],
      showNotifyDropdown: false
    }
  },
  computed: {
    isChanged() {
      if (!this.user) return false
      return (
        this.form.about !== (this.user.about || '') ||
        this.form.telegram !== (this.user.telegram || '') ||
        this.form.specification !== (this.user.specification || '') ||
        this.photoFile !== null
      )
    },
    initials() {
      if (!this.form.name) return ''
      return this.form.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0,2)
    }
  },
  async mounted() {
    try {
      const token = localStorage.getItem('token')
      if (!token) throw new Error('Не авторизован')
      const res = await fetch('/auth/api/me', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (!res.ok) throw new Error('Ошибка получения профиля')
      this.user = await res.json()
      this.form.name = this.user.Name || ''
      this.form.email = this.user.Email || ''
      this.form.specification = this.user.Specification || ''
      this.form.about = this.user.About || ''
      this.form.telegram = this.user.Telegram || ''
      this.form.resumeURL = this.user.ResumeURL || ''
      await this.fetchAchievements()
    } catch (e) {
      this.notification = e.message || 'Ошибка'
      this.notificationType = 'error'
    }
  },
  methods: {
    logout() {
      localStorage.removeItem('token')
      this.$router.push('/auth')
    },
    onPhotoChange(e) {
      const file = e.target.files[0]
      if (!file) return
      this.photoFile = file
      const reader = new FileReader()
      reader.onload = e => {
        this.photoPreview = e.target.result
      }
      reader.readAsDataURL(file)
    },
    removePhoto() {
      this.photoFile = null
      this.photoPreview = ''
      this.form.resumeURL = ''
    },
    async saveProfile() {
      this.loading = true
      this.notification = ''
      try {
        const token = localStorage.getItem('token')
        const body = {
          Name: this.form.name,
          Email: this.form.email,
          Specification: this.form.specification,
          About: this.form.about,
          Telegram: this.form.telegram,
          ResumeURL: this.form.resumeURL
        }
        const res = await fetch('/auth/api/update_user_info', {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(body)
        })
        if (!res.ok) throw new Error('Ошибка сохранения профиля')
        this.notification = 'Профиль обновлён!'
        this.notificationType = 'success'
        const meRes = await fetch('/auth/api/me', {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (meRes.ok) {
          this.user = await meRes.json()
          this.form.name = this.user.Name || ''
          this.form.email = this.user.Email || ''
          this.form.specification = this.user.Specification || ''
          this.form.about = this.user.About || ''
          this.form.telegram = this.user.Telegram || ''
          this.form.resumeURL = this.user.ResumeURL || ''
        }
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      } finally {
        this.loading = false
      }
    },
    async fetchAchievements() {
      try {
        const token = localStorage.getItem('token')
        let url = '/auth/api/get_user_achievements';
        if (this.user && this.user.ID) {
          url += `?user_id=${this.user.ID}`;
        }
        const res = await fetch(url, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (!res.ok) throw new Error('Ошибка получения достижений')
        const data = await res.json();
        this.achievements = Array.isArray(data) ? data : [];
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    },
    async createAchievement() {
      try {
        const token = localStorage.getItem('token')
        const body = {
          title: this.newAch.title,
          points: this.newAch.points,
          user_id: this.user.ID
        }
        const res = await fetch('/auth/api/create_achievement', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(body)
        })
        if (!res.ok) throw new Error('Ошибка создания достижения')
        this.notification = 'Достижение добавлено!'
        this.notificationType = 'success'
        this.showAchForm = false
        this.newAch = { event: '', title: '', date: '', points: '' }
        await this.fetchAchievements()
      } catch (e) {
        this.notification = e.message || 'Ошибка'
        this.notificationType = 'error'
      }
    },
    toggleNotifyDropdown() {
      this.showNotifyDropdown = !this.showNotifyDropdown
    },
    formatDate(dateStr) {
      if (!dateStr) return '-';
      const d = new Date(dateStr);
      if (isNaN(d)) return '-';
      const day = String(d.getDate()).padStart(2, '0');
      const month = String(d.getMonth() + 1).padStart(2, '0');
      const year = d.getFullYear();
      return `${day}:${month}:${year}`;
    }
  }
}
</script>

<style lang="scss" scoped>
.profile-full-bg {
  min-height: 100vh;
  width: 100vw;
  background: #181b1f;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding-bottom: 64px;
}
:global(body) {
  background: #181b1f !important;
}
.profile-page-new {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 48px;
}
.profile-title {
  color: #fff;
  font-size: 2.6em;
  font-weight: 700;
  margin-top: 48px;
  margin-bottom: 24px;
  text-align: center;
  width: 100%;
}
.profile-section {
  width: 100%;
  max-width: 950px;
  margin: 0 auto 32px auto;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.info-section, .achievements-section {
  max-width: 950px;
  margin: 0 auto 32px auto;
  background: #23262b;
  border-radius: 18px;
  box-shadow: 0 4px 32px rgba(0,0,0,0.13);
  padding: 36px 40px 32px 40px;
}
.section-title {
  color: #fff;
  font-size: 1.35em;
  font-weight: 600;
  margin-bottom: 24px;
  letter-spacing: 0.04em;
}
.info-flex {
  display: flex;
  gap: 40px;
  align-items: flex-start;
  flex-wrap: wrap;
}
.profile-avatar-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 180px;
  margin-bottom: 12px;
}
.profile-avatar {
  width: 110px;
  height: 110px;
  border-radius: 50%;
  object-fit: cover;
  background: #1976d2;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.2em;
  font-weight: 600;
  margin-bottom: 12px;
}
.profile-avatar-placeholder {
  background: #1976d2;
}
.profile-avatar-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
}
.add-photo-btn {
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 1em;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  &:hover { background: #1565c0; }
}
.delete-photo-btn {
  background: #fff;
  color: #23262b;
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 1em;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  &:hover { background: #eee; }
}
.profile-form {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  label {
    color: #b0b3b8 !important;
    font-size: 1em;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  input, select {
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
  input[disabled], select[disabled] {
    opacity: 0.7;
    cursor: not-allowed;
  }
  .form-actions {
    margin-top: 8px;
    display: flex;
    justify-content: flex-end;
  }
  .save-btn {
    background: #1976d2;
    color: #fff;
    border: none;
    border-radius: 8px;
    padding: 12px 0;
    font-size: 1.1em;
    font-weight: 500;
    cursor: pointer;
    width: 100%;
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
.ach-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
}
.add-ach-btn {
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 8px 18px;
  font-size: 1em;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  &:hover { background: #1565c0; }
}
.ach-table-wrap {
  overflow-x: auto;
  max-width: 1200px;
  width: 100%;
}
.ach-table {
  width: 100%;
  border-collapse: collapse;
  background: transparent;
  th, td {
    padding: 12px 16px;
    text-align: left;
    color: #fff;
    font-size: 1em;
  }
  th {
    font-weight: 600;
    background: #23262b;
    color: #fff;
    border-bottom: 2px solid #35373b;
  }
  tr:not(:last-child) td {
    border-bottom: 1px solid #35373b;
  }
  tr, td {
    background: none !important;
    border-radius: 0 !important;
    box-shadow: none !important;
    margin: 0 !important;
  }
  min-width: 700px;
}
.profile-header-actions-fixed {
  position: fixed;
  top: 32px;
  right: 48px;
  display: flex;
  align-items: center;
  gap: 12px;
  z-index: 100;
}
.notify-dropdown-wrap {
  position: relative;
}
.notify-dropdown {
  position: absolute;
  top: 36px;
  right: 0;
  min-width: 260px;
  background: #23262b;
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.18);
  padding: 12px 0;
  z-index: 300;
}
.notify-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}
.notify-item {
  padding: 10px 18px;
  color: #fff;
  font-size: 1em;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #35373b;
  &:last-child { border-bottom: none; }
}
.notify-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ff5252;
  margin-left: 8px;
  display: inline-block;
}
.notify-empty {
  color: #b0b3b8;
  text-align: center;
  padding: 18px 0;
}
.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  cursor: pointer;
  .logout-icon {
    margin-right: 2px;
    vertical-align: middle;
  }
}
.required {
  color: #ff5252;
  margin-left: 2px;
  font-size: 1.1em;
}
.ach-form-modal {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}
.ach-form {
  background: #23262b;
  border-radius: 16px;
  box-shadow: 0 4px 32px rgba(0,0,0,0.25);
  padding: 32px 28px 24px 28px;
  min-width: 320px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  label {
    color: #b0b3b8 !important;
    font-size: 1em;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
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
  .form-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    margin-top: 8px;
  }
  .save-btn {
    background: #1976d2;
    color: #fff;
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-size: 1.1em;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
    &:hover:not(:disabled) {
      background: #1565c0;
    }
    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }
  .cancel-btn {
    background: #fff;
    color: #23262b;
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-size: 1.1em;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
    &:hover { background: #eee; }
  }
}
input, select, option, label {
  color: #fff !important;
}
.notify-btn.no-bg {
  background: none !important;
  box-shadow: none !important;
}
</style> 