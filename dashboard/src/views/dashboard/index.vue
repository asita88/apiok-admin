<template>
  <div class="main gw-dash">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item
        ><i
          style="color: #448ef7; font-size: 30px"
          class="iconfont icon-yuntongbu"
        />大盘</a-breadcrumb-item
      >
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <div class="dashboard-content">
      <div class="time-filter gw-toolbar">
        <div class="gw-toolbar__label">时间范围</div>
        <div class="gw-toolbar__controls">
          <a-range-picker
            v-model:value="data.timeRange"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            :placeholder="['开始时间', '结束时间']"
            @change="fn.onTimeRangeChange"
            class="gw-toolbar__picker"
          />
          <a-button type="primary" class="gw-toolbar__btn" @click="fn.quickTimeRange('1h')">
            最近1小时
          </a-button>
          <a-button class="gw-toolbar__btn" @click="fn.quickTimeRange('24h')">
            最近24小时
          </a-button>
          <a-button class="gw-toolbar__btn" @click="fn.quickTimeRange('7d')">
            最近7天
          </a-button>
          <a-button class="gw-toolbar__btn" @click="fn.refresh" :loading="data.loading">
            刷新
          </a-button>
        </div>
      </div>

      <a-spin :spinning="data.loading">
        <div class="stats-grid" v-if="data.aggregation">
          <div class="stat-card stat-card--indigo">
            <div class="stat-icon">
              <ThunderboltOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">总请求数</div>
              <div class="stat-value stat-value--nums">{{ fn.formatNumber(data.aggregation.total_requests) }}</div>
            </div>
          </div>

          <div class="stat-card stat-card--emerald">
            <div class="stat-icon">
              <FieldTimeOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">平均响应时间</div>
              <div class="stat-value stat-value--nums">{{ fn.formatNumber(data.aggregation.avg_response_time, 2) }}ms</div>
            </div>
          </div>

          <div class="stat-card stat-card--amber">
            <div class="stat-icon">
              <WarningOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">错误数</div>
              <div class="stat-value stat-value--nums">{{ fn.formatNumber(data.aggregation.error_count) }}</div>
            </div>
          </div>

          <div class="stat-card stat-card--rose">
            <div class="stat-icon">
              <PieChartOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">错误率</div>
              <div class="stat-value stat-value--nums">{{ fn.formatNumber(data.aggregation.error_rate * 100, 2) }}%</div>
            </div>
          </div>

          <div class="stat-card stat-card--violet">
            <div class="stat-icon">
              <CloudUploadOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">总流量</div>
              <div class="stat-value stat-value--nums">{{ fn.formatBytes(data.aggregation.total_bytes_sent) }}</div>
            </div>
          </div>

          <div class="stat-card stat-card--cyan">
            <div class="stat-icon">
              <ClockCircleOutlined />
            </div>
            <div class="stat-content">
              <div class="stat-label">最大响应时间</div>
              <div class="stat-value stat-value--nums">{{ fn.formatNumber(data.aggregation.max_response_time, 2) }}ms</div>
            </div>
          </div>
        </div>

        <div class="charts-grid" v-if="data.aggregation">
          <div class="chart-card">
            <div class="chart-header">
              <h3>请求趋势</h3>
              <span class="chart-header__hint">时间序列 · 请求次数</span>
            </div>
            <div class="chart-content">
              <div class="time-series-chart">
                <div class="chart-area">
                  <div class="time-series-bars">
                    <div
                      v-for="(item, index) in data.aggregation.time_series"
                      :key="index"
                      class="chart-bar"
                      :style="{
                        height: (item.count / fn.getMaxTimeSeriesCount() * 100) + '%',
                        width: (100 / data.aggregation.time_series.length) + '%'
                      }"
                      :title="fn.formatTime(item.time) + ': ' + fn.formatNumber(item.count)"
                    ></div>
                  </div>
                  <div class="x-axis">
                    <div
                      v-for="(item, index) in fn.getXAxisLabels()"
                      :key="index"
                      class="x-label"
                      :style="{ left: (item.position * 100) + '%' }"
                    >
                      {{ item.label }}
                    </div>
                  </div>
                </div>
                <div class="y-axis">
                  <div
                    v-for="(tick, index) in fn.getYAxisTicks()"
                    :key="index"
                    class="y-tick"
                    :style="{ bottom: (tick.position * 100) + '%' }"
                  >
                    {{ fn.formatYAxisValue(tick.value) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>流量趋势</h3>
              <span class="chart-header__hint">出站字节 · 时间序列</span>
            </div>
            <div class="chart-content">
              <div class="time-series-chart">
                <div class="chart-area">
                  <div class="time-series-bars">
                    <div
                      v-for="(item, index) in data.aggregation.bytes_time_series"
                      :key="index"
                      class="chart-bar bytes-bar"
                      :style="{
                        height: (item.bytes / fn.getMaxBytesTimeSeries() * 100) + '%',
                        width: (100 / data.aggregation.bytes_time_series.length) + '%'
                      }"
                      :title="fn.formatTime(item.time) + ': ' + fn.formatBytes(item.bytes)"
                    ></div>
                  </div>
                  <div class="x-axis">
                    <div
                      v-for="(item, index) in fn.getBytesXAxisLabels()"
                      :key="index"
                      class="x-label"
                      :style="{ left: (item.position * 100) + '%' }"
                    >
                      {{ item.label }}
                    </div>
                  </div>
                </div>
                <div class="y-axis">
                  <div
                    v-for="(tick, index) in fn.getBytesYAxisTicks()"
                    :key="index"
                    class="y-tick"
                    :style="{ bottom: (tick.position * 100) + '%' }"
                  >
                    {{ fn.formatBytesYAxisValue(tick.value) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>状态码分布</h3>
            </div>
            <div class="chart-content">
              <div class="status-chart">
                <div
                  v-for="(item, index) in data.aggregation.status_stats"
                  :key="index"
                  class="status-bar"
                >
                  <div class="status-label">
                    <span class="status-code" :class="'status-' + item.status">
                      {{ item.status }}
                    </span>
                    <span class="status-count">{{ fn.formatNumber(item.count) }}</span>
                  </div>
                  <div class="status-progress">
                    <div
                      class="status-progress-bar"
                      :style="{
                        width: (item.count / fn.getMaxStatusCount() * 100) + '%',
                        backgroundColor: fn.getStatusColorByRange(item.status)
                      }"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>请求方法分布</h3>
            </div>
            <div class="chart-content">
              <div class="method-chart">
                <div
                  v-for="(item, index) in data.aggregation.method_stats"
                  :key="index"
                  class="method-item"
                >
                  <div class="method-label">{{ item.method }}</div>
                  <div class="method-bar-container">
                    <div
                      class="method-bar"
                      :style="{
                        width: (item.count / fn.getMaxMethodCount() * 100) + '%',
                        backgroundColor: fn.getMethodColor(item.method)
                      }"
                    >
                      <span class="method-count">{{ fn.formatNumber(item.count) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>Top 5 域名</h3>
            </div>
            <div class="chart-content">
              <div class="host-chart">
                <div
                  v-for="(item, index) in data.aggregation.host_stats"
                  :key="index"
                  class="host-item"
                >
                  <div class="host-rank">{{ index + 1 }}</div>
                  <div class="host-info">
                    <div class="host-label" :title="item.host">{{ item.host }}</div>
                    <div class="host-bar-container">
                      <div
                        class="host-bar"
                        :style="{
                          width: (item.count / fn.getMaxHostCount() * 100) + '%'
                        }"
                      >
                        <span class="host-count">{{ fn.formatNumber(item.count) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>Top 5 接口</h3>
            </div>
            <div class="chart-content">
              <div class="path-chart">
                <div
                  v-for="(item, index) in data.aggregation.path_stats"
                  :key="index"
                  class="path-item"
                >
                  <div class="path-rank">{{ index + 1 }}</div>
                  <div class="path-info">
                    <div class="path-label" :title="item.path">{{ item.path }}</div>
                    <div class="path-bar-container">
                      <div
                        class="path-bar"
                        :style="{
                          width: (item.count / fn.getMaxPathCount() * 100) + '%'
                        }"
                      >
                        <span class="path-count">{{ fn.formatNumber(item.count) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="chart-card">
            <div class="chart-header">
              <h3>Top 5 接口流量</h3>
            </div>
            <div class="chart-content">
              <div class="path-bytes-chart">
                <div
                  v-for="(item, index) in data.aggregation.path_bytes_stats"
                  :key="index"
                  class="path-bytes-item"
                >
                  <div class="path-bytes-rank">{{ index + 1 }}</div>
                  <div class="path-bytes-info">
                    <div class="path-bytes-label" :title="item.path">{{ item.path }}</div>
                    <div class="path-bytes-bar-container">
                      <div
                        class="path-bytes-bar"
                        :style="{
                          width: (item.bytes / fn.getMaxPathBytesCount() * 100) + '%'
                        }"
                      >
                        <span class="path-bytes-count">{{ fn.formatBytes(item.bytes) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-spin>
    </div>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import {
  ThunderboltOutlined,
  FieldTimeOutlined,
  WarningOutlined,
  PieChartOutlined,
  CloudUploadOutlined,
  ClockCircleOutlined
} from '@ant-design/icons-vue'
import { $accessLogAggregation } from '@/api/log'

const data = reactive({
  loading: false,
  timeRange: null,
  aggregation: null
})

const fn = {
  onTimeRangeChange: () => {
    fn.loadData()
  },

  quickTimeRange: (range) => {
    const now = dayjs()
    let start
    switch (range) {
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
    fn.loadData()
  },

  refresh: () => {
    fn.loadData()
  },

  loadData: async () => {
    if (!data.timeRange || data.timeRange.length !== 2) {
      message.warning('请选择时间范围')
      return
    }

    data.loading = true
    try {
      const params = {
        start_time: data.timeRange[0].unix(),
        end_time: data.timeRange[1].unix()
      }
      const res = await $accessLogAggregation(params)
      if (res.code === 0) {
        data.aggregation = res.data
      } else {
        message.error(res.msg || '获取数据失败')
      }
    } catch (error) {
      message.error('获取数据失败: ' + error.message)
    } finally {
      data.loading = false
    }
  },

  formatNumber: (num, decimals = 0) => {
    if (num === null || num === undefined) return '0'
    if (num >= 1000000) return (num / 1000000).toFixed(decimals) + 'M'
    if (num >= 1000) return (num / 1000).toFixed(decimals) + 'K'
    return num.toFixed(decimals)
  },

  formatBytes: (bytes) => {
    if (!bytes) return '0 B'
    if (bytes >= 1024 * 1024 * 1024) return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
    if (bytes >= 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
    if (bytes >= 1024) return (bytes / 1024).toFixed(2) + ' KB'
    return bytes + ' B'
  },

  formatTime: (timestamp) => {
    return dayjs.unix(timestamp).format('HH:mm:ss')
  },

  getMaxTimeSeriesCount: () => {
    if (!data.aggregation || !data.aggregation.time_series) return 1
    return Math.max(...data.aggregation.time_series.map(item => item.count), 1)
  },

  getMaxStatusCount: () => {
    if (!data.aggregation || !data.aggregation.status_stats) return 1
    return Math.max(...data.aggregation.status_stats.map(item => item.count), 1)
  },

  getMaxMethodCount: () => {
    if (!data.aggregation || !data.aggregation.method_stats) return 1
    return Math.max(...data.aggregation.method_stats.map(item => item.count), 1)
  },

  getMaxServiceCount: () => {
    if (!data.aggregation || !data.aggregation.service_stats) return 1
    return Math.max(...data.aggregation.service_stats.map(item => item.count), 1)
  },

  getMaxHostCount: () => {
    if (!data.aggregation || !data.aggregation.host_stats) return 1
    return Math.max(...data.aggregation.host_stats.map(item => item.count), 1)
  },

  getMaxPathCount: () => {
    if (!data.aggregation || !data.aggregation.path_stats) return 1
    return Math.max(...data.aggregation.path_stats.map(item => item.count), 1)
  },

  getMaxPathBytesCount: () => {
    if (!data.aggregation || !data.aggregation.path_bytes_stats) return 1
    return Math.max(...data.aggregation.path_bytes_stats.map(item => item.bytes), 1)
  },

  getMaxBytesTimeSeries: () => {
    if (!data.aggregation || !data.aggregation.bytes_time_series) return 1
    return Math.max(...data.aggregation.bytes_time_series.map(item => item.bytes), 1)
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

  formatYAxisValue: (value) => {
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

  getBytesXAxisLabels: () => {
    if (!data.aggregation || !data.aggregation.bytes_time_series) return []
    const labels = []
    const series = data.aggregation.bytes_time_series
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

  getBytesYAxisTicks: () => {
    const max = fn.getMaxBytesTimeSeries()
    const ticks = []
    const tickCount = 3
    for (let i = 0; i <= tickCount; i++) {
      const value = Math.floor((max / tickCount) * i)
      ticks.push({ value, position: i / tickCount })
    }
    return ticks
  },

  formatBytesYAxisValue: (value) => {
    if (value >= 1024 * 1024 * 1024) return (value / (1024 * 1024 * 1024)).toFixed(1) + 'GB'
    if (value >= 1024 * 1024) return (value / (1024 * 1024)).toFixed(1) + 'MB'
    if (value >= 1024) return (value / 1024).toFixed(1) + 'KB'
    return value + 'B'
  },

  getStatusColor: (status) => {
    const code = Math.floor(status / 100)
    switch (code) {
      case 2:
        return '#059669'
      case 3:
        return '#4f46e5'
      case 4:
        return '#d97706'
      case 5:
        return '#e11d48'
      default:
        return '#cbd5e1'
    }
  },

  getStatusColorByRange: (statusRange) => {
    switch (statusRange) {
      case '2xx':
        return '#059669'
      case '3xx':
        return '#4f46e5'
      case '4xx':
        return '#d97706'
      case '5xx':
        return '#e11d48'
      default:
        return '#cbd5e1'
    }
  },

  getMethodColor: (method) => {
    const colors = {
      GET: '#4f46e5',
      POST: '#059669',
      PUT: '#d97706',
      DELETE: '#e11d48',
      PATCH: '#7c3aed'
    }
    return colors[method] || '#94a3b8'
  }
}

onMounted(() => {
  fn.quickTimeRange('24h')
})
</script>

<style scoped>
.breadcrumb {
  font-size: 20px;
}

.gw-dash {
  --gw-slate-50: #f8fafc;
  --gw-slate-200: #e2e8f0;
  --gw-slate-500: #64748b;
  --gw-slate-700: #334155;
  --gw-slate-900: #0f172a;
  --gw-indigo: #4f46e5;
  --gw-indigo-light: #6366f1;
  --gw-card-shadow: 0 1px 2px rgba(15, 23, 42, 0.06), 0 4px 12px rgba(15, 23, 42, 0.06);
  --gw-card-shadow-hover: 0 2px 4px rgba(15, 23, 42, 0.08), 0 8px 20px rgba(15, 23, 42, 0.08);
}

.main {
  padding: 10px;
  height: calc(100vh - 64px - 20px);
  max-height: calc(100vh - 64px - 20px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
}

.dashboard-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.time-filter {
  margin-bottom: 0;
}

.gw-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px 20px;
  margin-bottom: 16px;
  padding: 14px 18px;
  background: #fff;
  border-radius: 10px;
  border: 1px solid var(--gw-slate-200);
  box-shadow: var(--gw-card-shadow);
}

.gw-toolbar__label {
  font-size: 12px;
  font-weight: 600;
  color: var(--gw-slate-500);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.gw-toolbar__controls {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.gw-toolbar__picker {
  width: 400px;
  max-width: 100%;
}

.gw-toolbar__btn {
  transition: color 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
  margin-bottom: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 18px 18px;
  background: #fff;
  border-radius: 10px;
  border: 1px solid var(--gw-slate-200);
  box-shadow: var(--gw-card-shadow);
  transition: box-shadow 0.2s ease, border-color 0.2s ease;
  cursor: default;
}

.stat-card:hover {
  box-shadow: var(--gw-card-shadow-hover);
  border-color: #cbd5e1;
}

.stat-card--indigo .stat-icon {
  background: linear-gradient(145deg, #4f46e5, #6366f1);
}
.stat-card--emerald .stat-icon {
  background: linear-gradient(145deg, #047857, #059669);
}
.stat-card--amber .stat-icon {
  background: linear-gradient(145deg, #b45309, #d97706);
}
.stat-card--rose .stat-icon {
  background: linear-gradient(145deg, #be123c, #e11d48);
}
.stat-card--violet .stat-icon {
  background: linear-gradient(145deg, #6d28d9, #7c3aed);
}
.stat-card--cyan .stat-icon {
  background: linear-gradient(145deg, #0e7490, #0891b2);
}

.stat-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 14px;
  color: #fff;
  font-size: 22px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: var(--gw-slate-500);
  margin-bottom: 6px;
  font-weight: 500;
}

.stat-value {
  font-size: 22px;
  font-weight: 600;
  color: var(--gw-slate-900);
  letter-spacing: -0.02em;
}

.stat-value--nums {
  font-variant-numeric: tabular-nums;
  font-feature-settings: 'tnum' 1;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: 14px;
}

.chart-card {
  background: #fff;
  border-radius: 10px;
  border: 1px solid var(--gw-slate-200);
  box-shadow: var(--gw-card-shadow);
  overflow: hidden;
}

.chart-header {
  padding: 14px 16px;
  border-bottom: 1px solid var(--gw-slate-200);
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
  background: var(--gw-slate-50);
}

.chart-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--gw-slate-900);
}

.chart-header__hint {
  font-size: 12px;
  color: var(--gw-slate-500);
  font-weight: 500;
}

.chart-content {
  padding: 14px 16px 16px;
}

.time-series-chart {
  display: flex;
  position: relative;
  height: 200px;
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

.chart-bar {
  background: linear-gradient(to top, #4338ca, #6366f1);
  border-radius: 3px 3px 0 0;
  cursor: pointer;
  transition: filter 0.2s ease, box-shadow 0.2s ease;
  min-height: 2px;
  position: relative;
  z-index: 1;
  box-shadow: 0 1px 2px rgba(79, 70, 229, 0.25);
}

.bytes-bar {
  background: linear-gradient(to top, #047857, #10b981);
  box-shadow: 0 1px 2px rgba(5, 150, 105, 0.25);
}

.bytes-bar:hover {
  filter: brightness(1.06);
}

.chart-bar:hover {
  filter: brightness(1.06);
  box-shadow: 0 2px 6px rgba(79, 70, 229, 0.35);
}

.y-axis {
  width: 60px;
  position: relative;
  border-right: 1px solid var(--gw-slate-200);
  padding-right: 8px;
  display: flex;
  flex-direction: column-reverse;
  justify-content: space-between;
  padding-top: 5px;
  padding-bottom: 5px;
}

.y-tick {
  position: relative;
  right: auto;
  transform: translateY(50%);
  font-size: 10px;
  color: var(--gw-slate-500);
  white-space: nowrap;
  text-align: right;
  height: 0;
  line-height: 0;
}

.x-axis {
  position: absolute;
  bottom: 0;
  left: 12px;
  right: 8px;
  height: 35px;
  border-top: 1px solid var(--gw-slate-200);
  padding-top: 4px;
  overflow: visible;
}

.x-label {
  position: absolute;
  transform: translateX(-50%);
  font-size: 10px;
  color: var(--gw-slate-500);
  white-space: nowrap;
  font-weight: 500;
  pointer-events: none;
}

.status-chart,
.method-chart,
.service-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-bar,
.method-item,
.service-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.status-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-code {
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
}

.status-2xx {
  background: #ecfdf5;
  color: #047857;
}

.status-3xx {
  background: #eef2ff;
  color: #4338ca;
}

.status-4xx {
  background: #fffbeb;
  color: #b45309;
}

.status-5xx {
  background: #fff1f2;
  color: #be123c;
}

.status-count {
  font-size: 13px;
  color: var(--gw-slate-900);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.status-progress {
  height: 8px;
  background: var(--gw-slate-50);
  border-radius: 999px;
  overflow: hidden;
  border: 1px solid var(--gw-slate-200);
}

.status-progress-bar {
  height: 100%;
  transition: width 0.3s;
}

.method-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--gw-slate-900);
  margin-bottom: 4px;
}

.method-bar-container,
.service-bar-container {
  height: 32px;
  background: var(--gw-slate-50);
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  border: 1px solid var(--gw-slate-200);
}

.method-bar,
.service-bar {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 8px;
  transition: width 0.3s;
}

.method-count,
.service-count {
  font-size: 12px;
  color: #fff;
  font-weight: 500;
  white-space: nowrap;
}

.service-label {
  font-size: 14px;
  color: #262626;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.host-chart,
.path-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.host-item,
.path-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.host-rank,
.path-rank {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gw-slate-50);
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  color: var(--gw-slate-700);
  flex-shrink: 0;
  border: 1px solid var(--gw-slate-200);
}

.host-info,
.path-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.host-label,
.path-label {
  font-size: 13px;
  color: var(--gw-slate-900);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.host-bar-container,
.path-bar-container {
  height: 24px;
  background: var(--gw-slate-50);
  border-radius: 999px;
  overflow: hidden;
  position: relative;
  border: 1px solid var(--gw-slate-200);
}

.host-bar,
.path-bar {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 10px;
  background: linear-gradient(to right, #4f46e5, #818cf8);
  transition: width 0.25s ease;
}

.host-count,
.path-count {
  font-size: 11px;
  color: #fff;
  font-weight: 500;
  white-space: nowrap;
}

.path-bytes-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.path-bytes-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.path-bytes-rank {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gw-slate-50);
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  color: var(--gw-slate-700);
  flex-shrink: 0;
  border: 1px solid var(--gw-slate-200);
}

.path-bytes-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.path-bytes-label {
  font-size: 13px;
  color: var(--gw-slate-900);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.path-bytes-bar-container {
  height: 24px;
  background: var(--gw-slate-50);
  border-radius: 999px;
  overflow: hidden;
  position: relative;
  border: 1px solid var(--gw-slate-200);
}

.path-bytes-bar {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 10px;
  background: linear-gradient(to right, #059669, #34d399);
  transition: width 0.25s ease;
}

.path-bytes-count {
  font-size: 11px;
  color: #fff;
  font-weight: 500;
  white-space: nowrap;
}

@media (prefers-reduced-motion: reduce) {
  .stat-card,
  .chart-bar,
  .bytes-bar,
  .method-bar,
  .host-bar,
  .path-bar,
  .path-bytes-bar,
  .status-progress-bar {
    transition: none !important;
  }
}
</style>

