<template>
  <div class="change-password">
    <a-card title="修改密码" :bordered="false">
      <a-form
        :model="formData"
        class="form"
        @finish="onSubmit"
        v-bind="{ labelCol: { xs: { span: 4 } }, wrapperCol: { xs: { span: 18 } } }"
      >
        <a-form-item label="旧密码：" name="old_password" :rules="schemaUser.oldPassword">
          <a-input-password v-model:value="formData.old_password" placeholder="请输入旧密码" />
        </a-form-item>

        <a-form-item label="新密码：" name="password" :rules="schemaUser.newPassword">
          <a-input-password v-model:value="formData.password" placeholder="请输入新密码（至少8个字符）" />
        </a-form-item>

        <a-form-item label="确认新密码：" name="re_password" :rules="schemaUser.rePassword">
          <a-input-password v-model:value="formData.re_password" placeholder="请再次输入新密码" />
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 4, span: 18 }">
          <a-button type="primary" html-type="submit" :loading="loading">提交</a-button>
          <a-button style="margin-left: 10px" @click="resetForm">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script>
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { $changePassword } from '@/api'
import { schemaUser } from '@/schema'

export default {
  setup() {
    const loading = ref(false)
    const formData = reactive({
      old_password: '',
      password: '',
      re_password: ''
    })

    const onSubmit = async () => {
      loading.value = true
      try {
        const { code, msg } = await $changePassword({
          old_password: formData.old_password,
          password: formData.password,
          re_password: formData.re_password
        })
        if (code === 0) {
          message.success('密码修改成功')
          resetForm()
        } else {
          message.error(msg || '密码修改失败')
        }
      } catch (error) {
        message.error('密码修改失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }

    const resetForm = () => {
      formData.old_password = ''
      formData.password = ''
      formData.re_password = ''
    }

    return {
      formData,
      onSubmit,
      resetForm,
      loading,
      schemaUser
    }
  }
}
</script>

<style lange="scss" scoped>
.change-password {
  padding: 20px;
}

.form {
  max-width: 600px;
  margin-top: 20px;
}
</style>

