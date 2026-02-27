<template>
  <div class="main">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item
        ><i
          style="color: #448ef7; font-size: 30px"
          class="iconfont icon-rizhi"
        />操作日志</a-breadcrumb-item
      >
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <!-- 筛选 -->
    <div class="filter">
      <span>用户名: </span>
      <a-input
        v-model:value="data.params.username"
        placeholder="请输入用户名"
        style="width: 150px"
        @pressEnter="fn.paramsChange(data.params)"
      />

      <span>操作类型: </span>
      <a-select
        v-model:value="data.params.action"
        placeholder="请选择"
        style="width: 120px"
        allowClear
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option value="create">创建</a-select-option>
        <a-select-option value="update">更新</a-select-option>
        <a-select-option value="delete">删除</a-select-option>
        <a-select-option value="list">列表</a-select-option>
        <a-select-option value="info">详情</a-select-option>
      </a-select>

      <span>资源类型: </span>
      <a-select
        v-model:value="data.params.resource_type"
        placeholder="请选择"
        style="width: 120px"
        allowClear
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option value="service">服务</a-select-option>
        <a-select-option value="router">路由</a-select-option>
        <a-select-option value="upstream">上游</a-select-option>
        <a-select-option value="user">用户</a-select-option>
        <a-select-option value="certificate">证书</a-select-option>
        <a-select-option value="plugin">插件</a-select-option>
      </a-select>

      <span>状态码: </span>
      <a-select
        v-model:value="data.params.status_code"
        placeholder="请选择"
        style="width: 120px"
        allowClear
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option :value="200">200 成功</a-select-option>
        <a-select-option :value="400">400 错误</a-select-option>
        <a-select-option :value="401">401 未授权</a-select-option>
        <a-select-option :value="500">500 服务器错误</a-select-option>
      </a-select>

      <a-input-search
        class="search"
        v-model:value="data.params.search"
        placeholder="搜索路径、资源ID或IP"
        enter-button
        @search="fn.paramsChange(data.params)"
        @pressEnter="fn.paramsChange(data.params)"
      />
    </div>

    <!-- 列表 -->
    <a-table
      class="table"
      size="small"
      bordered
      :columns="data.columns"
      :data-source="data.list"
      :pagination="false"
      :scroll="{ x: 'max-content' }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status_code'">
          <a-tag :color="record.status_code >= 400 ? 'red' : 'green'">
            {{ record.status_code }}
          </a-tag>
        </template>

        <template v-if="column.dataIndex === 'action'">
          <a-tag :color="fn.getActionColor(record.action)">
            {{ fn.getActionName(record.action) }}
          </a-tag>
        </template>

        <template v-if="column.dataIndex === 'resource_type'">
          <a-tag>{{ fn.getResourceTypeName(record.resource_type) }}</a-tag>
        </template>

        <template v-if="column.dataIndex === 'request_data'">
          <a-button type="link" size="small" @click="fn.showDetail(record, 'request')">
            查看请求
          </a-button>
        </template>

        <template v-if="column.dataIndex === 'response_data'">
          <a-button type="link" size="small" @click="fn.showDetail(record, 'response')">
            查看响应
          </a-button>
        </template>

        <template v-if="column.dataIndex === 'error_message'">
          <span v-if="record.error_message" class="color-red">
            {{ record.error_message }}
          </span>
          <span v-else>-</span>
        </template>
      </template>
    </a-table>

    <!-- 分页 -->
    <a-pagination
      v-model:current="data.params.page"
      v-model:page-size="data.params.page_size"
      :total="data.total"
      show-size-changer
      show-quick-jumper
      :show-total="total => `共 ${total} 条`"
      @change="fn.paramsChange(data.params)"
      @showSizeChange="fn.paramsChange(data.params)"
    />

    <!-- 详情弹窗 -->
    <a-modal
      v-model:visible="data.detailVisible"
      :title="data.detailTitle"
      width="800px"
      :footer="null"
    >
      <pre class="json-view">{{ data.detailContent }}</pre>
    </a-modal>
  </div>
