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

// 关闭加载图标
NProgress.configure({ showSpinner: false })

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/',
    name: 'layout',
    redirect: 'service',
    component: () => import('@/layouts/layout'),
    children: [
      {
        path: '/service',
        name: 'service',
        component: Service
      },
      {
        path: '/router',
        name: 'router',
        component: Router
      },
      {
        path: '/upstream',
        name: 'upstream',
        component: Upstream
      },
      {
        path: '/ssl',
        name: 'ssl',
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
        component: () => import('@/views/user')
      },
      {
        path: '/log',
        name: 'log',
        component: () => import('@/views/log')
      },
      {
        path: '/access-log',
        name: 'access-log',
        component: () => import('@/views/access-log')
      },
      {
        path: '/global-plugin',
        name: 'global-plugin',
        component: () => import('@/views/global-plugin')
      },
      {
        path: '/dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard')
      }
    ]
  },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: Error404 }
]

const router = new createRouter({
  routes,
  history: createWebHashHistory()
})

// 路由前置守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  const { userInfo } = store.state.user
  if (to.path == '/login') {
    if (userInfo.token) {
      router.replace('/')
    }
  } else if (to.path != '/login') {
    if (!userInfo.token) {
      router.replace('/login')
    }
  }
  next()
})

// 后置钩子
router.afterEach(async () => {
  NProgress.done()
})

export default router
