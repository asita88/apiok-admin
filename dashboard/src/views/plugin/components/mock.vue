<template>
  <a-form
    class="plugin-form-surface"
    :model="data.formData"
    name="formData"
    :label-col="{ span: 5 }"
    :wrapper-col="{ span: 18 }"
    autocomplete="off"
    label-align="right"
    @finish="fn.onSubmit"
  >
    <a-form-item label="配置名称" name="name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="插件描述" name="description">
      <a-textarea v-model:value="data.formData.description" :rows="2" placeholder="请输入插件配置描述" />
    </a-form-item>

    <a-form-item label="response_type" name="response_type" :rules="schemaPluginMock.response_type">
      <a-select class="select" v-model:value="data.formData.response_type">
        <a-select-option value="application/json">application/json</a-select-option>
        <a-select-option value="text/html">text/html</a-select-option>
        <a-select-option value="text/xml">text/xml</a-select-option>
      </a-select>
    </a-form-item>

    <a-form-item label="http_code" name="http_code" :rules="schemaPluginMock.http_code">
      <a-input-number v-model:value="data.formData.http_code" style="width: 100%" />
    </a-form-item>

    <a-form-item label="http_body" name="http_body" :rules="schemaPluginMock.http_body">
      <a-textarea v-model:value="data.formData.http_body" :rows="4" />
    </a-form-item>

    <a-form-item label="http_headers" name="http_headers">
      <div
        class="plugin-form-kv-row"
        v-for="(item, index) in data.formData.http_headers"
        :key="item.id"
      >
        <a-form-item
          class="plugin-form-kv-field"
          :name="['http_headers', index, 'key']"
          :rules="checkHttpHeader"
        >
          <a-input placeholder="key" v-model:value="item.key" />
        </a-form-item>
        <a-form-item class="plugin-form-kv-field" :name="['http_headers', index, 'value']">
          <a-input placeholder="value" v-model:value="item.value" />
        </a-form-item>
        <div class="plugin-form-kv-actions">
          <a-button type="link" size="small" @click="fn.addHttpHeaders">
            <template #icon>
              <PlusOutlined />
            </template>
          </a-button>
          <a-button
            v-if="index > 0"
            type="link"
            size="small"
            danger
            @click="fn.removeHttpHeaders(item)"
          >
            <template #icon>
              <MinusCircleOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </a-form-item>
    <a-form-item label="启用" name="enable" v-show="pluginOpType === 1">
      <a-switch v-model:checked="data.formData.enable" size="small" />
    </a-form-item>

    <a-form-item class="plugin-form-actions" :wrapper-col="{ offset: 5, span: 18 }">
      <a-space>
        <a-button html-type="submit" type="primary">保存</a-button>
        <a-button @click="fn.cancel(pluginConfigData?.key)">取消</a-button>
      </a-space>
    </a-form-item>
  </a-form>
</template>
<script>
import { reactive, onMounted } from 'vue'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons-vue'
import { Form, message } from 'ant-design-vue'
import { schemaPluginMock } from '@/schema'
import { $pluginConfigAdd, $pluginConfigUpdate } from '@/api'

