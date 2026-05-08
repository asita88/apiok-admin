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
    <a-form-item label="配置名称" name="name" :rules="schemaPluginTrafficTag.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="插件描述" name="description">
      <a-textarea v-model:value="data.formData.description" :rows="2" placeholder="请输入插件配置描述" />
    </a-form-item>

    <a-form-item label="匹配规则" name="match_rules">
      <a-card class="plugin-form-nested-card" size="small">
        <a-form-item label="path" name="match_rules.path">
          <a-input v-model:value="data.formData.match_rules.path" placeholder="路径匹配，如: /api/users" />
        </a-form-item>
        <a-form-item label="method" name="match_rules.method">
          <a-select
            mode="multiple"
            v-model:value="data.formData.match_rules.method"
            placeholder="选择HTTP方法"
            style="width: 100%"
          >
            <a-select-option value="GET">GET</a-select-option>
            <a-select-option value="POST">POST</a-select-option>
            <a-select-option value="PUT">PUT</a-select-option>
            <a-select-option value="DELETE">DELETE</a-select-option>
            <a-select-option value="PATCH">PATCH</a-select-option>
            <a-select-option value="OPTIONS">OPTIONS</a-select-option>
            <a-select-option value="HEAD">HEAD</a-select-option>
            <a-select-option value="TRACE">TRACE</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="headers" name="match_rules.headers">
          <div
            class="plugin-form-kv-row"
            v-for="(item, index) in data.formData.match_rules.headers"
            :key="item.id"
          >
            <a-form-item class="plugin-form-kv-field" :name="['match_rules', 'headers', index, 'key']">
              <a-input placeholder="header key" v-model:value="item.key" />
            </a-form-item>
            <a-form-item class="plugin-form-kv-field" :name="['match_rules', 'headers', index, 'value']">
              <a-input placeholder="header value" v-model:value="item.value" />
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
      </a-card>
    </a-form-item>

    <a-form-item label="标签" name="tags" :rules="schemaPluginTrafficTag.tags">
      <div class="plugin-form-kv-row" v-for="(item, index) in data.formData.tags" :key="item.id">
        <a-form-item
          class="plugin-form-kv-field"
          :name="['tags', index, 'key']"
          :rules="[{ required: true, message: '请输入标签key' }]"
        >
          <a-input placeholder="标签key" v-model:value="item.key" />
        </a-form-item>
        <a-form-item
          class="plugin-form-kv-field"
          :name="['tags', index, 'value']"
          :rules="[{ required: true, message: '请输入标签value' }]"
        >
          <a-input placeholder="标签value" v-model:value="item.value" />
        </a-form-item>
        <div class="plugin-form-kv-actions">
          <a-button type="link" size="small" @click="fn.addTag">
            <template #icon>
              <PlusOutlined />
            </template>
          </a-button>
          <a-button
            v-if="index > 0"
            type="link"
            size="small"
            danger
            @click="fn.removeTag(item)"
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
import { schemaPluginTrafficTag } from '@/schema'
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
        name: 'plugin-traffic-tag',
        description: '',
        match_rules: {
          path: '',
          method: [],
          headers: []
        },
        tags: [],
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    onMounted(() => {
      if (
        props.pluginConfigResId == null ||
        !data.formData.tags.length
      ) {
        addTag()
      }
      if (
        props.pluginConfigResId == null ||
        !data.formData.match_rules.headers.length
      ) {
        addHeader()
      }
    })

    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.description != null) {
        data.formData.description = props.pluginConfigData.description
      }
      if (props.pluginConfigData.match_rules != null) {
        if (props.pluginConfigData.match_rules.path != null) {
          data.formData.match_rules.path = props.pluginConfigData.match_rules.path
        }
        if (props.pluginConfigData.match_rules.method != null) {
          data.formData.match_rules.method = Array.isArray(props.pluginConfigData.match_rules.method)
            ? props.pluginConfigData.match_rules.method
            : [props.pluginConfigData.match_rules.method]
        }
        if (props.pluginConfigData.match_rules.headers != null) {
          const headers = props.pluginConfigData.match_rules.headers
          Object.keys(headers).forEach(key => {
            data.formData.match_rules.headers.push({
              key: key,
              value: headers[key],
              id: Date.now() + Math.random()
            })
          })
        }
      }
      if (props.pluginConfigData.tags != null) {
        const tags = props.pluginConfigData.tags
        Object.keys(tags).forEach(key => {
          data.formData.tags.push({
            key: key,
            value: tags[key],
            id: Date.now() + Math.random()
          })
        })
      }
    }

    const addTag = () => {
      data.formData.tags.push({
        key: '',
        value: '',
        id: Date.now() + Math.random()
      })
    }

    const removeTag = item => {
      const index = data.formData.tags.indexOf(item)
      if (index !== -1) {
        data.formData.tags.splice(index, 1)
      }
    }

    const addHeader = () => {
      data.formData.match_rules.headers.push({
        key: '',
        value: '',
        id: Date.now() + Math.random()
      })
    }

    const removeHeader = item => {
      const index = data.formData.match_rules.headers.indexOf(item)
      if (index !== -1) {
        data.formData.match_rules.headers.splice(index, 1)
      }
    }

    const onSubmit = async formData => {
      const tagsObj = {}
      formData.tags.forEach(item => {
        if (item.key) {
          tagsObj[item.key] = item.value || ''
        }
      })

      const headersObj = {}
      formData.match_rules.headers.forEach(item => {
        if (item.key) {
          headersObj[item.key] = item.value || ''
        }
      })

      const matchRules = {
        path: formData.match_rules.path || undefined,
        method: formData.match_rules.method && formData.match_rules.method.length > 0
          ? formData.match_rules.method.length === 1
            ? formData.match_rules.method[0]
            : formData.match_rules.method
          : undefined,
        headers: Object.keys(headersObj).length > 0 ? headersObj : undefined
      }

      Object.keys(matchRules).forEach(key => {
        if (matchRules[key] === undefined) {
          delete matchRules[key]
        }
      })

      if (props.pluginConfigResId == null) {
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: formData.name ?? '',
          description: formData.description ?? '',
          enable: formData.enable == true ? 1 : 2,
          config: reactive({
            match_rules: matchRules,
            tags: tagsObj
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
        let configData = reactive({
          name: formData.name ?? '',
          description: formData.description ?? '',
          config: reactive({
            match_rules: matchRules,
            tags: tagsObj
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

    const cancel = async key => {
      if (props.pluginOpType == 1) {
        emit('pluginAddVisible')
        resetFields()
      } else {
        emit('pluginEditVisibleOff', key)
      }
    }

    const fn = reactive({
      addTag,
      removeTag,
      addHeader,
      removeHeader,
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginTrafficTag }
  }
}
</script>


