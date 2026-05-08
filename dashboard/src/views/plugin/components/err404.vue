<template>
  <div class="plugin-form-empty">
    <p v-show="data.pluginType === 1" class="plugin-form-empty-title">请选择插件</p>
    <p v-show="data.pluginType === 1" class="plugin-form-empty-desc">请先在「插件信息」中选择要配置的插件类型。</p>
    <p v-show="data.pluginType === 2" class="plugin-form-empty-title">暂不支持</p>
    <p v-show="data.pluginType === 2" class="plugin-form-empty-desc">当前插件类型未提供配置表单。</p>
  </div>
</template>
<script>
import { reactive } from 'vue'

export default {
  props: {
    pluginOpType: {
      Number
    },
    pluginTag: {
      String
    }
  },
  setup(props) {
    const data = reactive({
      pluginType: 0
    })

    if (props.pluginTag !== '') {
      data.pluginType = 2
    } else {
      data.pluginType = 1
    }

    return { data }
  }
}
</script>
