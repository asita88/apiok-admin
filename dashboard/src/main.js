import { createApp } from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import router from '@/router'
import store from '@/store'
import '@/assets/font/iconfont.css'
import '@/assets/css/common.css'
import '@/assets/css/plugin-form.css'
import permissionDirective from '@/directives/permission'

;(function patchResizeObserverForStableLoop() {
  const OriginalResizeObserver = window.ResizeObserver
  if (!OriginalResizeObserver) return
  window.ResizeObserver = class extends OriginalResizeObserver {
    constructor(callback) {
      super((entries, observer) => {
        window.requestAnimationFrame(() => {
          window.requestAnimationFrame(() => {
            try {
              callback(entries, observer)
            } catch (e) {
              const m = String(e?.message ?? e ?? '')
              if (/ResizeObserver/i.test(m)) return
              throw e
            }
          })
        })
      })
    }
  }
})()

const RO_MSG = /ResizeObserver/i

const originalError = console.error
console.error = (...args) => {
  const first = args[0]
  const text =
    typeof first === 'string'
      ? first
      : first?.message ?? first?.reason?.message ?? ''
  if (text && RO_MSG.test(String(text))) return
  originalError.apply(console, args)
}

window.addEventListener(
  'error',
  event => {
    const m = event.message ?? ''
    if (m && RO_MSG.test(String(m))) {
      event.stopImmediatePropagation()
      event.preventDefault()
    }
  },
  true
)

window.addEventListener('unhandledrejection', event => {
  const msg = event.reason?.message ?? event.reason ?? ''
  if (msg && RO_MSG.test(String(msg))) {
    event.preventDefault()
  }
})

const app = createApp(App)
app.directive('permission', permissionDirective)
app.use(Antd)
app.use(store)
app.use(router).mount('#app')
