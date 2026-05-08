<template>
  <div class="main gw-plugin-shell">
    <header class="gw-plugin-hero gw-plugin-hero--compact">
      <div class="gw-plugin-hero__main">
        <div class="gw-plugin-hero__icon">
          <AppstoreOutlined />
        </div>
        <div>
          <h1 class="gw-plugin-hero__title">插件配置</h1>
          <p class="gw-plugin-hero__desc">绑定与管理网关插件实例，展开行可编辑参数。</p>
        </div>
      </div>
    </header>

    <a-card class="gw-plugin-card" size="small">
      <template #title>新增插件</template>
      <template #extra>
        <a-button type="primary" @click="fn.pluginAddVisible()">
          <template #icon>
            <PlusOutlined />
          </template>
          新增插件
        </a-button>
      </template>
      <div v-show="data.pluginAddVisible" class="gw-plugin-add-body">
        <a-tabs v-model:activeKey="data.addPluginTabActiveKey" type="card">
          <a-tab-pane key="1" tab="插件信息">
            <div class="gw-plugin-field">
              <span class="gw-plugin-field__label">选择插件</span>
              <a-select
                class="gw-plugin-field__control"
                :field-names="{
                  label: 'name',
                  value: 'res_id',
                  options: 'children'
                }"
                placeholder="请选择插件类型"
                :options="data.addPluginList"
                @change="fn.addPluginChange"
                style="width: 100%"
              />
            </div>

            <div v-show="data.addPluginComponent.infomationShow">
              <div class="gw-plugin-field">
                <span class="gw-plugin-field__label">插件类型</span>
                <a-input
                  class="gw-plugin-field__control"
                  v-model:value="data.addPluginType"
                  disabled
                  placeholder="插件类型"
                />
              </div>
              <div class="gw-plugin-field">
                <span class="gw-plugin-field__label">插件描述</span>
                <a-input
                  class="gw-plugin-field__control"
                  v-model:value="data.addPluginDesc"
                  disabled
                  placeholder="插件描述"
                />
              </div>
            </div>
          </a-tab-pane>
          <a-tab-pane key="2" tab="插件配置" :disabled="!data.addPluginComponent.pluginResId">
            <div class="plugin-add-form">
              <component
                :is="data.addPluginComponent.name"
                :pluginConfigData="data.addPluginComponent.pluginConfigData"
                :pluginOpType="data.pluginOpType"
                :pluginTag="data.addPluginComponent.tag"
                :pluginConfigType="pluginConfigType"
                :targetResId="currentResId"
                :pluginResId="data.addPluginComponent.pluginResId"
                :pluginConfigResId="null"
                @pluginAddVisible="fn.pluginAddVisible"
                @componentRefreshList="fn.componentRefreshList"
              />
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-card>

    <a-card class="gw-plugin-card gw-plugin-card--table" size="small">
      <template #title>已绑定插件</template>
      <a-table
        :columns="data.columns"
        :pagination="false"
        :data-source="data.list"
        v-model:expandedRowKeys="data.expandedRowKeys"
        :rowClassName="fn.tableRowClassName"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'icon'">
            <i class="iconfont" :class="[record.icon, record.color]" />
          </template>

          <template v-if="column.dataIndex === 'description'">
            <a-tooltip placement="topLeft">
              <template #title> {{ record.description }} </template>
              <span>
                {{ record.description }}
              </span>
            </a-tooltip>
          </template>

          <template v-if="column.dataIndex === 'enable'">
            <a-switch
              v-model:checked="record.enable"
              size="small"
              @click="fn.enableChange(record)"
            />
          </template>

          <template v-if="column.dataIndex === 'operation'">
            <a-space size="small">
              <a-tooltip placement="top">
                <template #title>展开编辑</template>
                <span
                  class="gw-plugin-op-icon"
                  role="button"
                  tabindex="0"
                  @click="fn.pluginConfigEditOn(record.key)"
                  @keydown.enter="fn.pluginConfigEditOn(record.key)"
                >
                  <EditOutlined />
                </span>
              </a-tooltip>
              <a-popconfirm
                placement="top"
                title="确认删除该插件配置？"
                ok-text="删除"
                cancel-text="取消"
                @confirm="fn.deleteFunc(record)"
              >
                <span class="gw-plugin-op-icon gw-plugin-op-icon--danger" role="button" tabindex="0">
                  <DeleteOutlined />
                </span>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
        <template #expandedRowRender="{ record }">
          <div class="plugin-edit-form">
            <component
              :is="record.component.name"
              :pluginConfigData="record.component.pluginConfigData"
              :pluginOpType="record.component.pluginOpType"
              :pluginTag="record.tag"
              :pluginConfigType="pluginConfigType"
              :targetResId="currentResId"
              :pluginResId="null"
              :pluginConfigResId="record.res_id"
              @pluginEditVisibleOff="fn.pluginEditVisibleOff"
              @componentRefreshList="fn.componentRefreshList"
            />
          </div>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import {
  AppstoreOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import '@/assets/css/plugin-shell.css'
import { $pluginConfigList, $pluginConfigEnable, $pluginConfigDelete, $globalPluginConfigList, $globalPluginConfigEnable, $globalPluginConfigDelete } from '@/api'
import { message } from 'ant-design-vue'
import { HookPluginKeyComponentMap, HookPluginTypeIdNameMap, HookPluginList } from '@/hooks'
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
    AppstoreOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
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
      String
    },
    pluginConfigType: {
      Number
    }
  },
  emits: ['componentCloseDrawer', 'componentRefreshList'],
  setup(props, { emit }) {
    onMounted(() => {
      if (props.pluginConfigType === 3 || props.currentResId !== null) {
        getList(props.currentResId || '')
      }
    })

    // 定义变量
    const data = reactive({
      expandedRowKeys: ref([]),
      columns: reactive([]),
      list: ref([]),
      pluginAddVisible: false, // 插件增加是否展示
      pluginOpType: 1, // 插件操作类型  1:增加   2:修改
      addPluginList: [], // 插件列表
      addPluginType: '',
      addPluginDesc: '',
      addPluginTabActiveKey: '1',
      addPluginComponent: reactive({
        name: 'Plugin404',
        tag: '',
        infomationShow: false,
        pluginResId: '',
        pluginConfigData: reactive({})
      })
    })

    // 定义hook数据
    const hookData = reactive({
      pluginResIdInfoMap: {}, // 插件信息hash列表
      pluginTypeIdNameMap: {}, // 插件类型hash列表
      pluginKeyToComponentMap: {} // 动态插件组件名称
    })

    const hookFnPluginKeyToComponent = pluginKey => {
      if (Object.values(hookData.pluginKeyToComponentMap).length == 0) {
        return 'Plugin404'
      }

      let componentName = hookData.pluginKeyToComponentMap[pluginKey]
      if (componentName == null) {
        componentName = 'Plugin404'
      }

      return componentName
    }

    // 定义列表头部
    data.columns = [
      { title: '', dataIndex: 'icon', width: 50 },
      { title: '名称', dataIndex: 'name' },
      { title: '标识', dataIndex: 'tag' },
      { title: '类型', dataIndex: 'type', width: 60 },
      { title: '描述', dataIndex: 'description', width: 190, ellipsis: true },
      { title: '启用', dataIndex: 'enable', width: 60 },
      { title: '操作', dataIndex: 'operation', width: 90 }
    ]

    // 获取插件列表
    const getList = async resId => {
      let code, dataList, msg
      let result
      if (props.pluginConfigType === 3) {
        // 全局插件
        result = await $globalPluginConfigList()
        code = result.code
        dataList = result.data
        msg = result.msg
      } else {
        result = await $pluginConfigList(resId, props.pluginConfigType)
        code = result.code
        dataList = result.data
        msg = result.msg
      }

      if (code !== 0) {
        message.error(msg)
        emit('componentCloseDrawer')
        return
      } else {
        if (Object.values(hookData.pluginTypeIdNameMap).length == 0) {
          hookData.pluginTypeIdNameMap = await HookPluginTypeIdNameMap()
        }
        if (Object.values(hookData.pluginKeyToComponentMap).length == 0) {
          hookData.pluginKeyToComponentMap = await HookPluginKeyComponentMap()
        }

        if (dataList.list && dataList.list.length > 0) {
          let pluginList = ref([])

          dataList.list.forEach(pluginConfigInfo => {
            let componentName = hookData.pluginKeyToComponentMap[pluginConfigInfo.plugin_key]
            if (componentName == null) {
              componentName = 'Plugin404'
            }

            pluginConfigInfo.config.key = pluginConfigInfo.res_id
            pluginConfigInfo.config.name = pluginConfigInfo.name
            let pluginConfigData = reactive(pluginConfigInfo.config)

            pluginList.value.push({
              key: pluginConfigInfo.res_id,
              res_id: pluginConfigInfo.res_id,
              icon: pluginConfigInfo.icon.length == 0 ? 'icon-apex_plugin1' : pluginConfigInfo.icon,
              name: pluginConfigInfo.name,
              tag: pluginConfigInfo.plugin_key,
              type: hookData.pluginTypeIdNameMap[pluginConfigInfo.plugin_type],
              description: pluginConfigInfo.plugin_description,
              enable: pluginConfigInfo.enable == 1 ? true : false,
              component: {
                name: hookFnPluginKeyToComponent(pluginConfigInfo.plugin_key),
                pluginOpType: 2,
                pluginConfigData: pluginConfigData
              }
            })
          })

          data.list = pluginList
        }
      }
    }

    // 增加插件是否展示动作
    const pluginAddVisible = async () => {
      if (data.pluginAddVisible === true) {
        data.pluginAddVisible = false
      } else {
        data.pluginAddVisible = true

        if (data.addPluginList.length == 0) {
          let pluginList = await HookPluginList()

          if (pluginList.length > 0) {
            pluginList.forEach(pluginInfo => {
              data.addPluginList.push({
                res_id: pluginInfo.res_id,
                name: pluginInfo.plugin_key
              })

              hookData.pluginResIdInfoMap[pluginInfo.res_id] = {
                type: hookData.pluginTypeIdNameMap[pluginInfo.type],
                description: pluginInfo.description
              }
            })
          }
        }
      }
    }

    // 编辑按钮——展示对应插件的配置信息
    const pluginConfigEditOn = async key => {
      let exist = false
      data.expandedRowKeys.forEach(k => {
        if (k == key) {
          exist = true
          return
        }
      })

      if (exist == false) {
        data.expandedRowKeys.push(key)
      } else {
        data.expandedRowKeys = data.expandedRowKeys.filter(t => t !== key)
      }
    }

    // 关闭对应的插件编辑模块
    const pluginEditVisibleOff = async key => {
      data.expandedRowKeys = data.expandedRowKeys.filter(t => t !== key)
    }

    // 刷新插件列表数据
    const componentRefreshList = async () => {
      getList(props.currentResId)
    }

    // 新增加插件时选择的插件基础数据（名称和描述）
    const addPluginChange = async (resId, option) => {
      let addPluginInfo = hookData.pluginResIdInfoMap[resId]
      data.addPluginComponent.name = hookFnPluginKeyToComponent(option.name)
      data.addPluginComponent.tag = option.name
      data.addPluginComponent.pluginResId = resId

      if (resId.length > 0) {
        data.addPluginComponent.infomationShow = true
      } else {
        data.addPluginComponent.infomationShow = false
      }

      if (addPluginInfo) {
        data.addPluginType = addPluginInfo.type
        data.addPluginDesc = addPluginInfo.description
      }
    }

    // 插件配置开关变化
    const enableChange = async record => {
      let enableData = reactive({
        enable: record.enable == true ? 1 : 2
      })
      let code, msg
      let result
      if (props.pluginConfigType === 3) {
        // 全局插件
        result = await $globalPluginConfigEnable(record.res_id, enableData)
        code = result.code
        msg = result.msg
      } else {
        result = await $pluginConfigEnable(
          record.res_id,
          enableData,
          props.pluginConfigType
        )
        code = result.code
        msg = result.msg
      }

      if (code !== 0) {
        message.error(msg)
        if (record.enable == true) {
          record.enable = false
        } else {
          record.enable = true
        }
        return
      } else {
        message.success(msg)
      }
    }

    // 插件配置删除
    const deleteFunc = async record => {
      let code, msg
      let result
      if (props.pluginConfigType === 3) {
        // 全局插件
        result = await $globalPluginConfigDelete(record.res_id)
        code = result.code
        msg = result.msg
      } else {
        result = await $pluginConfigDelete(record.res_id, props.pluginConfigType)
        code = result.code
        msg = result.msg
      }
      if (code !== 0) {
        message.error(msg)
        return
      } else {
        message.success(msg)
        getList(props.currentResId)
      }
    }

    // 定义函数
    const tableRowClassName = record =>
      record.enable ? 'gw-plugin-row--active' : 'gw-plugin-row--inactive'

    const fn = reactive({
      pluginAddVisible,
      pluginEditVisibleOff,
      addPluginChange,
      componentRefreshList,
      pluginConfigEditOn,
      enableChange,
      deleteFunc,
      tableRowClassName
    })

    return {
      data,
      fn,
      options: []
    }
  }
}
</script>

<style scoped>
.main {
  padding: 0 0 8px;
  min-height: 0;
}
.iconfont {
  font-size: 18px;
}
</style>
