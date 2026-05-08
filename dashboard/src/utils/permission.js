export const ROUTE_PERM_ORDER = [
  'dashboard',
  'service',
  'router',
  'upstream',
  'ssl',
  'global-plugin',
  'user',
  'log',
  'access-log'
]

export function normalizePermissions(loginData) {
  if (!loginData || typeof loginData !== 'object') return ['*']
  if ('permissions' in loginData && Array.isArray(loginData.permissions)) {
    return loginData.permissions.map(String)
  }
  const roles = loginData.roles
  if (Array.isArray(roles)) {
    if (roles.includes('admin') || roles.includes('super_admin')) {
      return ['*']
    }
    if (roles.includes('viewer')) {
      return ['dashboard', 'service', 'router', 'upstream', 'ssl', 'global-plugin']
    }
  }
  if (loginData.role === 'admin' || loginData.role === 'super_admin') {
    return ['*']
  }
  return ['*']
}

export function firstPermittedRouteName(permissions) {
  if (permissions === undefined || permissions === null) return 'dashboard'
  if (!Array.isArray(permissions) || permissions.length === 0) return null
  if (permissions.includes('*')) return 'dashboard'
  for (const name of ROUTE_PERM_ORDER) {
    if (permissions.includes(name)) return name
  }
  return null
}

export function canAccess(permissions, code) {
  if (permissions === undefined || permissions === null) return true
  if (!Array.isArray(permissions) || permissions.length === 0) return false
  if (permissions.includes('*')) return true
  return permissions.includes(code)
}
