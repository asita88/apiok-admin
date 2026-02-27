<template>
  <div class="main">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item
        ><i
          style="color: #448ef7; font-size: 30px"
          class="iconfont icon-rizhi"
        />访问日志</a-breadcrumb-item
      >
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <div class="main-layout">
      <div class="left-sidebar">
        <div class="sidebar-header">
          <span class="sidebar-title">Available fields</span>
          <span class="field-count" v-if="data.availableFields.length > 0">
            {{ data.availableFields.length }}
          </span>
        </div>
        <div class="field-search">
          <a-input v-model:value="data.fieldSearchText" placeholder="搜索字段" size="small" allowClear>
            <template #prefix>
              <SearchOutlined />
            </template>
          </a-input>
        </div>
        <div class="field-list">
          <div class="field-section" v-if="data.popularFields.length > 0">
            <div class="section-title">Popular</div>
            <div v-for="field in data.popularFields" :key="field.name" class="field-item"
              @mouseenter="fn.onFieldHover(field, true)" @mouseleave="fn.onFieldHover(field, false)">
              <span class="field-icon" :class="'field-type-' + field.type">
                {{ fn.getFieldIcon(field.type) }}
              </span>
              <span class="field-name" @click="fn.onFieldClick(field)">{{ field.name }}</span>
              <span class="field-type-label">{{ field.typeLabel }}</span>
              <span class="field-actions">
                <a class="field-remove-btn" @click.stop="fn.removeFromPopular(field)"
                  v-show="data.hoveredField === field.name">
                  <MinusOutlined />
                </a>
              </span>
            </div>
          </div>
          <div class="field-section">
            <div class="section-title">All fields</div>
            <div v-for="field in filteredFields" :key="field.name" class="field-item"
              @mouseenter="fn.onFieldHover(field, true)" @mouseleave="fn.onFieldHover(field, false)">
              <span class="field-icon" :class="'field-type-' + field.type">
                {{ fn.getFieldIcon(field.type) }}
              </span>
              <span class="field-name" @click="fn.onFieldClick(field)">{{ field.name }}</span>
              <span class="field-type-label">{{ field.typeLabel }}</span>
              <span class="field-actions">
                <a class="field-add-popular-btn" @click.stop="fn.addToPopular(field)" title="添加到 Popular">
                  <PlusOutlined />
                </a>
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="right-content">
        <div class="query-panel">
          <div class="time-range">
            <a-range-picker v-model:value="data.timeRange" show-time format="YYYY-MM-DD HH:mm:ss"
              :placeholder="['开始时间', '结束时间']" @change="fn.onTimeRangeChange" style="width: 400px" />
            <a-button type="primary" @click="fn.quickTimeRange('1h')" style="margin-left: 10px">
              最近1小时
            </a-button>
            <a-button @click="fn.quickTimeRange('24h')" style="margin-left: 10px">
              最近24小时
            </a-button>
            <a-button @click="fn.quickTimeRange('7d')" style="margin-left: 10px">
              最近7天
            </a-button>
          </div>

          <div class="query-bar">
            <a-input-search v-model:value="data.queryString" placeholder="输入查询条件" enter-button="查询"
              @search="fn.onQuery" style="width: 100%" />
          </div>
        </div>

        <div class="kibana-style-panel" v-if="data.aggregation">
          <div class="panel-header">
            <div class="header-left">
              <span class="total-hits">
                <strong>{{ fn.formatNumber(data.aggregation.total_requests) }}</strong> hits
              </span>
              <span class="time-range-display" v-if="data.timeRange && data.timeRange.length === 2">
                {{ fn.formatTimeRange(data.timeRange[0], data.timeRange[1]) }}
              </span>
              <span class="time-interval-selector" v-if="data.showChart && data.aggregation.time_series && data.aggregation.time_series.length > 0">
                <span style="margin-left: 16px; margin-right: 8px; color: #666;">@timestamp per</span>
                <a-select v-model:value="data.timeInterval" style="width: 120px" size="small">
                  <a-select-option value="auto">Auto</a-select-option>
                  <a-select-option value="second">Second</a-select-option>
                  <a-select-option value="minute">Minute</a-select-option>
                  <a-select-option value="hour">Hour</a-select-option>
                  <a-select-option value="day">Day</a-select-option>
                </a-select>
              </span>
            </div>
            <div class="header-right">
              <a-button type="link" size="small" @click="data.showChart = !data.showChart">
                {{ data.showChart ? '隐藏图表' : '显示图表' }}
              </a-button>
            </div>
          </div>

          <div class="time-series-chart-container"
            v-if="data.showChart && data.aggregation.time_series && data.aggregation.time_series.length > 0">
            <div class="kibana-chart-wrapper">
              <div class="y-axis">
                <div v-for="(tick, index) in fn.getYAxisTicks()" :key="index" class="y-tick"
                  :style="{ bottom: (tick.value / fn.getMaxTimeSeriesCount() * 100) + '%' }">
                  {{ fn.formatYAxisValue(tick.value) }}
                </div>
              </div>
              <div class="chart-area">
                <div class="time-series-bars">
                  <div v-for="(item, index) in data.aggregation.time_series" :key="index" class="kibana-bar" :style="{
                    height: (item.count / fn.getMaxTimeSeriesCount() * 100) + '%',
                    width: (100 / data.aggregation.time_series.length) + '%',
                    backgroundColor: item.count > 0 ? '#00b3a4' : '#e0e0e0'
                  }" :title="fn.formatTime(item.time) + ': ' + fn.formatNumber(item.count)"></div>
                </div>
                <div class="x-axis">
                  <div v-for="(item, index) in fn.getXAxisLabels()" :key="index" class="x-label"
                    :style="{ left: (item.position * 100) + '%' }">
                    {{ item.label }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <a-modal v-model:visible="data.aggregationModalVisible" :title="data.aggregationModalTitle" width="600px"
          :footer="null">
          <div class="aggregation-result" v-if="data.aggregationResult">
            <div class="aggregation-info" v-if="data.aggregationResult.exists_count !== undefined">
              <span>存在于 {{ data.aggregationResult.exists_count }} / {{ data.aggregationResult.total }} 条记录</span>
            </div>
            <div v-if="data.aggregationResult.type === 'count' || data.aggregationResult.type === 'terms'">
              <div class="aggregation-title">Top {{ data.aggregationResult.results.length }} values</div>
              <div class="terms-list">
                <div v-for="(item, index) in data.aggregationResult.results" :key="index" class="term-item">
                  <div class="term-header">
                    <span class="term-value">{{ item.value }}</span>
                    <span class="term-percentage" v-if="item.percentage">
                      {{ item.percentage }}%
                    </span>
                    <span class="term-count">{{ fn.formatNumber(item.count) }}</span>
                  </div>
                  <div class="term-bar-wrapper">
                    <div class="term-bar" :style="{
                      width: (item.count / data.aggregationResult.exists_count * 100) + '%',
                      backgroundColor: '#00b3a4'
                    }"></div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else-if="data.aggregationResult.type === 'stats'">
              <a-descriptions :column="2" bordered>
                <a-descriptions-item label="计数">{{ data.aggregationResult.results.count }}</a-descriptions-item>
                <a-descriptions-item label="最小值">{{ data.aggregationResult.results.min }}</a-descriptions-item>
                <a-descriptions-item label="最大值">{{ data.aggregationResult.results.max }}</a-descriptions-item>
                <a-descriptions-item label="平均值">{{ data.aggregationResult.results.avg?.toFixed(3)
                  }}</a-descriptions-item>
                <a-descriptions-item label="总和" :span="2">{{ fn.formatNumber(data.aggregationResult.results.sum)
                  }}</a-descriptions-item>
              </a-descriptions>
            </div>
            <div v-else>
              <a-descriptions :column="1" bordered>
                <a-descriptions-item v-for="(value, key) in data.aggregationResult.results" :key="key"
                  :label="key.toUpperCase()">
                  {{ typeof value === 'number' ? fn.formatNumber(value) : value }}
                </a-descriptions-item>
              </a-descriptions>
            </div>
          </div>
        </a-modal>

        <div class="log-list" ref="logListRef">
          <a-table :columns="data.columns" :data-source="data.list" :pagination="false" :loading="data.loading"
            :scroll="{ x: 'max-content', y: 'calc(100vh - 64px - 20px - 200px)' }" size="small" :expanded-row-keys="data.expandedRowKeys" @expand="fn.onExpand">
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'response_status'">
                <a-tag :color="fn.getStatusColor(record.response_status)">
                  {{ record.response_status }}
                </a-tag>
              </template>

              <template v-if="column.dataIndex === 'request_method'">
                <a-tag :color="fn.getMethodColor(record.request_method)">
                  {{ record.request_method }}
                </a-tag>
              </template>

              <template v-if="column.dataIndex === 'request_time'">
                <span :class="fn.getRequestTimeClass(record.request_time)">
                  {{ record.request_time }}s
                </span>
              </template>

              <template v-if="column.dataIndex === 'bytes_sent'">
                {{ fn.formatBytes(record.bytes_sent) }}
              </template>

              <template v-if="column.dataIndex === 'timestamp'">
                {{ fn.formatTime(record.timestamp) }}
              </template>
            </template>

            <template #expandedRowRender="{ record }">
              <div class="log-detail">
                <a-tabs>
                  <a-tab-pane key="request" tab="请求信息">
                    <div class="detail-section">
                      <div class="detail-item">
                        <strong>请求URI:</strong>
                        <div class="detail-value">{{ record.request_uri }}</div>
                      </div>
                      <div class="detail-item">
                        <strong>请求路径:</strong>
                        <div class="detail-value">{{ record.request_path }}</div>
                      </div>
                      <div class="detail-item" v-if="record.request_query_string">
                        <strong>查询字符串:</strong>
                        <div class="detail-value">{{ record.request_query_string }}</div>
                      </div>
                      <div class="detail-item" v-if="record.request_headers">
                        <strong>请求头:</strong>
                        <pre class="detail-value json-view">{{ fn.formatJson(record.request_headers) }}</pre>
                      </div>
                      <div class="detail-item" v-if="record.request_args">
                        <strong>请求参数:</strong>
                        <pre class="detail-value json-view">{{ fn.formatJson(record.request_args) }}</pre>
                      </div>
                      <div class="detail-item" v-if="record.request_body">
                        <strong>请求体:</strong>
                        <pre class="detail-value json-view">{{ fn.formatJson(record.request_body) }}</pre>
                      </div>
                    </div>
                  </a-tab-pane>

                  <a-tab-pane key="response" tab="响应信息">
                    <div class="detail-section">
                      <div class="detail-item" v-if="record.response_headers">
                        <strong>响应头:</strong>
                        <pre class="detail-value json-view">{{ fn.formatJson(record.response_headers) }}</pre>
                      </div>
                      <div class="detail-item" v-if="record.response_body">
                        <strong>响应体:</strong>
                        <pre class="detail-value json-view">{{ fn.formatJson(record.response_body) }}</pre>
                      </div>
                    </div>
                  </a-tab-pane>

                  <a-tab-pane key="upstream" tab="上游信息">
                    <div class="detail-section">
                      <div class="detail-item">
                        <strong>上游响应时间:</strong>
                        <div class="detail-value">{{ record.upstream_response_time || '-' }}</div>
                      </div>
                      <div class="detail-item">
                        <strong>上游连接时间:</strong>
                        <div class="detail-value">{{ record.upstream_connect_time || '-' }}</div>
                      </div>
                      <div class="detail-item">
                        <strong>服务器地址:</strong>
                        <div class="detail-value">{{ record.server_addr || '-' }}</div>
                      </div>
                      <div class="detail-item">
                        <strong>服务器端口:</strong>
                        <div class="detail-value">{{ record.server_port || '-' }}</div>
                      </div>
                    </div>
                  </a-tab-pane>
                </a-tabs>
              </div>
            </template>
          </a-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, onMounted, computed, ref, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, PlusOutlined, MinusOutlined } from '@ant-design/icons-vue'
