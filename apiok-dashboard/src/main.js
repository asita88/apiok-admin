import { createApp } from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import router from '@/router'
import store from '@/store'
import '@/assets/font/iconfont.css'
import '@/assets/css/common.css'

// 处理 ResizeObserver 警告
const resizeObserverErr = window.ResizeObserver
if (resizeObserverErr) {
  const resizeObserverErrHandler = window.ResizeObserver.prototype.observe
  window.ResizeObserver.prototype.observe = function (...args) {
    try {
      resizeObserverErrHandler.apply(this, args)
    } catch (e) {
      if (e.message !== 'ResizeObserver loop completed with undelivered notifications.') {
        throw e
      }
    }
  }
}

// 全局错误处理器 - 捕获 ResizeObserver 警告
const originalError = console.error
console.error = (...args) => {
  if (
    args[0] &&
    typeof args[0] === 'string' &&
    args[0].includes('ResizeObserver loop completed with undelivered notifications')
  ) {
    return
  }
  originalError.apply(console, args)
}

// 捕获未处理的 Promise rejection
window.addEventListener('error', event => {
  if (
    event.message &&
    event.message.includes('ResizeObserver loop completed with undelivered notifications')
  ) {
    event.preventDefault()
    return false
  }
})

// 捕获未处理的 Promise rejection
window.addEventListener('unhandledrejection', event => {
  if (
    event.reason &&
    event.reason.message &&
    event.reason.message.includes('ResizeObserver loop completed with undelivered notifications')
  ) {
    event.preventDefault()
    return false
  }
})

const app = createApp(App)
app.use(Antd)
app.use(store)
app.use(router).mount('#app')
