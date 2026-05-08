import store from '@/store'
import { canAccess } from '@/utils/permission'

export default {
  mounted(el, binding) {
    const code = binding.value
    if (!code) return
    const perms = store.state.user.userInfo.permissions
    if (!canAccess(perms, code)) {
      el.parentNode && el.parentNode.removeChild(el)
    }
  }
}
