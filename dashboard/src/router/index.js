import { createRouter, createWebHashHistory } from 'vue-router'
import Error404 from '@/layouts/error404'
import Login from '@/views/user/login'
import ChangePassword from '@/views/user/change-password'
import Service from '@/views/services'
import Router from '@/views/router'
import Upstream from '@/views/upstream'
import Ssl from '@/views/ssl'
import store from '@/store'

import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const routes = [
  {
    path: '/login',
    name: 'login',
    meta: { public: true },
    component: Login
  },
  {
    path: '/forbidden',
    name: 'forbidden',
    meta: { public: true },
    component: () => import('@/views/permission/forbidden.vue')
  },
  {
    path: '/',
    name: 'layout',
    redirect: () => {
      if (!store.state.user.userInfo.token) {
        return { path: '/login' }
      }
      const name = store.getters['user/firstRouteName']
      return { name: name || 'forbidden' }
    },
    component: () => import('@/layouts/layout'),
    children: [
      {
        path: '/service',
        name: 'service',
        meta: { perm: 'service' },
        component: Service
      },
      {
        path: '/router',
        name: 'router',
        meta: { perm: 'router' },
        component: Router
      },
      {
        path: '/upstream',
        name: 'upstream',
        meta: { perm: 'upstream' },
        component: Upstream
      },
      {
        path: '/ssl',
        name: 'ssl',
        meta: { perm: 'ssl' },
        component: Ssl
      },
      {
        path: '/user/change-password',
        name: 'change-password',
        component: ChangePassword
      },
      {
        path: '/user',
        name: 'user',
        meta: { perm: 'user' },
        component: () => import('@/views/user')
      },
      {
        path: '/log',
        name: 'log',
        meta: { perm: 'log' },
        component: () => import('@/views/log')
      },
      {
        path: '/access-log',
        name: 'access-log',
        meta: { perm: 'access-log' },
        component: () => import('@/views/access-log')
      },
      {
        path: '/global-plugin',
        name: 'global-plugin',
        meta: { perm: 'global-plugin' },
        component: () => import('@/views/global-plugin')
      },
      {
        path: '/dashboard',
        name: 'dashboard',
        meta: { perm: 'dashboard' },
        component: () => import('@/views/dashboard')
      }
    ]
  },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: Error404 }
]

const router = createRouter({
  routes,
  history: createWebHashHistory()
})

router.beforeEach((to, from, next) => {
  NProgress.start()
  const { userInfo } = store.state.user

  if (to.meta.public) {
    if (to.path === '/login' && userInfo.token) {
      const name = store.getters['user/firstRouteName']
      next({ name: name || 'dashboard' })
      return
    }
    next()
    return
  }

  if (!userInfo.token) {
    next({ path: '/login' })
    return
  }

  const needPerm = to.matched.map(r => r.meta?.perm).filter(Boolean)
  const last = needPerm[needPerm.length - 1]
  if (last && !store.getters['user/can'](last)) {
    next({ name: 'forbidden' })
    return
  }

  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