</template>

<script>
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import { $logList } from '@/api'

export default {
  setup() {
    const data = reactive({
      list: [],
      total: 0,
      params: {
        page: 1,
        page_size: 10,
        username: '',
        action: '',
        resource_type: '',
        status_code: 0,
        search: ''
      },
      columns: [
        {
          title: 'ID',
          dataIndex: 'res_id',
          width: 180
        },
        {
          title: '用户名',
          dataIndex: 'username',
          width: 120
        },
        {
          title: '操作类型',
          dataIndex: 'action',
          width: 100
        },
        {
          title: '资源类型',
          dataIndex: 'resource_type',
          width: 100
        },
        {
          title: '资源ID',
          dataIndex: 'resource_id',
          width: 150
        },
        {
          title: '请求方法',
          dataIndex: 'method',
          width: 80
        },
        {
          title: '请求路径',
          dataIndex: 'path',
          width: 250
        },
        {
          title: 'IP地址',
          dataIndex: 'ip',
          width: 120
        },
        {
          title: '状态码',
          dataIndex: 'status_code',
          width: 80
        },
        {
          title: '请求数据',
          dataIndex: 'request_data',
          width: 100
        },
        {
          title: '响应数据',
          dataIndex: 'response_data',
          width: 100
        },
        {
          title: '错误信息',
          dataIndex: 'error_message',
          width: 200
        },
        {
          title: '操作时间',
          dataIndex: 'created_at',
          width: 180
        }
      ],
      detailVisible: false,
      detailTitle: '',
      detailContent: ''
    })

    const fn = {
      paramsChange: async params => {
        const queryParams = { ...params }
        if (queryParams.status_code === 0) {
          delete queryParams.status_code
        }
        if (!queryParams.username) {
          delete queryParams.username
        }
        if (!queryParams.action) {
          delete queryParams.action
        }
        if (!queryParams.resource_type) {
          delete queryParams.resource_type
        }
        if (!queryParams.search) {
          delete queryParams.search
        }

        const { code, msg, data: resData } = await $logList(queryParams)
        if (code === 0) {
          data.list = resData.data || []
          data.total = resData.total || 0
        } else {
          message.error(msg || '获取日志列表失败')
        }
      },

      showDetail: (record, type) => {
        data.detailTitle = type === 'request' ? '请求数据' : '响应数据'
        const content = type === 'request' ? record.request_data : record.response_data
        try {
          const jsonObj = JSON.parse(content)
          data.detailContent = JSON.stringify(jsonObj, null, 2)
        } catch {
          data.detailContent = content || '无数据'
        }
        data.detailVisible = true
      },

      getActionName: action => {
        const actionMap = {
          create: '创建',
          update: '更新',
          delete: '删除',
          list: '列表',
          info: '详情'
        }
        return actionMap[action] || action
      },

      getActionColor: action => {
        const colorMap = {
          create: 'blue',
          update: 'orange',
          delete: 'red',
          list: 'green',
          info: 'cyan'
        }
        return colorMap[action] || 'default'
      },

      getResourceTypeName: type => {
        const typeMap = {
          service: '服务',
          router: '路由',
          upstream: '上游',
          user: '用户',
          certificate: '证书',
          plugin: '插件',
          cluster_node: '集群节点',
          letsencrypt: 'Let\'s Encrypt'
        }
        return typeMap[type] || type
      }
    }

    fn.paramsChange(data.params)

    return {
      data,
      fn
    }
  }
}
</script>

<style lange="scss" scoped>
.main {
  padding: 10px;
}

.breadcrumb {
  margin-bottom: 10px;
}

.filter {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 10px;
  flex-wrap: wrap;
}

.filter span {
  white-space: nowrap;
}

.search {
  width: 250px;
}

.table {
  margin-bottom: 16px;
}

.color-red {
  color: #ff4d4f;
}

.json-view {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  max-height: 500px;
  overflow: auto;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>

