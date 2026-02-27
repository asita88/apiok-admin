<template>
  <a-form
    :model="data.formData"
    name="formData"
    :label-col="{ span: 5 }"
    :wrapper-col="{ span: 18 }"
    autocomplete="off"
    @finish="fn.onSubmit"
  >
    <a-form-item label="配置名称" name="name" :rules="schemaPluginWaf.name">
      <a-input v-model:value="data.formData.name" />
    </a-form-item>

    <a-form-item label="启用" name="enabled">
      <a-switch v-model:checked="data.formData.enabled" size="small" />
    </a-form-item>

    <a-divider>IP白名单</a-divider>
    <a-form-item label="启用白名单" name="ip_whitelist.enabled">
      <a-switch v-model:checked="data.formData.ip_whitelist.enabled" size="small" />
    </a-form-item>
    <a-form-item
      label="IP列表"
      name="ip_whitelist.ip_list"
      :rules="schemaPluginWaf['ip_whitelist.ip_list']"
    >
      <a-select
        mode="tags"
        v-model:value="data.formData.ip_whitelist.ip_list"
        placeholder="输入IP地址，支持通配符，如: 192.168.1.*"
        style="width: 100%"
      />
    </a-form-item>

    <a-divider>IP黑名单</a-divider>
    <a-form-item label="启用黑名单" name="ip_blacklist.enabled">
      <a-switch v-model:checked="data.formData.ip_blacklist.enabled" size="small" />
    </a-form-item>
    <a-form-item
      label="IP列表"
      name="ip_blacklist.ip_list"
      :rules="schemaPluginWaf['ip_blacklist.ip_list']"
    >
      <a-select
        mode="tags"
        v-model:value="data.formData.ip_blacklist.ip_list"
        placeholder="输入IP地址，支持通配符，如: 192.168.1.*"
        style="width: 100%"
      />
    </a-form-item>

    <a-divider>规则配置</a-divider>
    <a-form-item label="规则列表" name="rules.rule_list">
      <a-card
        v-for="(rule, ruleIndex) in data.formData.rules.rule_list"
        :key="rule.id"
        size="small"
        style="margin-bottom: 16px"
      >
        <template #title>
          <a-input
            v-model:value="rule.name"
            placeholder="规则名称"
            style="width: 200px"
          />
          <a-button
            type="link"
            danger
            size="small"
            style="float: right"
            @click="fn.removeRule(rule)"
          >
            删除规则
          </a-button>
        </template>

        <a-form-item label="动作" :name="['rules', 'rule_list', ruleIndex, 'action']">
          <a-select v-model:value="rule.action" style="width: 100%">
            <a-select-option value="block">block</a-select-option>
            <a-select-option value="log">log</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="条件列表">
          <a-card
            v-for="(condition, condIndex) in rule.conditions"
            :key="condition.id"
            size="small"
            style="margin-bottom: 8px"
          >
            <template #title>
              <span>条件 {{ condIndex + 1 }}</span>
              <a-button
                type="link"
                danger
                size="small"
                style="float: right"
                @click="fn.removeCondition(rule, condition)"
              >
                删除
              </a-button>
            </template>

            <a-form-item
              label="匹配类型"
              :name="['rules', 'rule_list', ruleIndex, 'conditions', condIndex, 'match_type']"
            >
              <a-select v-model:value="condition.match_type" style="width: 100%">
                <a-select-option value="uri">uri</a-select-option>
                <a-select-option value="args">args</a-select-option>
                <a-select-option value="header">header</a-select-option>
                <a-select-option value="body">body</a-select-option>
                <a-select-option value="all">all</a-select-option>
                <a-select-option value="method">method</a-select-option>
                <a-select-option value="request_size">request_size</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item
              label="匹配模式"
              :name="['rules', 'rule_list', ruleIndex, 'conditions', condIndex, 'patterns']"
            >
              <a-select
                mode="tags"
                v-model:value="condition.patterns"
                placeholder="输入匹配模式（正则表达式或精确匹配）"
                style="width: 100%"
              />
            </a-form-item>

            <a-form-item
              label="操作符"
              :name="['rules', 'rule_list', ruleIndex, 'conditions', condIndex, 'operator']"
            >
              <a-select v-model:value="condition.operator" style="width: 100%">
                <a-select-option value="match">match</a-select-option>
                <a-select-option value="not_match">not_match</a-select-option>
              </a-select>
            </a-form-item>
          </a-card>
          <a-button type="dashed" block @click="fn.addCondition(rule)">
            <i class="iconfont icon-tianjia"></i> 添加条件
          </a-button>
        </a-form-item>
      </a-card>
      <a-button type="dashed" block @click="fn.addRule">
        <i class="iconfont icon-tianjia"></i> 添加规则
      </a-button>
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
import { schemaPluginWaf } from '@/schema'
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
        name: 'plugin-waf',
        enabled: true,
        ip_whitelist: {
          enabled: true,
          ip_list: []
        },
        ip_blacklist: {
          enabled: true,
          ip_list: []
        },
        rules: {
          rule_list: []
        },
        enable: false
      }
    })

    const { resetFields } = useForm(data.formData)

    onMounted(() => {
      if (props.pluginConfigResId == null) {
        addRule()
      }
    })

    if (props.pluginConfigData != null) {
      if (props.pluginConfigData.name != null) {
        data.formData.name = props.pluginConfigData.name
      }
      if (props.pluginConfigData.enabled != null) {
        data.formData.enabled = props.pluginConfigData.enabled
      }
      if (props.pluginConfigData.ip_whitelist != null) {
        data.formData.ip_whitelist = JSON.parse(JSON.stringify(props.pluginConfigData.ip_whitelist))
        if (!data.formData.ip_whitelist.ip_list || !Array.isArray(data.formData.ip_whitelist.ip_list)) {
          data.formData.ip_whitelist.ip_list = []
        }
        if (data.formData.ip_whitelist.enabled === undefined) {
          data.formData.ip_whitelist.enabled = true
        }
      }
      if (props.pluginConfigData.ip_blacklist != null) {
        data.formData.ip_blacklist = JSON.parse(JSON.stringify(props.pluginConfigData.ip_blacklist))
        if (!data.formData.ip_blacklist.ip_list || !Array.isArray(data.formData.ip_blacklist.ip_list)) {
          data.formData.ip_blacklist.ip_list = []
        }
        if (data.formData.ip_blacklist.enabled === undefined) {
          data.formData.ip_blacklist.enabled = true
        }
      }
      if (props.pluginConfigData.rules != null) {
        data.formData.rules = JSON.parse(JSON.stringify(props.pluginConfigData.rules))
        if (!data.formData.rules.rule_list) {
          data.formData.rules.rule_list = []
        } else {
          data.formData.rules.rule_list.forEach(rule => {
            rule.id = rule.id || Date.now() + Math.random()
            if (!rule.conditions) {
              rule.conditions = []
            } else {
              rule.conditions.forEach(cond => {
                cond.id = cond.id || Date.now() + Math.random()
                if (!cond.patterns || !Array.isArray(cond.patterns)) {
                  cond.patterns = []
                }
                if (!cond.operator) {
                  cond.operator = 'match'
                }
              })
            }
          })
        }
      }
    }

    const addRule = () => {
      data.formData.rules.rule_list.push({
        id: Date.now() + Math.random(),
        name: '',
        action: 'block',
        conditions: []
      })
      const newRule = data.formData.rules.rule_list[data.formData.rules.rule_list.length - 1]
      addCondition(newRule)
    }

    const removeRule = rule => {
      const index = data.formData.rules.rule_list.indexOf(rule)
      if (index !== -1) {
        data.formData.rules.rule_list.splice(index, 1)
      }
    }

    const addCondition = rule => {
      rule.conditions.push({
        id: Date.now() + Math.random(),
        match_type: 'all',
        patterns: [],
        operator: 'match'
      })
    }

    const removeCondition = (rule, condition) => {
      const index = rule.conditions.indexOf(condition)
      if (index !== -1) {
        rule.conditions.splice(index, 1)
      }
    }

    const onSubmit = async formData => {
      const config = {
        enabled: data.formData.enabled ?? true,
        ip_whitelist: {
          enabled: data.formData.ip_whitelist?.enabled ?? true,
          ip_list: Array.isArray(data.formData.ip_whitelist?.ip_list) ? data.formData.ip_whitelist.ip_list : []
        },
        ip_blacklist: {
          enabled: data.formData.ip_blacklist?.enabled ?? true,
          ip_list: Array.isArray(data.formData.ip_blacklist?.ip_list) ? data.formData.ip_blacklist.ip_list : []
        },
        rules: {
          rule_list: (data.formData.rules?.rule_list || []).map(rule => ({
            name: rule.name,
            action: rule.action || 'block',
            conditions: (rule.conditions || []).map(cond => ({
              match_type: cond.match_type,
              patterns: Array.isArray(cond.patterns) ? cond.patterns : [],
              operator: cond.operator || 'match'
            }))
          }))
        }
      }

      if (props.pluginConfigResId == null) {
        let configData = reactive({
          plugin_id: props.pluginResId ?? '',
          target_id: props.targetResId ?? '',
          type: props.pluginConfigType ?? '',
          name: data.formData.name ?? '',
          enable: data.formData.enable == true ? 1 : 2,
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
          name: data.formData.name ?? '',
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
      addRule,
      removeRule,
      addCondition,
      removeCondition,
      onSubmit,
      cancel
    })

    return { data, fn, schemaPluginWaf }
  }
}
</script>

<style scoped>
.ant-card {
  margin-bottom: 16px;
}
</style>
