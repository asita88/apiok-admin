<template>
  <!-- 左侧部分 -->
  <a-layout>
    <!-- 固定菜单 -->
    <a-layout-sider
      :style="{
        marginTop: '60px',
        overflow: 'auto',
        height: '100vh',
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0
      }"
    >
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline">
        <a-menu-item key="dashboard" @click="$router.push('/dashboard')">
          <DashboardOutlined />
          <span>大盘</span>
        </a-menu-item>
        <a-menu-item key="service" @click="$router.push('/service')">
          <AppstoreOutlined />
          <span>服务管理</span>
        </a-menu-item>
        <a-menu-item key="router" @click="$router.push('/router')">
          <LinkOutlined />
          <span>路由管理</span>
        </a-menu-item>
        <a-menu-item key="upstream" @click="$router.push('/upstream')">
          <ClusterOutlined />
          <span>上游管理</span>
        </a-menu-item>
        <a-menu-item key="ssl" @click="$router.push('/ssl')">
          <SafetyCertificateOutlined />
          <span>证书管理</span>
        </a-menu-item>
        <a-menu-item key="global-plugin" @click="$router.push('/global-plugin')">
          <GlobalOutlined />
          <span>全局插件</span>
        </a-menu-item>
        <a-sub-menu key="system">
          <template #title>
            <SettingOutlined />
            <span>系统管理</span>
          </template>
          <a-menu-item key="user" @click="$router.push('/user')">
            <UserOutlined />
            <span>用户管理</span>
          </a-menu-item>
          <a-menu-item key="log" @click="$router.push('/log')">
            <FileTextOutlined />
            <span>操作日志</span>
          </a-menu-item>
          <a-menu-item key="access-log" @click="$router.push('/access-log')">
            <FileTextOutlined />
            <span>访问日志</span>
          </a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <!-- 右侧 -->
    <a-layout>
      <!-- 固定头部 -->
      <a-layout-header class="header" :style="{ position: 'fixed', zIndex: 1, width: '100%' }">
        <div class="logo">
          <div>
            <img src="@/assets/img/apiok-logo-no-text.png" />
          </div>
          <div>
            <span>APIOK-Admin</span>
          </div>
        </div>

        <!-- 用户信息 -->
        <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="horizontal">
          <a-sub-menu key="sub1">
            <template #title>
              <div class="user">
                <span>{{ userName }}</span>
                <img src="../assets/img/user.gif" />
              </div>
            </template>
            <a-menu-item key="change-password" @click="changePassword">修改密码</a-menu-item>
            <a-menu-item key="logout" @click="logout">退出登录</a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-header>

      <!-- 主体 -->
      <a-layout :style="{ marginLeft: '200px', marginTop: '64px', height: 'calc(100vh - 64px)' }">
        <!-- 内容 -->
        <a-layout-content class="content">
          <router-view></router-view>
        </a-layout-content>

        <!-- 页脚 -->
        <a-layout-footer :style="{ textAlign: 'center' }">
          <!-- Ant Design ©2018 Created by Ant UED -->
        </a-layout-footer>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<script>
import { watch, ref } from 'vue'
import { $logout } from '@/api'
import { message } from 'ant-design-vue'
import {
  DashboardOutlined,
  AppstoreOutlined,
  LinkOutlined,
  ClusterOutlined,
  SafetyCertificateOutlined,
  GlobalOutlined,
  SettingOutlined,
  UserOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import router from '@/router'
import store from '@/store'
export default {
        components: {
          DashboardOutlined,
          AppstoreOutlined,
          LinkOutlined,
          ClusterOutlined,
          SafetyCertificateOutlined,
          GlobalOutlined,
          SettingOutlined,
          UserOutlined,
          FileTextOutlined
        },
  setup() {
    const selectedKeys = ref([router.currentRoute.value.name])
    const collapsed = ref(false)
    const { userInfo } = store.state.user

    const logout = async () => {
      const { code, msg } = await $logout()
      if (code == 0 || code == 401) {
        // 退出登录后清空本地用户信息和token缓存
        store.commit('user/setToken', {})
        // 跳转到登录页面
        router.replace('/login')
      } else {
        message.error(msg)
      }
    }

    const changePassword = () => {
      router.push('/user/change-password')
    }

    // 监听路由变化更新选中菜单
    watch(
      () => router.currentRoute.value.name,
      newName => {
        selectedKeys.value = [newName]
      }
    )

    return { selectedKeys, collapsed, logout, changePassword, userName: userInfo.username || userInfo.email || '' }
  }
}
</script>

<style lange="scss" scoped>
.content {
  margin: 10px;
  padding: 10px;
  background: #fff;
  height: calc(100%);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}
.logo {
  height: 30px;
  margin: 15px;
  display: flex;
  justify-content: left;
  align-items: center;
}
.logo img {
  height: 30px;
  margin-left: 6px;
}
.logo span {
  margin-left: 10px;
  color: #fff;
  font-size: 17px;
}
.header {
  background: #001529;
  padding: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.user span {
  margin-right: 10px;
}
.user img {
  width: 40px;
  height: 40px;
  border-radius: 20px;
  margin-right: 10px;
}
.iconfont {
  margin-right: 10px;
}
</style>
