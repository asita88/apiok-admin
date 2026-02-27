<template>
  <a-form
    :model="data.formData"
    name="formData"
    :label-col="{ span: 5 }"
    :wrapper-col="{ span: 18 }"
    autocomplete="off"
    @finish="fn.onSubmit"
  >
    <a-form-item label="配置名称" name="name" :rules="schemaPluginLogMysql.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="host" name="host" :rules="schemaPluginLogMysql.host">
      <a-input v-model:value="data.formData.host" placeholder="MySQL主机地址" />
    </a-form-item>

    <a-form-item label="port" name="port" :rules="schemaPluginLogMysql.port">
      <a-input-number
        v-model:value="data.formData.port"
        :min="1"
        :max="65535"
        style="width: 100%"
        placeholder="MySQL端口"
      />
    </a-form-item>

    <a-form-item label="database" name="database" :rules="schemaPluginLogMysql.database">
      <a-input v-model:value="data.formData.database" placeholder="数据库名称" />
    </a-form-item>

    <a-form-item label="user" name="user" :rules="schemaPluginLogMysql.user">
      <a-input v-model:value="data.formData.user" placeholder="数据库用户名" />
    </a-form-item>

    <a-form-item label="password" name="password">
      <a-input-password v-model:value="data.formData.password" placeholder="数据库密码" />
    </a-form-item>

    <a-form-item label="table_name" name="table_name" :rules="schemaPluginLogMysql.table_name">
      <a-input v-model:value="data.formData.table_name" placeholder="表名" />
    </a-form-item>

    <a-form-item label="timeout" name="timeout" :rules="schemaPluginLogMysql.timeout">
      <a-input-number
        v-model:value="data.formData.timeout"
        :min="1000"
        :max="60000"
        style="width: 100%"
        placeholder="连接超时时间(毫秒)"
      />
    </a-form-item>

    <a-form-item label="pool_size" name="pool_size" :rules="schemaPluginLogMysql.pool_size">
      <a-input-number
        v-model:value="data.formData.pool_size"
        :min="1"
        :max="1000"
        style="width: 100%"
        placeholder="连接池大小"
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

    <a-form-item label="batch_size" name="batch_size" :rules="schemaPluginLogMysql.batch_size">
      <a-input-number
        v-model:value="data.formData.batch_size"
        :min="1"
        :max="1000"
        style="width: 100%"
        placeholder="批量插入大小"
      />
    </a-form-item>

    <a-form-item label="batch_timeout" name="batch_timeout" :rules="schemaPluginLogMysql.batch_timeout">
      <a-input-number
        v-model:value="data.formData.batch_timeout"
        :min="1000"
        :max="60000"
        style="width: 100%"
        placeholder="批量插入超时时间(毫秒)"
      />
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
import { schemaPluginLogMysql } from '@/schema'
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
        name: 'plugin-log-mysql',
        enabled: true,
        host: '127.0.0.1',
        port: 3306,
        database: '',
        user: '',
        password: '',
        table_name: 'apiok_access_log',
        timeout: 5000,
        pool_size: 100,
        include_request_body: false,
        include_response_body: false,
        include_headers: [],
        exclude_headers: [],
        batch_size: 100,
        batch_timeout: 5000,
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
      if (props.pluginConfigData.host != null) {
        data.formData.host = props.pluginConfigData.host
      }
      if (props.pluginConfigData.port != null) {
        data.formData.port = props.pluginConfigData.port
      }
      if (props.pluginConfigData.database != null) {
        data.formData.database = props.pluginConfigData.database
      }
      if (props.pluginConfigData.user != null) {
        data.formData.user = props.pluginConfigData.user
      }
      if (props.pluginConfigData.password != null) {
        data.formData.password = props.pluginConfigData.password
      }
      if (props.pluginConfigData.table_name != null) {
        data.formData.table_name = props.pluginConfigData.table_name
      }
      if (props.pluginConfigData.timeout != null) {
        data.formData.timeout = props.pluginConfigData.timeout
      }
      if (props.pluginConfigData.pool_size != null) {
        data.formData.pool_size = props.pluginConfigData.pool_size
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
      if (props.pluginConfigData.batch_size != null) {
        data.formData.batch_size = props.pluginConfigData.batch_size
      }
      if (props.pluginConfigData.batch_timeout != null) {
        data.formData.batch_timeout = props.pluginConfigData.batch_timeout
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
            host: formData.host ?? '127.0.0.1',
            port: formData.port ?? 3306,
            database: formData.database ?? '',
            user: formData.user ?? '',
            password: formData.password ?? '',
            table_name: formData.table_name ?? 'apiok_access_log',
            timeout: formData.timeout ?? 5000,
            pool_size: formData.pool_size ?? 100,
            include_request_body: formData.include_request_body ?? false,
            include_response_body: formData.include_response_body ?? false,
            include_headers: formData.include_headers ?? [],
            exclude_headers: formData.exclude_headers ?? [],
            batch_size: formData.batch_size ?? 100,
            batch_timeout: formData.batch_timeout ?? 5000
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
            host: formData.host ?? '127.0.0.1',
            port: formData.port ?? 3306,
            database: formData.database ?? '',
            user: formData.user ?? '',
            password: formData.password ?? '',
            table_name: formData.table_name ?? 'apiok_access_log',
            timeout: formData.timeout ?? 5000,
            pool_size: formData.pool_size ?? 100,
            include_request_body: formData.include_request_body ?? false,
            include_response_body: formData.include_response_body ?? false,
            include_headers: formData.include_headers ?? [],
            exclude_headers: formData.exclude_headers ?? [],
            batch_size: formData.batch_size ?? 100,
            batch_timeout: formData.batch_timeout ?? 5000
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

    return { data, fn, schemaPluginLogMysql }
  }
}
</script>

