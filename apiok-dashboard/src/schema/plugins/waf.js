import { nameFormat, nameMax } from '../common'

const ipList = {
  type: 'array',
  min: 1,
  message: '至少需要一个IP地址'
}

const ruleName = {
  required: true,
  message: '请输入规则名称'
}

const patterns = {
  type: 'array',
  min: 1,
  message: '至少需要一个匹配模式'
}

const matchType = {
  required: true,
  message: '请选择匹配类型'
}

const action = {
  required: true,
  message: '请选择动作'
}

export const schemaPluginWaf = {
  name: [nameFormat, nameMax],
  'ip_whitelist.ip_list': [ipList],
  'ip_blacklist.ip_list': [ipList],
  'rules.rule_list[].name': [ruleName],
  'rules.rule_list[].conditions[].patterns': [patterns],
  'rules.rule_list[].conditions[].match_type': [matchType],
  'rules.rule_list[].action': [action]
}
