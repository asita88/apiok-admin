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
    <a-form-item label="配置名称" name="name" :rules="schemaPluginRequestRewrite.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="插件描述" name="description">
      <a-textarea v-model:value="data.formData.description" :rows="2" placeholder="请输入插件配置描述" />
    </a-form-item>

    <a-form-item label="URI重写" name="uri_rewrite">
      <a-card class="plugin-form-nested-card" size="small">
        <a-form-item label="类型" name="uri_rewrite.type">
          <a-select v-model:value="data.formData.uri_rewrite.type" style="width: 100%">
            <a-select-option value="regex">regex</a-select-option>
            <a-select-option value="replace">replace</a-select-option>
            <a-select-option value="prefix">prefix</a-select-option>
            <a-select-option value="suffix">suffix</a-select-option>
          </a-select>
        </a-form-item>
        <template v-if="data.formData.uri_rewrite.type === 'regex'">
          <a-form-item label="pattern" name="uri_rewrite.value.pattern">
            <a-input v-model:value="data.formData.uri_rewrite.value.pattern" placeholder="正则表达式模式" />
          </a-form-item>
          <a-form-item label="replacement" name="uri_rewrite.value.replacement">
            <a-input v-model:value="data.formData.uri_rewrite.value.replacement" placeholder="替换值" />
          </a-form-item>
          <a-form-item label="flags" name="uri_rewrite.value.flags">
            <a-input v-model:value="data.formData.uri_rewrite.value.flags" placeholder="正则标志(可选)" />
          </a-form-item>
        </template>
        <template v-else-if="data.formData.uri_rewrite.type === 'replace'">
          <a-form-item label="from" name="uri_rewrite.value.from">
            <a-input v-model:value="data.formData.uri_rewrite.value.from" placeholder="要替换的字符串" />
          </a-form-item>
          <a-form-item label="to" name="uri_rewrite.value.to">
            <a-input v-model:value="data.formData.uri_rewrite.value.to" placeholder="替换为" />
          </a-form-item>
        </template>
        <template v-else-if="data.formData.uri_rewrite.type === 'prefix' || data.formData.uri_rewrite.type === 'suffix'">
          <a-form-item label="remove" name="uri_rewrite.value.remove">
            <a-input v-model:value="data.formData.uri_rewrite.value.remove" placeholder="要移除的前缀/后缀" />
          </a-form-item>
          <a-form-item label="add" name="uri_rewrite.value.add">
            <a-input v-model:value="data.formData.uri_rewrite.value.add" placeholder="要添加的前缀/后缀" />
          </a-form-item>
        </template>
      </a-card>
    </a-form-item>

    <a-form-item label="请求头" name="headers">
      <div
        class="plugin-form-kv-row"
        v-for="(item, index) in data.formData.headers"
        :key="item.id"
      >
        <a-form-item class="plugin-form-kv-field" :name="['headers', index, 'key']">
          <a-input placeholder="header key" v-model:value="item.key" />
        </a-form-item>
        <a-form-item class="plugin-form-kv-field" :name="['headers', index, 'value']">
          <a-input placeholder="header value (留空表示删除)" v-model:value="item.value" />
        </a-form-item>
        <div class="plugin-form-kv-actions">
          <a-button type="link" size="small" @click="fn.addHeader">
            <template #icon>
              <PlusOutlined />
            </template>
          </a-button>
          <a-button
            v-if="index > 0"
            type="link"
            size="small"
            danger
            @click="fn.removeHeader(item)"
          >
            <template #icon>
              <MinusCircleOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </a-form-item>

    <a-form-item label="查询参数" name="query_args">
      <div
        class="plugin-form-kv-row"
        v-for="(item, index) in data.formData.query_args"
        :key="item.id"
      >
        <a-form-item class="plugin-form-kv-field" :name="['query_args', index, 'key']">
          <a-input placeholder="参数名" v-model:value="item.key" />
        </a-form-item>
        <a-form-item class="plugin-form-kv-field" :name="['query_args', index, 'value']">
          <a-input placeholder="参数值 (留空表示删除)" v-model:value="item.value" />
        </a-form-item>
        <div class="plugin-form-kv-actions">
          <a-button type="link" size="small" @click="fn.addQueryArg">
            <template #icon>
              <PlusOutlined />
            </template>
          </a-button>
          <a-button
            v-if="index > 0"
            type="link"
            size="small"
            danger
            @click="fn.removeQueryArg(item)"
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
import { schemaPluginRequestRewrite } from '@/schema'
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
    const data = reactive({
      formData: {
        name: 'plugin-request-rewrite',
        description: '',
        enabled: true,
        uri_rewrite: {
          type: 'regex',
          value: {
            pattern: '',
            replacement: '',
            flags: ''
          }
        },
        headers: [],
        query_args: [],
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    onMounted(() => {
      if (props.pluginConfigResId == null) {
        addHeader()
        addQueryArg()
      }
    })

    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.description != null) {
        data.formData.description = props.pluginConfigData.description
      }
      if (props.pluginConfigData.enabled != null) {
        data.formData.enabled = props.pluginConfigData.enabled
      }
      if (props.pluginConfigData.uri_rewrite != null) {
        data.formData.uri_rewrite = JSON.parse(JSON.stringify(props.pluginConfigData.uri_rewrite))
        if (!data.formData.uri_rewrite.value) {
          data.formData.uri_rewrite.value = {}
        }
      }
      if (props.pluginConfigData.headers != null) {
        const headers = props.pluginConfigData.headers
        Object.keys(headers).forEach(key => {
          data.formData.headers.push({
            key: key,
            value: headers[key],
            id: Date.now() + Math.random()
          })
        })
      }
      if (props.pluginConfigData.query_args != null) {
        const queryArgs = props.pluginConfigData.query_args
        Object.keys(queryArgs).forEach(key => {
          data.formData.query_args.push({
            key: key,
            value: queryArgs[key],
            id: Date.now() + Math.random()
          })
        })
      }
    }

    const addHeader = () => {
      data.formData.headers.push({
        key: '',
        value: '',
        id: Date.now() + Math.random()
      })
    }

    const removeHeader = item => {
      const index = data.formData.headers.indexOf(item)
      if (index !== -1) {
        data.formData.headers.splice(index, 1)
      }
    }

    const addQueryArg = () => {
      data.formData.query_args.push({
        key: '',
        value: '',
        id: Date.now() + Math.random()
      })
    }

    const removeQueryArg = item => {
      const index = data.formData.query_args.indexOf(item)
      if (index !== -1) {
        data.formData.query_args.splice(index, 1)
      }
    }

    const onSubmit = async formData => {
      const headersObj = {}
      formData.headers.forEach(item => {
        if (item.key) {
          headersObj[item.key] = item.value === '' || item.value === null ? null : item.value
        }
      })

      const queryArgsObj = {}
      formData.query_args.forEach(item => {
        if (item.key) {
          const value = item.value === '' || item.value === null ? null : item.value
          if (typeof value === 'string' && !isNaN(value) && value !== '') {
            queryArgsObj[item.key] = Number(value)
          } else {
            queryArgsObj[item.key] = value
          }
        }
      })

      const uriRewrite = formData.uri_rewrite.type && formData.uri_rewrite.value
        ? {
            type: formData.uri_rewrite.type,
            value: formData.uri_rewrite.value
          }
        : undefined

      const config = {
        enabled: formData.enabled ?? true
      }

      if (uriRewrite) {
        config.uri_rewrite = uriRewrite
      }
      if (Object.keys(headersObj).length > 0) {
        config.headers = headersObj
      }
      if (Object.keys(queryArgsObj).length > 0) {
        config.query_args = queryArgsObj
      }

      if (props.pluginConfigResId == null) {
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: formData.name ?? '',
          description: formData.description ?? '',
          enable: formData.enable == true ? 1 : 2,
          config: reactive(config)
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
        let configData = reactive({
          name: formData.name ?? '',
          description: formData.description ?? '',
          config: reactive(config)
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

    const cancel = async key => {
      if (props.pluginOpType == 1) {
        emit('pluginAddVisible')
        resetFields()
      } else {
        emit('pluginEditVisibleOff', key)
      }
    }

    const fn = reactive({
      addHeader,
      removeHeader,
      addQueryArg,
      removeQueryArg,
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginRequestRewrite }
  }
}
</script>


