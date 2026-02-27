import request from '@/utils/requestAxios'

const logList = '/admin/log/list'
const accessLogList = '/admin/log/access/list'

/**
 * 日志列表
 * @param {*} params
 */
export const $logList = params => {
  return request.request({
    url: logList,
    method: 'GET',
    params: params
  })
}

/**
 * 访问日志列表
 * @param {*} params
 */
export const $accessLogList = params => {
  return request.request({
    url: accessLogList,
    method: 'GET',
    params: params
  })
}

const accessLogAggregation = '/admin/log/access/aggregation'

/**
 * 访问日志聚合统计
 * @param {*} params
 */
export const $accessLogAggregation = params => {
  return request.request({
    url: accessLogAggregation,
    method: 'GET',
    params: params
  })
}

const fieldAggregation = '/admin/log/access/field/aggregation'

/**
 * 字段聚合统计
 * @param {*} params
 */
export const $fieldAggregation = params => {
  return request.request({
    url: fieldAggregation,
    method: 'GET',
    params: params
  })
}

