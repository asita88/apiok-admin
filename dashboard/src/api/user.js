import request from '@/utils/requestAxios'

const login = '/admin/user/login'
const logout = '/admin/user/logout'

/**
 * 登录
 * @param {*} data
 */
export const $login = data => {
  return request.request({
    url: login,
    method: 'POST',
    data: data
  })
}

/**
 * 退出登录
 * @returns
 */
export const $logout = () => {
  return request.request({
    url: logout,
    method: 'POST'
  })
}

const changePassword = '/admin/user/change-password'

/**
 * 修改密码
 * @param {*} data
 */
export const $changePassword = data => {
  return request.request({
    url: changePassword,
    method: 'PUT',
    data: data
  })
}

const userList = '/admin/user/list'
const userInfo = '/admin/user/info'
const userAdd = '/admin/user/add'
const userUpdate = '/admin/user/update'
const userDelete = '/admin/user/delete'

/**
 * 用户列表
 * @param {*} params
 */
export const $userList = params => {
  return request.request({
    url: userList,
    method: 'GET',
    params: params
  })
}

/**
 * 用户详情
 * @param {*} resId
 */
export const $userInfo = resId => {
  return request.request({
    url: userInfo + '/' + resId,
    method: 'GET'
  })
}

/**
 * 添加用户
 * @param {*} data
 */
export const $userAdd = data => {
  return request.request({
    url: userAdd,
    method: 'POST',
    data: data
  })
}

/**
 * 更新用户
 * @param {*} resId
 * @param {*} data
 */
export const $userUpdate = (resId, data) => {
  return request.request({
    url: userUpdate + '/' + resId,
    method: 'PUT',
    data: data
  })
}

/**
 * 删除用户
 * @param {*} resId
 */
export const $userDelete = resId => {
  return request.request({
    url: userDelete + '/' + resId,
    method: 'DELETE'
  })
}
