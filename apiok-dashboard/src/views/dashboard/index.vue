<template>
  <div class="main">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>
        <i style="color: #448ef7; font-size: 30px" class="iconfont icon-rizhi"></i>大盘
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <div class="dashboard-content">
      <div class="time-filter">
        <a-range-picker
          v-model:value="data.timeRange"
          show-time
          format="YYYY-MM-DD HH:mm:ss"
          :placeholder="['开始时间', '结束时间']"
          @change="fn.onTimeRangeChange"
          style="width: 400px"
        />
        <a-button type="primary" @click="fn.quickTimeRange('1h')" style="margin-left: 10px">
          最近1小时
        </a-button>
        <a-button @click="fn.quickTimeRange('24h')" style="margin-left: 10px">
          最近24小时
        </a-button>
        <a-button @click="fn.quickTimeRange('7d')" style="margin-left: 10px">
          最近7天
        </a-button>
        <a-button @click="fn.refresh" style="margin-left: 10px" :loading="data.loading">
          刷新
        </a-button>
      </div>

      <a-spin :spinning="data.loading">
        <div class="stats-grid" v-if="data.aggregation">
          <div class="stat-card">
            <div class="stat-icon" style="background: #1890ff">
              <i class="iconfont icon-rizhi"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">总请求数</div>
              <div class="stat-value">{{ fn.formatNumber(data.aggregation.total_requests) }}</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon" style="background: #52c41a">
              <i class="iconfont icon-shijian"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">平均响应时间</div>
              <div class="stat-value">{{ fn.formatNumber(data.aggregation.avg_response_time, 2) }}ms</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon" style="background: #faad14">
              <i class="iconfont icon-cuowu"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">错误数</div>
              <div class="stat-value">{{ fn.formatNumber(data.aggregation.error_count) }}</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon" style="background: #f5222d">
              <i class="iconfont icon-baifenbi"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">错误率</div>
              <div class="stat-value">{{ fn.formatNumber(data.aggregation.error_rate * 100, 2) }}%</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon" style="background: #722ed1">
              <i class="iconfont icon-shujuliang"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">总流量</div>
              <div class="stat-value">{{ fn.formatBytes(data.aggregation.total_bytes_sent) }}</div>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon" style="background: #13c2c2">
              <i class="iconfont icon-shijian"></i>
            </div>
            <div class="stat-content">
              <div class="stat-label">最大响应时间</div>
              <div class="stat-value">{{ fn.formatNumber(data.aggregation.max_response_time, 2) }}ms</div>
            </div>
          </div>
        </div>

        <div class="charts-grid" v-if="data.aggregation">
          <div class="chart-card">
            <div class="chart-header">
              <h3>请求趋势</h3>
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
import { reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
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
        return '#52c41a'
      case 3:
        return '#1890ff'
      case 4:
        return '#faad14'
      case 5:
        return '#f5222d'
      default:
        return '#d9d9d9'
    }
  },

  getStatusColorByRange: (statusRange) => {
    switch (statusRange) {
      case '2xx':
        return '#52c41a'
      case '3xx':
        return '#1890ff'
      case '4xx':
        return '#faad14'
      case '5xx':
        return '#f5222d'
      default:
        return '#d9d9d9'
    }
  },

  getMethodColor: (method) => {
    const colors = {
      GET: '#1890ff',
      POST: '#52c41a',
      PUT: '#faad14',
      DELETE: '#f5222d',
      PATCH: '#722ed1'
    }
    return colors[method] || '#d9d9d9'
  }
}

onMounted(() => {
  fn.quickTimeRange('24h')
})
</script>

<style scoped>
.main {
  padding: 10px;
  height: calc(100vh - 64px - 20px);
  max-height: calc(100vh - 64px - 20px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
}

.breadcrumb {
  flex-shrink: 0;
}

.dashboard-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.time-filter {
  margin-bottom: 20px;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: #fff;
  font-size: 24px;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #8c8c8c;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 16px;
}

.chart-card {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.chart-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.chart-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.chart-content {
  padding: 16px;
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
  background: linear-gradient(to top, #1890ff, #40a9ff);
  border-radius: 2px 2px 0 0;
  cursor: pointer;
  transition: all 0.2s;
  min-height: 2px;
  position: relative;
  z-index: 1;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.bytes-bar {
  background: linear-gradient(to top, #52c41a, #73d13d);
}

.bytes-bar:hover {
  background: linear-gradient(to top, #389e0d, #52c41a);
}

.chart-bar:hover {
  background: linear-gradient(to top, #096dd9, #1890ff);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  transform: translateY(-1px);
}

.y-axis {
  width: 60px;
  position: relative;
  border-right: 1px solid #e8e8e8;
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
  color: #999;
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
  border-radius: 4px;
  font-size: 12px;
}

.status-2xx {
  background: #f6ffed;
  color: #52c41a;
}

.status-3xx {
  background: #e6f7ff;
  color: #1890ff;
}

.status-4xx {
  background: #fffbe6;
  color: #faad14;
}

.status-5xx {
  background: #fff1f0;
  color: #f5222d;
}

.status-count {
  font-size: 14px;
  color: #262626;
  font-weight: 500;
}

.status-progress {
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.status-progress-bar {
  height: 100%;
  transition: width 0.3s;
}

.method-label {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  margin-bottom: 4px;
}

.method-bar-container,
.service-bar-container {
  height: 32px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
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
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #595959;
  flex-shrink: 0;
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
  color: #262626;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.host-bar-container,
.path-bar-container {
  height: 24px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.host-bar,
.path-bar {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 8px;
  background: linear-gradient(to right, #1890ff, #40a9ff);
  transition: width 0.3s;
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
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #595959;
  flex-shrink: 0;
}

.path-bytes-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.path-bytes-label {
  font-size: 13px;
  color: #262626;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.path-bytes-bar-container {
  height: 24px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.path-bytes-bar {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 8px;
  background: linear-gradient(to right, #52c41a, #73d13d);
  transition: width 0.3s;
}

.path-bytes-count {
  font-size: 11px;
  color: #fff;
  font-weight: 500;
  white-space: nowrap;
}
</style>

