const host = {
  required: true,
  message: '请输入MySQL主机地址'
}

const port = {
  type: 'number',
  min: 1,
  max: 65535,
  message: '端口号必须在1-65535之间'
}

const database = {
  required: true,
  message: '请输入数据库名称'
}

const user = {
  required: true,
  message: '请输入数据库用户名'
}

const tableName = {
  required: true,
  message: '请输入表名'
}

const timeout = {
  type: 'number',
  min: 1000,
  max: 60000,
  message: '超时时间必须在1000-60000毫秒之间'
}

const poolSize = {
  type: 'number',
  min: 1,
  max: 1000,
  message: '连接池大小必须在1-1000之间'
}

const batchSize = {
  type: 'number',
  min: 1,
  max: 1000,
  message: '批量插入大小必须在1-1000之间'
}

const batchTimeout = {
  type: 'number',
  min: 1000,
  max: 60000,
  message: '批量插入超时时间必须在1000-60000毫秒之间'
}

export const schemaPluginLogMysql = {
  name: [{ required: true, message: '请输入插件名称' }],
  host: [host],
  port: [port],
  database: [database],
  user: [user],
  table_name: [tableName],
  timeout: [timeout],
  pool_size: [poolSize],
  batch_size: [batchSize],
  batch_timeout: [batchTimeout]
}

