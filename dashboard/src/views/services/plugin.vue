<template>
  <div class="service-plugin-main gw-plugin-shell">
    <div class="gw-plugin-split gw-plugin-split--lg-row plugin-row">
      <div class="gw-plugin-split__col gw-plugin-split__col--nav">
        <a-card class="gw-plugin-card gw-plugin-nav-card" size="small" title="插件类型">
          <a-input-search
            v-model:value="data.searchKeyword"
            placeholder="搜索插件标识或描述"
            allow-clear
            style="margin-bottom: 12px"
            @search="fn.searchPlugin"
          />

          <div class="gw-plugin-list-scroll">
            <a-list
              :data-source="filteredPluginList"
              :loading="data.loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <a-list-item
                  :class="{ 'plugin-item-active': data.selectedPlugin?.res_id === item.res_id }"
                  @click="fn.selectPlugin(item)"
                >
                  <a-list-item-meta>
                    <template #title>
                      <span class="gw-plugin-list-title">
                        <i class="iconfont" :class="item.icon || 'icon-apex_plugin1'" />
                        {{ item.plugin_key }}
                      </span>
                    </template>
                    <template #description>
                      <span class="gw-plugin-list-desc">{{ item.description }}</span>
                    </template>
                  </a-list-item-meta>
                  <template #actions>
                    <a-tag :color="getTypeColor(item.type)">{{ item.typeName }}</a-tag>
                  </template>
                </a-list-item>
              </template>
            </a-list>
          </div>
        </a-card>
      </div>

      <div class="gw-plugin-split__col">
        <a-card
          class="gw-plugin-card gw-plugin-detail-card"
          size="small"
          :title="data.selectedPlugin ? `配置 · ${data.selectedPlugin.plugin_key}` : '插件配置'"
        >
          <div v-if="!data.selectedPlugin" class="gw-plugin-empty">
            <a-empty description="请从左侧选择一个插件" />
          </div>

          <div v-else class="gw-plugin-config-scroll">
            <div v-if="data.configuredPlugins.length > 0" style="margin-bottom: 16px">
              <a-table
                :columns="data.configColumns"
                :data-source="data.configuredPlugins"
                :pagination="false"
                size="small"
                bordered
              >
                <template #bodyCell="{ column, record }">
                  <template v-if="column.dataIndex === 'config'">
                    <span v-html="fn.getConfigSummary(record.plugin_key, record.config)"></span>
                  </template>
                  <template v-if="column.dataIndex === 'description'">
                    <a-tooltip :title="record.description" v-if="record.description">
                      <span>{{ record.description }}</span>
                    </a-tooltip>
                    <span v-else style="color: #999;">-</span>
                  </template>
                  <template v-if="column.dataIndex === 'enable'">
                    <a-switch
                      v-model:checked="record.enable"
                      size="small"
                      @change="fn.enableChange(record)"
                    />
                  </template>
                  <template v-if="column.dataIndex === 'operation'">
                    <a-space>
                      <a-button type="link" size="small" @click="fn.editConfig(record)">编辑</a-button>
                      <a-popconfirm
                        title="确认删除?"
                        ok-text="是"
                        cancel-text="否"
                        @confirm="fn.deleteConfig(record)"
                      >
                        <a-button type="link" size="small" danger>删除</a-button>
                      </a-popconfirm>
                    </a-space>
                  </template>
                </template>
              </a-table>
            </div>

            <div v-if="data.showConfigForm && data.selectedPlugin">
              <p class="gw-plugin-section-title">编辑表单</p>
              <component
                :key="`${data.configComponentName}-${data.editConfigResId || 'new'}-${data.selectedPlugin.res_id}`"
                :is="data.configComponentName"
                :pluginConfigData="data.editConfigData"
                :pluginOpType="data.pluginOpType"
                :pluginTag="data.selectedPlugin.plugin_key"
                :pluginConfigType="1"
                :targetResId="currentResId"
                :pluginResId="data.selectedPlugin.res_id"
                :pluginConfigResId="data.editConfigResId"
                @pluginAddVisible="fn.hideConfigForm"
                @pluginEditVisibleOff="fn.hideConfigForm"
                @componentRefreshList="fn.loadConfiguredPlugins"
              />
            </div>

            <div v-if="!data.showConfigForm" class="gw-plugin-empty" style="padding: 24px 0">
              <a-button type="primary" size="large" @click="fn.showAddForm">
                <template #icon>
                  <PlusOutlined />
                </template>
                添加配置
              </a-button>
            </div>
          </div>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import '@/assets/css/plugin-shell.css'