import { $accessLogList, $accessLogAggregation, $fieldAggregation } from '@/api/log'
import dayjs from 'dayjs'
import 'dayjs/plugin/utc'

export default {
  components: {
    SearchOutlined,
    PlusOutlined,
    MinusOutlined
  },
  setup() {
    const logListRef = ref(null)
    const isLoadingMore = ref(false)
    const data = reactive({
      timeRange: null,
      queryString: '',
      params: {
        page: 1,
        page_size: 500,
        start_time: 0,
        end_time: 0,
        search: ''
      },
      list: [],
      total: 0,
      loading: false,
      expandedRowKeys: [],
      aggregation: null,
      showChart: true,
      timeInterval: 'auto',
      availableFields: [],
      popularFields: [],
      fieldSearchText: '',
      aggregationModalVisible: false,
      aggregationModalTitle: '',
      aggregationResult: null,
      currentField: null,
      hoveredField: null,
      columns: [
        {
          title: '时间',
          dataIndex: 'timestamp',
          width: 180,
          fixed: 'left'
        },
        {
          title: '方法',
          dataIndex: 'request_method',
          width: 80
        },
        {
          title: '状态',
          dataIndex: 'response_status',
          width: 80
        },
        {
          title: '请求路径',
          dataIndex: 'request_path',
          width: 300,
          ellipsis: true
        },
        {
          title: 'Host',
          dataIndex: 'request_host',
          width: 200,
          ellipsis: true
        },
        {
          title: '客户端IP',
          dataIndex: 'remote_addr',
          width: 140
        },
        {
          title: '服务名称',
          dataIndex: 'service_name',
          width: 150,
          ellipsis: true
        },
        {
          title: '路由名称',
          dataIndex: 'router_name',
          width: 150,
          ellipsis: true
        },
        {
          title: '响应时间',
          dataIndex: 'request_time',
          width: 100
        },
        {
          title: '响应大小',
          dataIndex: 'bytes_sent',
          width: 100
        }
      ]
    })

    const fn = {
      onTimeRangeChange: () => {
        if (data.timeRange && data.timeRange.length === 2) {
          data.params.start_time = Math.floor(data.timeRange[0].valueOf() / 1000)
          data.params.end_time = Math.floor(data.timeRange[1].valueOf() / 1000)
        } else {
          data.params.start_time = 0
          data.params.end_time = 0
        }
        fn.onQuery()
      },

      quickTimeRange: type => {
        const now = dayjs()
        let start
        switch (type) {
          case '1h':
            start = now.subtract(1, 'hour')
            break
          case '24h':
            start = now.subtract(24, 'hour')
            break
          case '7d':
            start = now.subtract(7, 'day')
            break
          default:
            start = now.subtract(1, 'hour')
        }
        data.timeRange = [start, now]
        fn.onTimeRangeChange()
      },

      onQuery: async (isLoadMore = false) => {
        if (isLoadMore) {
          isLoadingMore.value = true
        } else {
          data.loading = true
        }
        const queryParams = { ...data.params }
        if (!queryParams.search) {
          delete queryParams.search
        }
        if (queryParams.start_time === 0) {
          delete queryParams.start_time
        }
        if (queryParams.end_time === 0) {
          delete queryParams.end_time
        }

        try {
          const [listRes, aggRes] = await Promise.all([
            $accessLogList(queryParams),
            isLoadMore ? Promise.resolve({ code: 0, data: data.aggregation }) : $accessLogAggregation(queryParams)
          ])

          if (listRes.code === 0) {
            if (isLoadMore) {
              data.list = [...data.list, ...(listRes.data.data || [])]
            } else {
              data.list = listRes.data.data || []
            }
            data.total = listRes.data.total || 0
            fn.extractFields()
            nextTick(() => {
              if (logListRef.value) {
                const tableBody = logListRef.value.querySelector('.ant-table-body')
                if (tableBody) {
                  if (data.params.page === 1 && !isLoadMore) {
                    tableBody.scrollTop = 0
                  }
                  if (!isLoadMore) {
                    tableBody.removeEventListener('scroll', fn.onTableScroll)
                    tableBody.addEventListener('scroll', fn.onTableScroll)
                  }
                }
              }
            })
          } else {
            message.error(listRes.msg || '获取访问日志失败')
          }

          if (aggRes.code === 0 && !isLoadMore) {
            data.aggregation = aggRes.data || null
          }
        } catch (error) {
          message.error('查询失败: ' + error.message)
        } finally {
          data.loading = false
          isLoadingMore.value = false
        }
      },

      onTableScroll: (event) => {
        const tableBody = event.target
        const scrollTop = tableBody.scrollTop
        const scrollHeight = tableBody.scrollHeight
        const clientHeight = tableBody.clientHeight
        
        if (scrollTop + clientHeight >= scrollHeight - 10) {
          const currentPage = data.params.page
          const totalPages = Math.ceil(data.total / data.params.page_size)
          
          if (currentPage < totalPages && !isLoadingMore.value && !data.loading) {
            data.params.page = currentPage + 1
            fn.onQuery(true)
          }
        }
      },

      onExpand: (expanded, record) => {
        if (expanded) {
          data.expandedRowKeys = [record.id]
        } else {
          data.expandedRowKeys = []
        }
      },

      formatTime: timestamp => {
        if (!timestamp) return '-'
        return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')
      },

      formatBytes: bytes => {
        if (!bytes) return '0 B'
        const k = 1024
        const sizes = ['B', 'KB', 'MB', 'GB']
        const i = Math.floor(Math.log(bytes) / Math.log(k))
        return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
      },

      formatJson: str => {
        if (!str) return ''
        try {
          const obj = JSON.parse(str)
          return JSON.stringify(obj, null, 2)
        } catch {
          return str
        }
      },

      getStatusColor: status => {
        if (status >= 500) return 'red'
        if (status >= 400) return 'orange'
        if (status >= 300) return 'blue'
        return 'green'
      },

      getMethodColor: method => {
        const colors = {
          GET: 'blue',
          POST: 'green',
          PUT: 'orange',
          DELETE: 'red',
          PATCH: 'purple'
        }
        return colors[method] || 'default'
      },

      getRequestTimeClass: time => {
        if (time > 1) return 'request-time-slow'
        if (time > 0.5) return 'request-time-medium'
        return 'request-time-fast'
      },

      getMaxTimeSeriesCount: () => {
        if (!data.aggregation || !data.aggregation.time_series) return 1
        return Math.max(...data.aggregation.time_series.map(item => item.count), 1)
      },

      formatNumber: num => {
        if (!num && num !== 0) return '0'
        return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
      },

      formatTimeRange: (start, end) => {
        if (!start || !end) return ''
        const startStr = dayjs(start).format('MMM DD, YYYY@HH:mm:ss.SSS')
        const endStr = dayjs(end).format('MMM DD, YYYY@HH:mm:ss.SSS')
        return `${startStr} - ${endStr}`
      },

      getTimeInterval: () => {
        if (data.timeInterval !== 'auto') {
          return data.timeInterval
        }
        if (!data.aggregation || !data.aggregation.time_series || data.aggregation.time_series.length < 2) return '30 seconds'
        const interval = data.aggregation.time_series[1].time - data.aggregation.time_series[0].time
        if (interval < 60) return `${interval} seconds`
        if (interval < 3600) return `${Math.floor(interval / 60)} minutes`
        return `${Math.floor(interval / 3600)} hours`
      },

      getTimeIntervalDisplay: () => {
        if (data.timeInterval === 'auto') {
          return fn.getTimeInterval()
        }
        const labels = {
          second: 'Second',
          minute: 'Minute',
          hour: 'Hour',
          day: 'Day'
        }
        return labels[data.timeInterval] || 'Auto'
      },

      getYAxisTicks: () => {
        const max = fn.getMaxTimeSeriesCount()
        const ticks = []
        const tickCount = 3
        for (let i = 0; i <= tickCount; i++) {
          const value = Math.floor((max / tickCount) * i)
          ticks.push({ value, position: i / tickCount })
        }
        return ticks
      },

      formatYAxisValue: value => {
        if (value >= 1000000) return (value / 1000000).toFixed(1) + 'M'
        if (value >= 1000) return (value / 1000).toFixed(1) + 'K'
        return value.toString()
      },

      getXAxisLabels: () => {
        if (!data.aggregation || !data.aggregation.time_series) return []
        const labels = []
        const series = data.aggregation.time_series
        if (series.length === 0) return []
        
        const labelCount = Math.min(5, series.length)
        if (labelCount === 1) {
          return []
        } else {
          const step = Math.floor(series.length / (labelCount - 1))
          for (let i = step; i < series.length; i += step) {
            if (i >= series.length) break
            labels.push({
              label: dayjs.unix(series[i].time).format('HH:mm:ss'),
              position: i / (series.length - 1)
            })
          }
          if (series.length > 0 && labels[labels.length - 1].position < 0.99) {
            labels.push({
              label: dayjs.unix(series[series.length - 1].time).format('HH:mm:ss'),
              position: 1
            })
          }
        }
        return labels
      },

      extractFields: () => {
        if (!data.list || data.list.length === 0) {
          data.availableFields = []
          data.popularFields = []
          return
        }

        const fieldMap = new Map()
        const popularFieldNames = [
          'timestamp',
          'request_method',
          'response_status',
          'request_path',
          'request_host',
          'remote_addr',
          'service_name',
          'router_name',
          'request_time',
          'bytes_sent'
        ]

        data.list.forEach(item => {
          Object.keys(item).forEach(key => {
            if (!fieldMap.has(key)) {
              const value = item[key]
              let type = 'text'
              let typeLabel = 'text'

              if (key === 'timestamp' || key.includes('time') || key.includes('Time')) {
                type = 'timestamp'
                typeLabel = 'date'
              } else if (typeof value === 'number' || key.includes('count') || key.includes('status') || key.includes('port')) {
                type = 'number'
                typeLabel = 'number'
              } else if (typeof value === 'boolean') {
                type = 'boolean'
                typeLabel = 'boolean'
              }

              fieldMap.set(key, {
                name: key,
                type,
                typeLabel,
                isPopular: popularFieldNames.includes(key)
              })
            }
          })
        })

        const allFields = Array.from(fieldMap.values()).sort((a, b) => {
          if (a.isPopular && !b.isPopular) return -1
          if (!a.isPopular && b.isPopular) return 1
          return a.name.localeCompare(b.name)
        })

        data.availableFields = allFields
        data.popularFields = allFields.filter(f => f.isPopular)
      },

      getFieldIcon: type => {
        const icons = {
          text: 't',
          number: '#',
          timestamp: '📅',
          boolean: '✓'
        }
        return icons[type] || 't'
      },

      onFieldClick: field => {
        if (field.type === 'number') {
          fn.onAggregationSelect(field, { key: 'stats' })
        } else {
          fn.onAggregationSelect(field, { key: 'terms' })
        }
      },

      onFieldHover: (field, isHover) => {
        if (isHover) {
          data.hoveredField = field.name
        } else {
          data.hoveredField = null
        }
      },

      isInPopular: field => {
        return data.popularFields.some(f => f.name === field.name)
      },

      addToPopular: field => {
        if (fn.isInPopular(field)) {
          return
        }
        data.popularFields.push(field)
        field.isPopular = true
        message.success(`已将 ${field.name} 添加到 Popular`)
      },

      removeFromPopular: field => {
        const index = data.popularFields.findIndex(f => f.name === field.name)
        if (index > -1) {
          data.popularFields.splice(index, 1)
          field.isPopular = false
          message.success(`已将 ${field.name} 从 Popular 移除`)
        }
      },

      onFieldMenuVisible: (field, visible) => {
        if (visible) {
          data.currentField = field
        }
      },

      onAggregationSelect: (field, event) => {
        const aggregationType = event.key
        const typeNames = {
          count: '计数',
          terms: '分组统计',
          avg: '平均值',
          max: '最大值',
          min: '最小值',
          sum: '求和',
          stats: '统计'
        }

        data.aggregationModalTitle = `${field.name} - ${typeNames[aggregationType] || aggregationType}`
        data.aggregationModalVisible = true
        data.aggregationResult = null

        const result = fn.calculateFieldAggregation(field, aggregationType)
        data.aggregationResult = result
      },

      calculateFieldAggregation: (field, aggregationType) => {
        if (!data.list || data.list.length === 0) {
          message.warning('没有可聚合的数据')
          return null
        }

        const fieldName = field.name
        const fieldValues = data.list.map(item => item[fieldName]).filter(v => v !== null && v !== undefined && v !== '')

        if (fieldValues.length === 0) {
          message.warning(`字段 ${fieldName} 在当前数据中没有值`)
          return null
        }

        const total = data.list.length
        const existsCount = fieldValues.length

        let results = null

        switch (aggregationType) {
          case 'count':
            results = [
              { value: '总计', count: total },
              { value: '有值', count: existsCount },
              { value: '空值', count: total - existsCount }
            ]
            break

          case 'terms':
            const valueCountMap = new Map()
            fieldValues.forEach(value => {
              const key = String(value)
              valueCountMap.set(key, (valueCountMap.get(key) || 0) + 1)
            })

            const terms = Array.from(valueCountMap.entries())
              .map(([value, count]) => ({
                value: value.length > 50 ? value.substring(0, 50) + '...' : value,
                count,
                percentage: ((count / existsCount) * 100).toFixed(1)
              }))
              .sort((a, b) => b.count - a.count)
              .slice(0, 20)

            results = terms
            break

          case 'avg':
            if (field.type !== 'number') {
              message.warning('平均值只能用于数字类型字段')
              return null
            }
            const numericValues = fieldValues.filter(v => !isNaN(Number(v))).map(v => Number(v))
            if (numericValues.length === 0) {
              message.warning('字段中没有有效的数字值')
              return null
            }
            const avg = numericValues.reduce((sum, val) => sum + val, 0) / numericValues.length
            results = { avg: avg.toFixed(3) }
            break

          case 'max':
            if (field.type !== 'number') {
              message.warning('最大值只能用于数字类型字段')
              return null
            }
            const maxValues = fieldValues.filter(v => !isNaN(Number(v))).map(v => Number(v))
            if (maxValues.length === 0) {
              message.warning('字段中没有有效的数字值')
              return null
            }
            results = { max: Math.max(...maxValues) }
            break

          case 'min':
            if (field.type !== 'number') {
              message.warning('最小值只能用于数字类型字段')
              return null
            }
            const minValues = fieldValues.filter(v => !isNaN(Number(v))).map(v => Number(v))
            if (minValues.length === 0) {
              message.warning('字段中没有有效的数字值')
              return null
            }
            results = { min: Math.min(...minValues) }
            break

          case 'sum':
            if (field.type !== 'number') {
              message.warning('求和只能用于数字类型字段')
              return null
            }
            const sumValues = fieldValues.filter(v => !isNaN(Number(v))).map(v => Number(v))
            if (sumValues.length === 0) {
              message.warning('字段中没有有效的数字值')
              return null
            }
            const sum = sumValues.reduce((s, val) => s + val, 0)
            results = { sum }
            break

          case 'stats':
            if (field.type !== 'number') {
              message.warning('统计只能用于数字类型字段')
              return null
            }
            const statsValues = fieldValues.filter(v => !isNaN(Number(v))).map(v => Number(v))
            if (statsValues.length === 0) {
              message.warning('字段中没有有效的数字值')
              return null
            }
            const sorted = [...statsValues].sort((a, b) => a - b)
            const statsAvg = statsValues.reduce((s, val) => s + val, 0) / statsValues.length
            const statsSum = statsValues.reduce((s, val) => s + val, 0)
            results = {
              count: statsValues.length,
              min: sorted[0],
              max: sorted[sorted.length - 1],
              avg: statsAvg,
              sum: statsSum
            }
            break

          default:
            results = []
        }

        return {
          field_name: fieldName,
          type: aggregationType,
          results,
          total: total,
          exists_count: existsCount
        }
      }
    }

    const filteredFields = computed(() => {
      if (!data.fieldSearchText) {
        return data.availableFields.filter(f => !f.isPopular)
      }
      const searchLower = data.fieldSearchText.toLowerCase()
      return data.availableFields.filter(
        f => !f.isPopular && f.name.toLowerCase().includes(searchLower)
      )
    })

    onMounted(() => {
      fn.quickTimeRange('1h')
    })

    return {
      data,
      fn,
      filteredFields,
      logListRef
    }
  }
}
</script>

