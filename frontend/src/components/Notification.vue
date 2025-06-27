<template>
  <div :class="['notification', type]" v-if="visible">
    <slot>
      <span class="notification__message">{{ message }}</span>
    </slot>
    <button class="notification__close" @click="close">&times;</button>
  </div>
</template>

<script>
export default {
  name: 'Notification',
  props: {
    message: {
      type: String,
      default: ''
    },
    type: {
      type: String,
      default: 'info' // info, success, warning, error
    },
    duration: {
      type: Number,
      default: 3000 // ms, 0 = не скрывать автоматически
    }
  },
  data() {
    return {
      visible: true,
      timer: null
    }
  },
  mounted() {
    if (this.duration > 0) {
      this.timer = setTimeout(this.close, this.duration)
    }
  },
  beforeUnmount() {
    if (this.timer) clearTimeout(this.timer)
  },
  methods: {
    close() {
      this.visible = false
      this.$emit('close')
    }
  }
}
</script>

<style lang="scss" scoped>
.notification {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 16px 24px;
  margin: 8px 0;
  font-size: 1rem;
  position: relative;
  min-width: 240px;
  max-width: 400px;
  border-left: 4px solid #2196f3;
  transition: all 0.2s;

  &.success { border-color: #4caf50; }
  &.error { border-color: #f44336; }
  &.warning { border-color: #ff9800; }
  &.info { border-color: #2196f3; }

  &__message {
    flex: 1;
    color: #333;
  }
  &__close {
    background: none;
    border: none;
    font-size: 1.2em;
    color: #888;
    cursor: pointer;
    margin-left: 12px;
    transition: color 0.2s;
    &:hover { color: #333; }
  }
}
</style> 