import { HookPluginList, HookPluginTypeIdNameMap, HookPluginKeyComponentMap } from '@/hooks'
import { $pluginConfigList, $pluginConfigEnable, $pluginConfigDelete } from '@/api'
import Plugin404 from '../plugin/components/err404.vue'
import Cors from '../plugin/components/cors.vue'
import Mock from '../plugin/components/mock.vue'
import KeyAuth from '../plugin/components/keyAuth.vue'
import JwtAuth from '../plugin/components/jwtAuth.vue'
import LimitReq from '../plugin/components/limitReq.vue'
import LimitConn from '../plugin/components/limitConn.vue'
import LimitCount from '../plugin/components/limitCount.vue'
import Waf from '../plugin/components/waf.vue'
import LogKafka from '../plugin/components/logKafka.vue'
import LogMysql from '../plugin/components/logMysql.vue'
import TrafficTag from '../plugin/components/trafficTag.vue'
import RequestRewrite from '../plugin/components/requestRewrite.vue'
import ResponseRewrite from '../plugin/components/responseRewrite.vue'

export default {
  components: {
    PlusOutlined,
    Plugin404,
    Cors,
    Mock,
    KeyAuth,
    JwtAuth,
    LimitReq,
    LimitConn,
    LimitCount,
    Waf,
    LogKafka,
    LogMysql,
    TrafficTag,
    RequestRewrite,
    ResponseRewrite
  },
  props: {
    currentResId: {
      type: String,
      default: ''
    }
  },
  emits: ['componentCloseDrawer', 'componentRefreshList'],
  setup(props, { emit }) {
    const data = reactive({
      searchKeyword: '',
      loading: false,
      pluginList: [],
      pluginTypeMap: {},
      selectedPlugin: null,
      configuredPlugins: [],
      showConfigForm: false,
      configComponentName: 'Plugin404',
      editConfigData: {},
      editConfigResId: null,
      pluginOpType: 1,
      configColumns: [
        { title: '配置名称', dataIndex: 'name', width: 150 },
        { title: '核心参数', dataIndex: 'config', width: 300 },
        { title: '备注', dataIndex: 'description', width: 200 },
        { title: '启用', dataIndex: 'enable', width: 80 },
        { title: '操作', dataIndex: 'operation', width: 120 }
      ]
    })

    const filteredPluginList = computed(() => {
      if (!data.pluginList || !Array.isArray(data.pluginList)) {
        return []
      }
      if (!data.searchKeyword || data.searchKeyword.trim() === '') {
        return [...data.pluginList]
      }
      const keyword = data.searchKeyword.toLowerCase().trim()
      return data.pluginList.filter(plugin => {
        if (!plugin || typeof plugin !== 'object') return false
        return (
          (plugin.plugin_key && String(plugin.plugin_key).toLowerCase().includes(keyword)) ||
          (plugin.description && String(plugin.description).toLowerCase().includes(keyword)) ||
          (plugin.typeName && String(plugin.typeName).toLowerCase().includes(keyword))
        )
      })
    })

    const getTypeColor = type => {
      const colorMap = {
        1: 'blue',
        2: 'orange',
        3: 'red',
        4: 'green',
        99: 'default'
      }
      return colorMap[type] || 'default'
    }

    const loadPluginList = async () => {
      data.loading = true
      try {
        const pluginList = await HookPluginList()
        const typeMap = await HookPluginTypeIdNameMap()

        data.pluginTypeMap = typeMap || {}
        if (Array.isArray(pluginList)) {
          data.pluginList = pluginList.map(plugin => ({
            ...plugin,
            typeName: typeMap[plugin.type] || '未知',
            icon: plugin.icon || 'icon-apex_plugin1'
          }))
        } else {
          data.pluginList = []
        }
      } catch (error) {
        message.error('加载插件列表失败')
        data.pluginList = []
      } finally {
        data.loading = false
      }
    }

    const loadConfiguredPlugins = async () => {
      if (!data.selectedPlugin || !props.currentResId) return

      try {
        const result = await $pluginConfigList(props.currentResId, 1)
        if (result.code === 0 && result.data?.list) {
          data.configuredPlugins = result.data.list
            .filter(item => item.plugin_key === data.selectedPlugin.plugin_key)
            .map(item => ({
              ...item,
              enable: item.enable === 1,
              key: item.res_id
            }))
        } else {
          data.configuredPlugins = []
        }
      } catch (error) {
        message.error('加载配置列表失败')
        data.configuredPlugins = []
      }
    }

    const selectPlugin = async plugin => {
      data.selectedPlugin = plugin
      data.showConfigForm = false
      data.editConfigResId = null
      data.editConfigData = {}
      await loadConfiguredPlugins()
    }

    const getComponentName = pluginKey => {
      const componentMap = HookPluginKeyComponentMap()
      return componentMap[pluginKey] || 'Plugin404'
    }

    const showAddForm = () => {
      if (!data.selectedPlugin) return
      data.pluginOpType = 1
      data.editConfigResId = null
      data.editConfigData = {}
      data.configComponentName = getComponentName(data.selectedPlugin.plugin_key)
      data.showConfigForm = true
    }

    const editConfig = record => {
      if (!data.selectedPlugin) return
      data.pluginOpType = 2
      data.editConfigResId = record.res_id
      data.editConfigData = {
        ...(record.config || {}),
        name: record.name || '',
        description: record.description || ''
      }
      data.configComponentName = getComponentName(data.selectedPlugin.plugin_key)
      data.showConfigForm = true
    }

    const hideConfigForm = () => {
      data.showConfigForm = false
      setTimeout(() => {
        data.editConfigResId = null
        data.editConfigData = {}
      }, 100)
    }

    const enableChange = async record => {
      const enableData = {
        enable: record.enable ? 1 : 2
      }
      try {
        const result = await $pluginConfigEnable(record.res_id, enableData, 1)
        if (result.code !== 0) {
          message.error(result.msg || '操作失败')
          record.enable = !record.enable
        } else {
          message.success('操作成功')
          emit('componentRefreshList')
        }
      } catch (error) {
        message.error('操作失败')
        record.enable = !record.enable
      }
    }

    const deleteConfig = async record => {
      try {
        const result = await $pluginConfigDelete(record.res_id, 1)
        if (result.code !== 0) {
          message.error(result.msg || '删除失败')
        } else {
          message.success('删除成功')
          await loadConfiguredPlugins()
          emit('componentRefreshList')
        }
      } catch (error) {
        message.error('删除失败')
      }
    }

    const searchPlugin = () => {
      // 搜索功能已通过 computed 实现
    }

    const getConfigSummary = (pluginKey, config) => {
      if (!config || !pluginKey) return '-'
      
      const parts = []
      
      switch (pluginKey) {
        case 'cors':
          if (config.allow_origins) parts.push(`源: ${config.allow_origins}`)
          if (config.allow_methods) parts.push(`方法: ${config.allow_methods}`)
          break
        case 'mock':
          if (config.response_type) parts.push(`类型: ${config.response_type}`)
          if (config.http_code) parts.push(`状态码: ${config.http_code}`)
          break
        case 'key-auth':
          if (config.secret) parts.push(`密钥: ${config.secret.substring(0, 8)}***`)
          break
        case 'jwt-auth':
          if (config.jwt_key) parts.push(`密钥: ${config.jwt_key.substring(0, 8)}***`)
          break
        case 'limit-req':
          if (config.rate) parts.push(`速率: ${config.rate}/s`)
          if (config.burst) parts.push(`突发: ${config.burst}`)
          break
        case 'limit-conn':
          if (config.rate) parts.push(`速率: ${config.rate}`)
          if (config.burst) parts.push(`突发: ${config.burst}`)
          if (config.default_conn_delay) parts.push(`延迟: ${config.default_conn_delay}s`)
          break
        case 'limit-count':
          if (config.count) parts.push(`次数: ${config.count}`)
          if (config.time_window) parts.push(`窗口: ${config.time_window}s`)
          break
        case 'waf':
          if (config.ip_whitelist?.ip_list?.length) {
            parts.push(`白名单: ${config.ip_whitelist.ip_list.length}个`)
          }
          if (config.ip_blacklist?.ip_list?.length) {
            parts.push(`黑名单: ${config.ip_blacklist.ip_list.length}个`)
          }
          if (config.rules?.rule_list?.length) {
            parts.push(`规则: ${config.rules.rule_list.length}条`)
          }
          break
        case 'log-kafka':
          if (config.brokers?.length) parts.push(`Brokers: ${config.brokers.length}个`)
          if (config.topic) parts.push(`Topic: ${config.topic}`)
          break
        case 'log-mysql':
          if (config.host) parts.push(`主机: ${config.host}`)
          if (config.database) parts.push(`数据库: ${config.database}`)
          if (config.table_name) parts.push(`表: ${config.table_name}`)
          break
        case 'traffic-tag':
          if (config.tags) {
            const tagKeys = Object.keys(config.tags)
            if (tagKeys.length) parts.push(`标签: ${tagKeys.join(', ')}`)
          }
          break
        case 'request-rewrite':
          if (config.uri_rewrite?.type) parts.push(`URI重写: ${config.uri_rewrite.type}`)
          if (config.headers && Object.keys(config.headers).length) {
            parts.push(`请求头: ${Object.keys(config.headers).length}个`)
          }
          if (config.query_args && Object.keys(config.query_args).length) {
            parts.push(`查询参数: ${Object.keys(config.query_args).length}个`)
          }
          break
        case 'response-rewrite':
          if (config.status_code) parts.push(`状态码: ${config.status_code}`)
          if (config.body_rewrite?.type) parts.push(`体重写: ${config.body_rewrite.type}`)
          if (config.headers && Object.keys(config.headers).length) {
            parts.push(`响应头: ${Object.keys(config.headers).length}个`)
          }
          break
        default:
          return '-'
      }
      
      return parts.length > 0 ? parts.join(' | ') : '-'
    }

    onMounted(async () => {
      await loadPluginList()
    })

    const fn = reactive({
      searchPlugin,
      selectPlugin,
      showAddForm,
      editConfig,
      hideConfigForm,
      enableChange,
      deleteConfig,
      loadConfiguredPlugins,
      getConfigSummary
    })

    return {
      data,
      filteredPluginList,
      fn,
      getTypeColor
    }
  }
}
</script>

<style scoped>
.service-plugin-main {
  padding: 0;
  height: 100%;
  min-height: 0;
  max-height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
}

.plugin-row {
  flex: 1;
  min-height: 0;
}

.gw-plugin-list-title .iconfont {
  margin-right: 8px;
  font-size: 16px;
}

.gw-plugin-list-desc {
  font-size: 12px;
  color: #64748b;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

:deep(.gw-plugin-nav-card.ant-card) {
  min-height: 0;
}

:deep(.gw-plugin-detail-card.ant-card) {
  height: 100%;
  min-height: 0;
}

:deep(.gw-plugin-nav-card .ant-card-body),
:deep(.gw-plugin-detail-card .ant-card-body) {
  padding: 12px 14px 16px;
}
</style>

