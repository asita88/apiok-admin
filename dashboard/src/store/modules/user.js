import { normalizePermissions, firstPermittedRouteName, canAccess } from '@/utils/permission'

function emptyUser() {
  return {
    userId: null,
    userName: null,
    email: null,
    token: null,
    roles: [],
    permissions: []
  }
}

export default {
  namespaced: true,
  state() {
    return {
      userInfo: emptyUser()
    }
  },
  getters: {
    isLoggedIn: state => !!state.userInfo.token,
    permissions: state => state.userInfo.permissions || [],
    can: state => code => {
      if (!state.userInfo.token) return false
      return canAccess(state.userInfo.permissions, code)
    },
    firstRouteName: state => {
      if (!state.userInfo.token) return null
      const name = firstPermittedRouteName(state.userInfo.permissions)
      return name || 'forbidden'
    }
  },
  mutations: {
    setToken(state, loginData) {
      if (!loginData || !loginData.token) {
        state.userInfo = emptyUser()
        return Promise.resolve()
      }
      const permissions = normalizePermissions(loginData)
      state.userInfo = {
        userId: loginData.userId ?? null,
        userName: loginData.username ?? loginData.userName ?? null,
        email: loginData.email ?? null,
        token: loginData.token,
        roles: Array.isArray(loginData.roles) ? loginData.roles : loginData.role ? [loginData.role] : [],
        permissions
      }
      return Promise.resolve()
    }
  },
  actions: {}
}
