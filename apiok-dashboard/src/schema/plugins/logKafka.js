const brokers = {
  required: true,
  type: 'array',
  min: 1,
  message: '至少需要一个Kafka broker地址'
}

const topic = {
  required: true,
  message: '请输入Kafka topic名称'
}

const timeout = {
  type: 'number',
  min: 1000,
  max: 60000,
  message: '超时时间必须在1000-60000毫秒之间'
}

const keepaliveTimeout = {
  type: 'number',
  min: 1000,
  max: 600000,
  message: 'Keepalive超时时间必须在1000-600000毫秒之间'
}

export const schemaPluginLogKafka = {
  name: [{ required: true, message: '请输入插件名称' }],
  brokers: [brokers],
  topic: [topic],
  timeout: [timeout],
  keepalive_timeout: [keepaliveTimeout]
}

