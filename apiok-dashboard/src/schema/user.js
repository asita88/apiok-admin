const emailFormat = { type: 'email', message: '邮箱格式错误!' }

const passwordMin = { min: 8, message: '密码最小8个字符!' }
const passwordMax = { max: 16, message: '密码最长16个字符!' }

const nameFormat = {
  pattern: /^[A-Za-z\u4E00-\u9FA5\d_\-./*@#（）()&~+=[\]]+$/,
  message: '可包含汉字、字母、数字、_ + - . / * @ # （ ） ( ) [ ] & ~ ='
}
const nameMax = { max: 15, message: '昵称最长15个字符!' }

const userNameFormat = {
  pattern: /^[A-Za-z\u4E00-\u9FA5\d_\-./*@#（）()&~+=[\]]+$/,
  message: '可包含汉字、字母、数字、_ + - . / * @ # （ ） ( ) [ ] & ~ ='
}

export const schemaUser = {
  email: [{ required: true, message: '请输入邮箱!' }, emailFormat],
  password: [{ required: true, message: '请输入密码!' }, passwordMin, passwordMax],
  userName: [{ required: true, message: '请输入用户名!' }, userNameFormat, nameMax],
  oldPassword: [{ required: true, message: '请输入旧密码!' }, passwordMin, passwordMax],
  newPassword: [{ required: true, message: '请输入新密码!' }, passwordMin, passwordMax],
  rePassword: [
    { required: true, message: '请确认新密码!' },
    ({ getFieldValue }) => ({
      validator(_, value) {
        if (!value || getFieldValue('password') === value) {
          return Promise.resolve()
        }
        return Promise.reject(new Error('两次输入的密码不一致!'))
      }
    })
  ]
}
