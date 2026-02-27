<template>
  <div class="main">
    <a-form
      :model="data.formData"
      class="form"
      v-bind="{ labelCol: { xs: { span: 4 } }, wrapperCol: { xs: { span: 18 } } }"
      @finish="fn.onSubmit"
    >
      <a-form-item label="上游名称：" name="name" :rules="schemaUpstream.name">
        <a-input v-model:value="data.formData.name" />
      </a-form-item>

      <a-form-item label="负载均衡：" name="load_balance" :rules="schemaUpstream.load_balance">
        <a-select
          class="select"
          ref="select"
          v-model:value="data.formData.load_balance"
          placeholder="请选择"
        >
          <a-select-option value="1">加权轮询（Round Robin）</a-select-option>
          <a-select-option value="2">一致性哈希（CHash）</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item
        class="upstream_nodes_main"
        label="上游节点："
        name="upstream_nodes"
        :rules="schemaUpstream.upstream_nodes"
      >
        <div
          v-for="(item, index) in data.formData.upstream_nodes"
          :key="item.id"
          style="margin-bottom: 16px; padding: 12px; border: 1px solid #e8e8e8; border-radius: 4px;"
        >
          <a-space style="display: flex; margin-bottom: 8px" align="baseline">
            <a-form-item
              class="upstream_nodes_item"
              :name="['upstream_nodes', index, 'node_ip']"
              :rules="schemaUpstream.node_ip"
            >
              <a-input
                placeholder="IPV4"
                v-model:value="data.formData.upstream_nodes[index].node_ip"
              />
            </a-form-item>

            <a-form-item
              class="upstream_nodes_item"
              :name="['upstream_nodes', index, 'node_port']"
              :rules="schemaUpstream.node_port"
            >
              <a-input-number
                placeholder="Port"
                v-model:value="data.formData.upstream_nodes[index].node_port"
              />
            </a-form-item>

            <a-form-item
              class="upstream_nodes_item"
              :name="['upstream_nodes', index, 'node_weight']"
              :rules="schemaUpstream.node_weight"
            >
              <a-input-number
                placeholder="Weight"
                v-model:value="data.formData.upstream_nodes[index].node_weight"
              />
            </a-form-item>

            <a-form-item class="upstream_nodes_item">
              <a-button type="link" size="small" @click="fn.editTags(item, index)" title="编辑标签">
                标签
              </a-button>

              <a v-if="index > 0" @click="fn.removeIP(item)">
                <i class="iconfont color-red icon-jian"></i>
              </a>
            </a-form-item>
          </a-space>
          <div v-if="item.tags && Object.keys(item.tags).length > 0" style="margin-top: 8px;">
            <span style="margin-right: 8px; color: #666; font-size: 12px;">标签：</span>
            <a-tag
              v-for="(value, key) in item.tags"
              :key="key"
              color="blue"
              style="margin-right: 8px; margin-bottom: 4px;"
            >
              {{ key }}: {{ value }}
            </a-tag>
          </div>
        </div>
        <a-button type="dashed" @click="fn.addIP()" style="width: 100%; margin-top: 8px;">
          <i class="iconfont icon-tianjia"></i> 添加节点
        </a-button>
      </a-form-item>

      <a-modal
        v-model:visible="data.tagModalVisible"
        title="编辑标签"
        @ok="fn.saveTags"
        @cancel="fn.cancelTags"
      >
        <div v-for="(tag, tagIndex) in data.currentTagList" :key="tagIndex" style="margin-bottom: 12px">
          <a-space style="width: 100%">
            <a-input
              v-model:value="tag.key"
              placeholder="标签键"
              style="width: 200px"
            />
            <a-input
              v-model:value="tag.value"
              placeholder="标签值"
              style="width: 200px"
            />
            <a @click="fn.removeTag(tagIndex)">
              <i class="iconfont color-red icon-jian"></i>
            </a>
          </a-space>
        </div>
        <a-button type="dashed" @click="fn.addTag" style="width: 100%">
          <i class="iconfont icon-tianjia"></i> 添加标签
        </a-button>
      </a-modal>

      <a-divider orientation="left">健康检测</a-divider>

      <a-form-item label="健康检测：" style="margin-bottom: 16px">
        <a-switch v-model:checked="data.formData.checkEnabled" @change="fn.onCheckEnabledChange" />
      </a-form-item>

      <div v-if="data.formData.checkEnabled">
        <a-form-item 
          label="检测类型：" 
          style="margin-bottom: 16px"
        >
          <a-radio-group v-model:value="data.formData.checkTcp">
            <a-radio :value="true">TCP</a-radio>
            <a-radio :value="false">HTTP</a-radio>
          </a-radio-group>
        </a-form-item>

        <div v-if="!data.formData.checkTcp">
          <a-form-item 
            label="HTTP方法：" 
            style="margin-bottom: 16px"
          >
            <a-select
              v-model:value="data.formData.checkMethod"
              placeholder="请选择（可选）"
              allow-clear
              style="width: 100%"
            >
              <a-select-option value="GET">GET</a-select-option>
              <a-select-option value="POST">POST</a-select-option>
              <a-select-option value="HEADER">HEADER</a-select-option>
              <a-select-option value="OPTIONS">OPTIONS</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item 
            label="HTTP Host：" 
            style="margin-bottom: 16px"
          >
            <a-input
              v-model:value="data.formData.checkHost"
              placeholder="HTTP Host头（可选）"
            />
          </a-form-item>

          <a-form-item 
            label="HTTP URI：" 
            style="margin-bottom: 16px"
          >
            <a-input
              v-model:value="data.formData.checkUri"
              placeholder="/"
            />
            <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
              健康检测的HTTP路径，默认为 /
            </div>
          </a-form-item>
        </div>

        <a-form-item 
          label="检测间隔：" 
          style="margin-bottom: 16px"
        >
          <a-input-number
            v-model:value="data.formData.checkInterval"
            :min="0"
            :max="86400"
            placeholder="1"
            style="width: 100%"
          />
          <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
            健康检测间隔时间（秒），0-86400
          </div>
        </a-form-item>

        <a-form-item 
          label="超时时间：" 
          style="margin-bottom: 16px"
        >
          <a-input-number
            v-model:value="data.formData.checkTimeout"
            :min="0"
            :max="86400"
            placeholder="1"
            style="width: 100%"
          />
          <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
            健康检测超时时间（秒），0-86400
          </div>
        </a-form-item>

      </div>

      <a-form-item label="连接超时：" name="connect_timeout">
        <a-input v-model:value="data.formData.connect_timeout" placeholder="请输入（毫秒）" />
      </a-form-item>

      <a-form-item label="写超时：" name="write_timeout">
        <a-input v-model:value="data.formData.write_timeout" placeholder="请输入（毫秒）" />
      </a-form-item>

      <a-form-item label="读超时：" name="read_timeout">
        <a-input v-model:value="data.formData.read_timeout" placeholder="请输入（毫秒）" />
      </a-form-item>

      <a-form-item label="启用：" v-show="currentResId == null">
        <a-switch v-model:checked="data.formData.enable" />
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 10, offset: 16 }">
        <a-button type="primary" html-type="submit">保存</a-button>
        <a-button style="margin-left: 15px" @click="fn.cancel">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { $upstreamInfo, $upstreamAdd, $upstreamUpdate } from '@/api'
