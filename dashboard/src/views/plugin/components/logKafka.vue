<template>
  <a-form
    :model="data.formData"
    name="formData"
    :label-col="{ span: 5 }"
    :wrapper-col="{ span: 18 }"
    autocomplete="off"
    @finish="fn.onSubmit"
  >
    <a-form-item label="配置名称" name="name" :rules="schemaPluginLogKafka.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="brokers" name="brokers" :rules="schemaPluginLogKafka.brokers">
      <a-select
        mode="tags"
        v-model:value="data.formData.brokers"
        placeholder="输入Kafka broker地址，如: 127.0.0.1:9092"
        style="width: 100%"
      />
    </a-form-item>

    <a-form-item label="topic" name="topic" :rules="schemaPluginLogKafka.topic">
      <a-input v-model:value="data.formData.topic" placeholder="Kafka topic名称" />
    </a-form-item>

    <a-form-item label="timeout" name="timeout" :rules="schemaPluginLogKafka.timeout">
      <a-input-number
        v-model:value="data.formData.timeout"
        :min="1000"
        :max="60000"
        style="width: 100%"
        placeholder="连接超时时间(毫秒)"
      />
    </a-form-item>

    <a-form-item label="keepalive_timeout" name="keepalive_timeout" :rules="schemaPluginLogKafka.keepalive_timeout">
      <a-input-number
        v-model:value="data.formData.keepalive_timeout"
        :min="1000"
        :max="600000"
        style="width: 100%"
        placeholder="Keepalive超时时间(毫秒)"
      />
    </a-form-item>

    <a-form-item label="include_request_body" name="include_request_body">
      <a-switch v-model:checked="data.formData.include_request_body" size="small" />
    </a-form-item>

    <a-form-item label="include_response_body" name="include_response_body">
      <a-switch v-model:checked="data.formData.include_response_body" size="small" />
    </a-form-item>

    <a-form-item label="include_headers" name="include_headers">
      <a-select
        mode="tags"
        v-model:value="data.formData.include_headers"
        placeholder="输入要包含的请求/响应头名称"
        style="width: 100%"
      />
    </a-form-item>

    <a-form-item label="exclude_headers" name="exclude_headers">
      <a-select
        mode="tags"
        v-model:value="data.formData.exclude_headers"
        placeholder="输入要排除的请求头名称(不区分大小写)"
        style="width: 100%"
      />
    </a-form-item>

    <a-form-item label="log_format" name="log_format">
      <a-select v-model:value="data.formData.log_format" style="width: 100%">
        <a-select-option value="json">json</a-select-option>
        <a-select-option value="text">text</a-select-option>
      </a-select>
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
import { reactive } from 'vue'
import { Form, message } from 'ant-design-vue'
import { schemaPluginLogKafka } from '@/schema'
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
        name: 'plugin-log-kafka',
        enabled: true,
        brokers: [],
        topic: '',
        timeout: 5000,
        keepalive_timeout: 60000,
        include_request_body: false,
        include_response_body: false,
        include_headers: [],
        exclude_headers: [],
        log_format: 'json',
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.enabled != null) {
        data.formData.enabled = props.pluginConfigData.enabled
      }
      if (props.pluginConfigData.brokers != null) {
        data.formData.brokers = Array.isArray(props.pluginConfigData.brokers)
          ? props.pluginConfigData.brokers
          : []
      }
      if (props.pluginConfigData.topic != null) {
        data.formData.topic = props.pluginConfigData.topic
      }
      if (props.pluginConfigData.timeout != null) {
        data.formData.timeout = props.pluginConfigData.timeout
      }
      if (props.pluginConfigData.keepalive_timeout != null) {
        data.formData.keepalive_timeout = props.pluginConfigData.keepalive_timeout
      }
      if (props.pluginConfigData.include_request_body != null) {
        data.formData.include_request_body = props.pluginConfigData.include_request_body
      }
      if (props.pluginConfigData.include_response_body != null) {
        data.formData.include_response_body = props.pluginConfigData.include_response_body
      }
      if (props.pluginConfigData.include_headers != null) {
        data.formData.include_headers = Array.isArray(props.pluginConfigData.include_headers)
          ? props.pluginConfigData.include_headers
          : []
      }
      if (props.pluginConfigData.exclude_headers != null) {
        data.formData.exclude_headers = Array.isArray(props.pluginConfigData.exclude_headers)
          ? props.pluginConfigData.exclude_headers
          : []
      }
      if (props.pluginConfigData.log_format != null) {
        data.formData.log_format = props.pluginConfigData.log_format
      }
    }

    const onSubmit = async formData => {
      if (props.pluginConfigResId == null) {
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: formData.name ?? '',
          enable: formData.enable == true ? 1 : 2,
          config: reactive({
            enabled: formData.enabled ?? true,
            brokers: formData.brokers ?? [],
            topic: formData.topic ?? '',
            timeout: formData.timeout ?? 5000,
            keepalive_timeout: formData.keepalive_timeout ?? 60000,
            include_request_body: formData.include_request_body ?? false,
            include_response_body: formData.include_response_body ?? false,
            include_headers: formData.include_headers ?? [],
            exclude_headers: formData.exclude_headers ?? [],
            log_format: formData.log_format ?? 'json'
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
          config: reactive({
            enabled: formData.enabled ?? true,
            brokers: formData.brokers ?? [],
            topic: formData.topic ?? '',
            timeout: formData.timeout ?? 5000,
            keepalive_timeout: formData.keepalive_timeout ?? 60000,
            include_request_body: formData.include_request_body ?? false,
            include_response_body: formData.include_response_body ?? false,
            include_headers: formData.include_headers ?? [],
            exclude_headers: formData.exclude_headers ?? [],
            log_format: formData.log_format ?? 'json'
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
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginLogKafka }
  }
}
</script>