const useForm = Form.useForm
export default {
  components: {
    PlusOutlined,
    MinusCircleOutlined
  },
  props: {
    pluginConfigData: {
      Object
    },
    pluginConfigType: {
      Number
    },
    targetResId: {
      String
    },
    pluginConfigResId: {
      String
    },
    pluginOpType: {
      Number
    },
    pluginResId: {
      String
    }
  },
  emits: ['pluginAddVisible', 'pluginEditVisibleOff', 'componentRefreshList'],
  setup(props, { emit }) {
    onMounted(async () => {
      if (
        props.pluginConfigResId == null ||
        Object.keys(props.pluginConfigData.http_headers).length === 0
      ) {
        // 初始化一个空 http_headers
        addHttpHeaders()
      }
    })

    const data = reactive({
      formData: {
        name: 'plugin-mock',
        description: '',
        response_type: 'application/json',
        http_code: 200,
        http_body: '',
        http_headers: [],
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    // 接收的父级参数进行表单dom赋值，不需要监听其变化反应
    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.description != null) {
        data.formData.description = props.pluginConfigData.description
      }
      if (props.pluginConfigData.response_type != null) {
        data.formData.response_type = props.pluginConfigData.response_type
      }

      if (props.pluginConfigData.http_code != null) {
        data.formData.http_code = props.pluginConfigData.http_code
      }

      if (props.pluginConfigData.http_body != null) {
        data.formData.http_body = props.pluginConfigData.http_body
      }

      if (props.pluginConfigData.http_headers != null) {
        let hh = JSON.parse(JSON.stringify(props.pluginConfigData.http_headers))
        Object.getOwnPropertyNames(hh).forEach(function (k) {
          data.formData.http_headers.push({
            key: k,
            value: hh[k],
            id: Date.now()
          })
        })
      }
    }

    const addHttpHeaders = () => {
      data.formData.http_headers.push({
        key: undefined,
        vaule: undefined,
        id: Date.now()
      })
    }

    const removeHttpHeaders = item => {
      let index = data.formData.http_headers.indexOf(item)
      if (index !== -1) {
        data.formData.http_headers.splice(index, 1)
      }
    }

    // 提交当前插件的表单数据
    const onSubmit = async formData => {
      let hh = reactive({})
      formData.http_headers.forEach(item => {
        if (item.key != null) {
          let hv = ''
          if (item.value !== undefined) {
            hv = item.value.toString()
          }
          hh[item.key.toString()] = hv
        }
      })

      if (props.pluginConfigResId == null) {
        // 新增插件配置
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: formData.name ?? '',
          description: formData.description ?? '',
          enable: formData.enable == true ? 1 : 2,
          config: reactive({
            response_type: formData.response_type ?? '',
            http_code: formData.http_code ?? '',
            http_body: formData.http_body ?? '',
            http_headers: hh ?? []
          })
        })

        let { code, msg } = await $pluginConfigAdd(configData, props.pluginConfigType)
        if (code !== 0) {
          message.error(msg)
          return
        } else {
          message.success(msg)
          emit('pluginAddVisible')
          emit('componentRefreshList')
        }

        resetFields()
      } else {
        // 更新插件配置
        let configData = reactive({
          name: formData.name ?? '',
          description: formData.description ?? '',
          config: reactive({
            response_type: formData.response_type ?? '',
            http_code: formData.http_code ?? '',
            http_body: formData.http_body ?? '',
            http_headers: hh ?? []
          })
        })

        let { code, msg } = await $pluginConfigUpdate(
          props.pluginConfigResId,
          configData,
          props.pluginConfigType
        )
        if (code !== 0) {
          message.error(msg)
          return
        } else {
          message.success(msg)
          emit('pluginEditVisibleOff', props.pluginConfigData?.key)
          emit('componentRefreshList')
        }
      }
    }

    // 取消按钮
    const cancel = async key => {
      if (props.pluginOpType == 1) {
        // 调用父组件方法，收起增加插件的表单
        emit('pluginAddVisible')

        resetFields()
      } else {
        // 调用父组件方法，收起编辑插件的表单
        emit('pluginEditVisibleOff', key)
      }
    }

    const checkHttpHeader = [
      {
        validator: async (_, value) => {
          let pattern = /^[A-Za-z1-9_-]+$/
          if (value !== undefined && value.length !== 0 && !pattern.test(value)) {
            return Promise.reject('当前值仅包含字母、数字、划线')
          } else {
            return Promise.resolve().callback
          }
        }
      }
    ]

    const fn = reactive({
      addHttpHeaders,
      removeHttpHeaders,
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginMock, checkHttpHeader }
  }
}
</script>

