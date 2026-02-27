<template>
  <a-form
    :model="data.formData"
    name="formData"
    :label-col="{ span: 5 }"
    :wrapper-col="{ span: 18 }"
    autocomplete="off"
    @finish="fn.onSubmit"
  >
    <a-form-item label="配置名称" name="name" :rules="schemaPluginResponseRewrite.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="响应头" name="headers">
      <div
        class="headers-list"
        v-for="(item, index) in data.formData.headers"
        :key="item.id"
      >
        <a-form-item
          :name="['headers', index, 'key']"
          style="width: 30%"
        >
          <a-input placeholder="header key" v-model:value="item.key" />
        </a-form-item>
        <a-form-item
          :name="['headers', index, 'value']"
          style="width: 50%"
        >
          <a-input placeholder="header value (留空表示删除)" v-model:value="item.value" />
        </a-form-item>
        <a-form-item>
          <a @click="fn.addHeader()">
            <i class="iconfont icon-tianjia"></i>
          </a>
          <a v-if="index > 0" @click="fn.removeHeader(item)">
            <i class="iconfont color-red icon-jian"></i>
          </a>
        </a-form-item>
      </div>
    </a-form-item>

    <a-form-item label="状态码" name="status_code">
      <a-input-number
        v-model:value="data.formData.status_code"
        :min="100"
        :max="599"
        style="width: 100%"
        placeholder="HTTP状态码(100-599)"
      />
    </a-form-item>

    <a-form-item label="响应体重写" name="body_rewrite">
      <a-card size="small">
        <a-form-item label="类型" name="body_rewrite.type">
          <a-select v-model:value="data.formData.body_rewrite.type" style="width: 100%">
            <a-select-option value="regex">regex</a-select-option>
            <a-select-option value="replace">replace</a-select-option>
            <a-select-option value="prefix">prefix</a-select-option>
            <a-select-option value="suffix">suffix</a-select-option>
          </a-select>
        </a-form-item>
        <template v-if="data.formData.body_rewrite.type === 'regex'">
          <a-form-item label="pattern" name="body_rewrite.value.pattern">
            <a-input v-model:value="data.formData.body_rewrite.value.pattern" placeholder="正则表达式模式" />
          </a-form-item>
          <a-form-item label="replacement" name="body_rewrite.value.replacement">
            <a-input v-model:value="data.formData.body_rewrite.value.replacement" placeholder="替换值" />
          </a-form-item>
          <a-form-item label="flags" name="body_rewrite.value.flags">
            <a-input v-model:value="data.formData.body_rewrite.value.flags" placeholder="正则标志(可选)" />
          </a-form-item>
        </template>
        <template v-else-if="data.formData.body_rewrite.type === 'replace'">
          <a-form-item label="from" name="body_rewrite.value.from">
            <a-input v-model:value="data.formData.body_rewrite.value.from" placeholder="要替换的字符串" />
          </a-form-item>
          <a-form-item label="to" name="body_rewrite.value.to">
            <a-input v-model:value="data.formData.body_rewrite.value.to" placeholder="替换为" />
          </a-form-item>
        </template>
        <template v-else-if="data.formData.body_rewrite.type === 'prefix' || data.formData.body_rewrite.type === 'suffix'">
          <a-form-item label="remove" name="body_rewrite.value.remove">
            <a-input v-model:value="data.formData.body_rewrite.value.remove" placeholder="要移除的前缀/后缀" />
          </a-form-item>
          <a-form-item label="add" name="body_rewrite.value.add">
            <a-input v-model:value="data.formData.body_rewrite.value.add" placeholder="要添加的前缀/后缀" />
          </a-form-item>
        </template>
      </a-card>
    </a-form-item>

    <a-form-item label="启用" name="enable" v-show="pluginOpType === 1">
      <a-switch v-model:checked="data.formData.enable" size="small" />
    </a-form-item>

    <a-form-item :wrapper-col="{ offset: 10, span: 16 }">
      <a-button html-type="submit" type="primary">保存</a-button>
      <a-button style="margin-left: 20px" @click="fn.cancel(pluginConfigData?.key)">取消</a-button>
    </a-form-item>
  </a-form>
</template>
<script>
import { reactive, onMounted } from 'vue'
import { Form, message } from 'ant-design-vue'
import { schemaPluginResponseRewrite } from '@/schema'
import { $pluginConfigAdd, $pluginConfigUpdate } from '@/api'

const useForm = Form.useForm
export default {
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
        name: 'plugin-response-rewrite',
        enabled: true,
        headers: [],
        status_code: undefined,
        body_rewrite: {
          type: 'regex',
          value: {
            pattern: '',
            replacement: '',
            flags: ''
          }
        },
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    onMounted(() => {
      if (props.pluginConfigResId == null) {
        addHeader()
      }
    })

    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.enabled != null) {
        data.formData.enabled = props.pluginConfigData.enabled
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
      if (props.pluginConfigData.status_code != null) {
        data.formData.status_code = props.pluginConfigData.status_code
      }
      if (props.pluginConfigData.body_rewrite != null) {
        data.formData.body_rewrite = JSON.parse(JSON.stringify(props.pluginConfigData.body_rewrite))
        if (!data.formData.body_rewrite.value) {
          data.formData.body_rewrite.value = {}
        }
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

    const onSubmit = async formData => {
      const headersObj = {}
      formData.headers.forEach(item => {
        if (item.key) {
          headersObj[item.key] = item.value === '' || item.value === null ? null : item.value
        }
      })

      const config = {
        enabled: formData.enabled ?? true
      }

      if (Object.keys(headersObj).length > 0) {
        config.headers = headersObj
      }
      if (formData.status_code != null) {
        config.status_code = formData.status_code
      }
      if (formData.body_rewrite.type && formData.body_rewrite.value) {
        config.body_rewrite = {
          type: formData.body_rewrite.type,
          value: formData.body_rewrite.value
        }
      }

      if (props.pluginConfigResId == null) {
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: formData.name ?? '',
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
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginResponseRewrite }
  }
}
</script>

<style scoped>
.headers-list {
  display: flex;
  margin-bottom: 10px;
}
</style>

