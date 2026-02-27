<template>
  <div class="main">
    <a-form
      :model="formData"
      class="form"
      @finish="onSubmit"
      v-bind="{ labelCol: { xs: { span: 4 } }, wrapperCol: { xs: { span: 18 } } }"
    >
      <a-form-item label="用户名：" name="name" :rules="schemaUser.userName">
        <a-input v-model:value="formData.name" placeholder="请输入用户名" />
      </a-form-item>

      <a-form-item label="邮箱：" name="email" :rules="schemaUser.email">
        <a-input v-model:value="formData.email" placeholder="请输入邮箱" />
      </a-form-item>

      <a-form-item
        label="密码："
        name="password"
        :rules="isEdit ? schemaUser.passwordEdit : schemaUser.password"
      >
        <a-input-password
          v-model:value="formData.password"
          :placeholder="isEdit ? '留空则不修改密码' : '请输入密码（至少8个字符）'"
        />
      </a-form-item>

      <a-form-item :wrapper-col="{ offset: 4, span: 18 }">
        <a-button type="primary" html-type="submit" :loading="loading">提交</a-button>
        <a-button style="margin-left: 10px" @click="handleCancel">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { $userAdd, $userUpdate, $userInfo } from '@/api'
import { schemaUser } from '@/schema'

export default {
  props: {
    resId: {
      type: String,
      default: null
    }
  },
  emits: ['close', 'refresh'],
  setup(props, { emit }) {
    const loading = ref(false)
    const isEdit = ref(!!props.resId)
    const formData = reactive({
      name: '',
      email: '',
      password: ''
    })

    const schemaUserPasswordEdit = [
      { min: 8, message: '密码最小8个字符!' },
      { max: 16, message: '密码最长16个字符!' }
    ]

    const loadUserInfo = async () => {
      if (!props.resId) return

      const { code, msg, data } = await $userInfo(props.resId)
      if (code === 0) {
        formData.name = data.name || ''
        formData.email = data.email || ''
      } else {
        message.error(msg || '获取用户信息失败')
      }
    }

    const onSubmit = async () => {
      if (!isEdit.value && !formData.password) {
        message.error('密码不能为空')
        return
      }

      loading.value = true
      try {
        const submitData = {
          name: formData.name,
          email: formData.email
        }

        if (formData.password) {
          submitData.password = formData.password
        }

        let result
        if (isEdit.value) {
          result = await $userUpdate(props.resId, submitData)
        } else {
          result = await $userAdd(submitData)
        }

        if (result.code === 0) {
          message.success(isEdit.value ? '更新成功' : '添加成功')
          emit('refresh')
          emit('close')
        } else {
          message.error(result.msg || (isEdit.value ? '更新失败' : '添加失败'))
        }
      } catch (error) {
        message.error(isEdit.value ? '更新失败，请稍后重试' : '添加失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }

    const handleCancel = () => {
      emit('close')
    }

    onMounted(() => {
      loadUserInfo()
    })

    return {
      formData,
      onSubmit,
      handleCancel,
      loading,
      isEdit,
      schemaUser: {
        ...schemaUser,
        passwordEdit: schemaUserPasswordEdit
      }
    }
  }
}
</script>

<style lange="scss" scoped>
.main {
  padding: 20px 0;
}

.form {
  max-width: 600px;
}
</style>