<style scoped>
.main {
  padding: 10px;
}

.breadcrumb {
  margin-bottom: 10px;
}

.main-layout {
  display: flex;
  gap: 16px;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.left-sidebar {
  width: 280px;
  background: #fff;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  border: 1px solid #e8e8e8;
  overflow: hidden;
}

.sidebar-header {
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.field-count {
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 12px;
}

.field-search {
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
}

.field-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.field-section {
  margin-bottom: 16px;
}

.section-title {
  padding: 8px 16px;
  font-size: 12px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.field-item {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.field-item:hover {
  background-color: #f5f5f5;
}

.field-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  border-radius: 2px;
  flex-shrink: 0;
}

.field-type-text {
  background: #e3f2fd;
  color: #1976d2;
}

.field-type-number {
  background: #e8f5e9;
  color: #388e3c;
}

.field-type-timestamp {
  background: #fff3e0;
  color: #f57c00;
}

.field-type-boolean {
  background: #f3e5f5;
  color: #7b1fa2;
}

.field-name {
  flex: 1;
  font-size: 13px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.field-type-label {
  font-size: 11px;
  color: #999;
  text-transform: lowercase;
}

.field-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.field-item:hover .field-actions {
  opacity: 1;
}

.field-item:hover .field-add-popular-btn {
  opacity: 1;
}

.field-add-btn {
  color: #1890ff;
  font-size: 14px;
  padding: 2px 4px;
  cursor: pointer;
}

.field-add-btn:hover {
  color: #40a9ff;
}

.field-add-popular-btn {
  color: #1890ff;
  font-size: 14px;
  padding: 2px 4px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.field-add-popular-btn:hover {
  color: #40a9ff;
  background-color: #e6f7ff;
  border-radius: 4px;
}

.field-remove-btn {
  color: #ff4d4f;
  font-size: 14px;
  padding: 2px 4px;
  cursor: pointer;
}

.field-remove-btn:hover {
  color: #ff7875;
}

.aggregation-result {
  max-height: 600px;
  overflow-y: auto;
}

.aggregation-info {
  margin-bottom: 16px;
  padding: 8px 12px;
  background: #f5f5f5;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
}

.aggregation-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.terms-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.term-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.term-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.term-value {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #333;
}

.term-percentage {
  color: #666;
  font-size: 12px;
  min-width: 50px;
  text-align: right;
}

.term-count {
  color: #666;
  font-size: 12px;
  min-width: 60px;
  text-align: right;
}

.term-bar-wrapper {
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.term-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.right-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  background: #fff;
}

.query-panel {
  background: #fff;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 0px;
  flex-shrink: 0;
}

.time-range {
  margin-bottom: 16px;
}

.query-bar {
  margin-bottom: 0px;
}

.log-list {
  background: #fff;
  padding: 16px;
  border-radius: 4px;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.log-detail {
  padding: 16px;
  background: #f5f5f5;
}

.detail-section {
  max-height: 600px;
  overflow-y: auto;
}

.detail-item {
  margin-bottom: 16px;
}

.detail-item strong {
  display: block;
  margin-bottom: 8px;
  color: #333;
}

.detail-value {
  background: #fff;
  padding: 8px;
  border-radius: 4px;
  word-break: break-all;
  white-space: pre-wrap;
}

.json-view {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
}

.request-time-fast {
  color: #52c41a;
}

.request-time-medium {
  color: #faad14;
}

.request-time-slow {
  color: #ff4d4f;
}

.aggregation-panel {
  background: #fff;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
}

.chart-container {
  background: #fafafa;
  padding: 16px;
  border-radius: 4px;
}

.chart-container h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 500;
}

.chart-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bar-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.bar-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.bar-label .service-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 8px;
}

.bar {
  height: 20px;
  background: #e8e8e8;
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: 2px;
}

.time-series-chart {
  display: flex;
  align-items: flex-end;
  height: 200px;
  gap: 2px;
  margin-bottom: 8px;
  padding: 8px 0;
}

.time-bar {
  background: #1890ff;
  border-radius: 2px;
  cursor: pointer;
  transition: opacity 0.2s;
  min-height: 2px;
}

.time-bar:hover {
  opacity: 0.8;
}

.time-series-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #666;
  margin-top: 4px;
}

.time-label {
  flex: 1;
  text-align: center;
}

.kibana-style-panel {
  background: #fff;
  border-radius: 4px;
  margin-bottom: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.total-hits {
  font-size: 14px;
  color: #333;
}

.total-hits strong {
  font-size: 16px;
  font-weight: 600;
  color: #000;
}

.time-range-display {
  font-size: 12px;
  color: #666;
  font-family: 'Courier New', monospace;
}

.header-right {
  display: flex;
  align-items: center;
}

.time-series-chart-container {
  padding: 20px;
  background: #fafafa;
  border-radius: 4px;
  margin-top: 8px;
}

.chart-header {
  margin-bottom: 8px;
}

.chart-title {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.kibana-chart-wrapper {
  display: flex;
  position: relative;
  height: 150px;
  background: #fff;
  border-radius: 4px;
  padding: 12px;
  border: 1px solid #e8e8e8;
}

.y-axis {
  width: 50px;
  position: relative;
  padding-right: 12px;
  flex-shrink: 0;
}

.y-tick {
  position: absolute;
  right: 0;
  transform: translateY(50%);
  font-size: 12px;
  color: #8c8c8c;
  white-space: nowrap;
  font-weight: 500;
}

.chart-area {
  flex: 1;
  position: relative;
  padding-left: 12px;
  padding-right: 8px;
}

.time-series-bars {
  display: flex;
  align-items: flex-end;
  height: calc(100% - 30px);
  gap: 2px;
  position: relative;
  margin-bottom: 30px;
}

.time-series-bars::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    repeating-linear-gradient(to bottom, transparent, transparent 29px, #f0f0f0 29px, #f0f0f0 30px);
  pointer-events: none;
  z-index: 0;
}

.kibana-bar {
  background: linear-gradient(to top, #1890ff, #40a9ff);
  border-radius: 2px 2px 0 0;
  cursor: pointer;
  transition: all 0.2s;
  min-height: 2px;
  position: relative;
  z-index: 1;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.kibana-bar:hover {
  background: linear-gradient(to top, #096dd9, #1890ff);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  transform: translateY(-1px);
}

.x-axis {
  position: absolute;
  bottom: 0;
  left: 12px;
  right: 8px;
  height: 35px;
  border-top: 1px solid #d9d9d9;
  padding-top: 4px;
  overflow: visible;
}

.x-label {
  position: absolute;
  transform: translateX(-50%);
  font-size: 10px;
  color: #8c8c8c;
  white-space: nowrap;
  font-weight: 500;
  pointer-events: none;
}
</style>