import { schemaUpstream } from '@/schema'

export default {
  props: {
    currentResId: null
  },
  emits: ['componentCloseDrawer', 'componentRefreshList'],
  setup(props, { emit }) {
    // 初始化——服务详情数据
    onMounted(async () => {
      if (props.currentResId !== null) {
        getInfo(props.currentResId)
      }
    })

    // 默认节点数据
    const defaultNode = {
      id: 0,
      node_ip: null,
      node_port: null,
      node_weight: null,
      tags: {}
    }

    // 定义变量
    const data = reactive({
      formData: {
        name: null,
        load_balance: null,
        connect_timeout: null,
        write_timeout: null,
        read_timeout: null,
        upstream_nodes: ref([defaultNode]),
        enable: false,
        checkEnabled: false,
        checkTcp: true,
        checkMethod: '',
        checkHost: '',
        checkUri: '/',
        checkInterval: 1,
        checkTimeout: 1
      },
      tagModalVisible: false,
      currentTagIndex: -1,
      currentTagList: []
    })

    // 获取详情
    const getInfo = async resId => {
      let { code, data: dataInfo, msg } = await $upstreamInfo(resId)

      if (code !== 0) {
        message.error(msg)
        emit('componentCloseDrawer')
        return
      } else {
        data.formData.name = dataInfo.name
        data.formData.load_balance = dataInfo.algorithm.toString()
        data.formData.enable = dataInfo.enable == 1 ? true : false
        data.formData.connect_timeout = dataInfo.connect_timeout
        data.formData.write_timeout = dataInfo.write_timeout
        data.formData.read_timeout = dataInfo.read_timeout

        if (dataInfo.node_list && dataInfo.node_list.length > 0) {
          data.formData.upstream_nodes = []

          dataInfo.node_list.forEach((item, index) => {
            let nodesInfo = JSON.parse(JSON.stringify(defaultNode))
            nodesInfo.id = index
            nodesInfo.node_ip = item.node_ip
            nodesInfo.node_port = item.node_port
            nodesInfo.node_weight = item.node_weight
            nodesInfo.tags = item.tags || {}
            data.formData.upstream_nodes.push(nodesInfo)
          })
        }

        // 从 upstream 级别读取健康检测配置
        if (dataInfo.check) {
          const check = dataInfo.check
          data.formData.checkEnabled = check.enabled || false
          data.formData.checkTcp = check.tcp !== undefined ? check.tcp : true
          data.formData.checkMethod = check.method || ''
          data.formData.checkHost = check.host || ''
          data.formData.checkUri = check.uri || '/'
          data.formData.checkInterval = check.interval || 1
          data.formData.checkTimeout = check.timeout || 1
        }
      }
    }

    // 增加节点元素
    const addIP = () => {
      data.formData.upstream_nodes.push({
        id: data.formData.upstream_nodes.length,
        node_ip: null,
        node_port: null,
        node_weight: null,
        tags: {}
      })
    }

    // 编辑标签
    const editTags = (item, index) => {
      data.currentTagIndex = index
      const tags = item.tags || {}
      data.currentTagList = Object.keys(tags).map(key => ({
        key: key,
        value: tags[key]
      }))
      if (data.currentTagList.length === 0) {
        data.currentTagList.push({ key: '', value: '' })
      }
      data.tagModalVisible = true
    }

    // 添加标签
    const addTag = () => {
      data.currentTagList.push({ key: '', value: '' })
    }

    // 删除标签
    const removeTag = (index) => {
      data.currentTagList.splice(index, 1)
    }

    // 保存标签
    const saveTags = () => {
      const tags = {}
      for (let i = 0; i < data.currentTagList.length; i++) {
        const key = data.currentTagList[i].key?.trim()
        const value = data.currentTagList[i].value?.trim()
        if (key && value) {
          tags[key] = value
        }
      }
      if (data.currentTagIndex >= 0) {
        data.formData.upstream_nodes[data.currentTagIndex].tags = tags
      }
      data.tagModalVisible = false
    }

    // 取消编辑标签
    const cancelTags = () => {
      data.tagModalVisible = false
      data.currentTagIndex = -1
      data.currentTagList = []
    }

    // 健康检测开关变化
    const onCheckEnabledChange = () => {
      if (!data.formData.checkEnabled) {
        // 如果关闭健康检测，重置为默认值
        data.formData.checkTcp = true
        data.formData.checkMethod = ''
        data.formData.checkHost = ''
        data.formData.checkUri = '/'
        data.formData.checkInterval = 1
        data.formData.checkTimeout = 1
      }
    }

    // 删除域名元素
    const removeIP = async item => {
      let index = data.formData.upstream_nodes.indexOf(item)
      if (index !== -1) {
        data.formData.upstream_nodes.splice(index, 1)
      }
    }

    // 点击提交保存当前数据
    const onSubmit = async () => {
      let formData = JSON.parse(JSON.stringify(data.formData))
      formData.load_balance = parseInt(formData.load_balance)
      formData.connect_timeout = parseInt(formData.connect_timeout)
      formData.write_timeout = parseInt(formData.write_timeout)
      formData.read_timeout = parseInt(formData.read_timeout)
      formData.enable = formData.enable == true ? 1 : 2

      // 处理健康检测配置（放在 upstream 级别）
      if (formData.checkEnabled) {
        formData.check = {
          enabled: true,
          tcp: formData.checkTcp,
          method: formData.checkMethod || '',
          host: formData.checkHost || '',
          uri: formData.checkUri || '/',
          interval: formData.checkInterval || 1,
          timeout: formData.checkTimeout || 1
        }
      } else {
        formData.check = {
          enabled: false,
          tcp: true,
          method: '',
          host: '',
          uri: '/',
          interval: 1,
          timeout: 1
        }
      }

      // 移除临时字段
      delete formData.checkEnabled
      delete formData.checkTcp
      delete formData.checkMethod
      delete formData.checkHost
      delete formData.checkUri
      delete formData.checkInterval
      delete formData.checkTimeout

      // 处理节点数据，包含标签
      if (formData.upstream_nodes && formData.upstream_nodes.length > 0) {
        formData.upstream_nodes = formData.upstream_nodes.map(node => {
          const nodeData = {
            node_ip: node.node_ip,
            node_port: node.node_port,
            node_weight: node.node_weight
          }
          if (node.tags && Object.keys(node.tags).length > 0) {
            nodeData.tags = node.tags
          }
          return nodeData
        })
      }

      // 调用增加/修改接口
      let res
      if (props.currentResId !== null) {
        res = await $upstreamUpdate(props.currentResId, formData)
      } else {
        res = await $upstreamAdd(formData)
      }

      if (res.code && res.code != 0) {
        message.error(res.msg)
        return
      } else {
        message.success(res.msg)
        emit('componentRefreshList')
      }
    }

    // 表单取消，关闭抽屉
    const cancel = async () => {
      emit('componentCloseDrawer')
    }

    // 定义函数
    const fn = reactive({
      addIP,
      removeIP,
      onSubmit,
      cancel,
      onCheckEnabledChange,
      editTags,
      addTag,
      removeTag,
      saveTags,
      cancelTags
    })

    return {
      schemaUpstream,
      data,
      fn
    }
  }
}
</script>

<style lange="scss" scoped>
.main {
  padding: 10px;
}
.form {
  width: '150px';
}
.upstream_nodes_main {
  margin-bottom: 0px;
}
.upstream_nodes_item {
  /* display: inline-block;
  justify-content: center;
  margin-right: 10px; */
}
</style>
