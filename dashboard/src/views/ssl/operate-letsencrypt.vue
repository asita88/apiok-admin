<template>
  <div class="main">
    <a-form
      :model="data.formData"
      class="form"
      :label-col="{ style: { width: '120px' } }"
      :wrapper-col="{ span: 18 }"
      @finish="fn.onSubmit"
    >
      <a-alert
        message="Let's Encrypt 证书申请"
        description="请输入要申请证书的域名，系统将自动通过 Let's Encrypt 申请免费 SSL 证书。"
        type="info"
        show-icon
        style="margin-bottom: 20px"
      />

      <a-form-item label="域名：" name="domain" :rules="schemaLetsEncrypt.domain">
        <a-input v-model:value="data.formData.domain" placeholder="例如：example.com" />
      </a-form-item>

      <a-form-item label="申请后启用：">
        <a-switch v-model:checked="data.formData.enable" />
        <span style="margin-left: 10px; color: #999">申请成功后自动启用证书</span>
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 10, offset: 16 }">
        <a-button type="primary" html-type="submit" :loading="data.loading">申请证书</a-button>
        <a-button style="margin-left: 15px" @click="fn.cancel">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import { $letsencryptRequest } from '@/api'
import { schemaLetsEncrypt } from '@/schema'

export default {
  emits: ['componentCloseDrawer', 'componentRefreshList'],
  setup(props, { emit }) {
    const data = reactive({
      formData: {
        domain: '',
        enable: false
      },
      loading: false
    })

    const onSubmit = async () => {
      if (!data.formData.domain || data.formData.domain.trim() === '') {
        message.error('请输入域名')
        return
      }

      data.loading = true
      const formData = {
        domain: data.formData.domain.trim(),
        enable: data.formData.enable ? 1 : 2
      }

      try {
        const res = await $letsencryptRequest(formData)

        if (res.code != 0) {
          message.error(res.msg)
          return
        } else {
          message.success('证书申请成功，请稍后刷新列表查看')
          emit('componentRefreshList')
        }
      } catch (error) {
        message.error('申请证书失败：' + (error.message || '未知错误'))
      } finally {
        data.loading = false
      }
    }

    const cancel = () => {
      emit('componentCloseDrawer')
    }

    const fn = reactive({
      onSubmit,
      cancel
    })

    return {
      schemaLetsEncrypt,
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
</style>

