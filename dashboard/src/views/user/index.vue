<template>
  <div class="main">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item
        ><i
          style="color: #448ef7; font-size: 30px"
          class="iconfont icon-user"
        />用户管理</a-breadcrumb-item
      >
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <!-- 筛选 -->
    <div class="filter">
      <a-input-search
        class="search"
        v-model:value="data.params.search"
        placeholder="搜索用户名或邮箱"
        enter-button
        @search="fn.paramsChange(data.params)"
        @pressEnter="fn.paramsChange(data.params)"
      />

      <!-- 新增 -->
      <a-button type="primary" @click="fn.drawerOperate(null)"
        ><i class="iconfont icon-tianjia" />新增用户</a-button
      >
    </div>

    <!-- 列表 -->
    <a-table
      class="table"
      size="small"
      bordered
      :columns="data.columns"
      :data-source="data.list"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'res_id'">
          <a @click="fn.drawerOperate(record.res_id)">
            {{ record.res_id }}
          </a>
        </template>

        <template v-if="column.dataIndex === 'action'">
          <a-space>
            <a @click="fn.drawerOperate(record.res_id)">编辑</a>
            <a-popconfirm
              title="确定要删除这个用户吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="fn.delete(record.res_id)"
            >
              <a class="color-red">删除</a>
            </a-popconfirm>
          </a-space>
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

    <!-- 抽屉 -->
    <a-drawer
      :title="data.drawerTitle"
      :width="600"
      :visible="data.drawerVisible"
      :body-style="{ paddingBottom: '80px' }"
      @close="fn.drawerClose"
    >
      <user-operate
        v-if="data.drawerVisible"
        :res-id="data.drawerResId"
        @close="fn.drawerClose"
        @refresh="fn.paramsChange(data.params)"
      />
    </a-drawer>
  </div>
</template>

<script>
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import { $userList, $userDelete } from '@/api'
import UserOperate from './operate.vue'

export default {
  components: {
    UserOperate
  },
  setup() {
    const data = reactive({
      list: [],
      total: 0,
      params: {
        page: 1,
        page_size: 10,
        search: ''
      },
      columns: [
        {
          title: '用户ID',
          dataIndex: 'res_id',
          width: 200
        },
        {
          title: '用户名',
          dataIndex: 'name',
          width: 150
        },
        {
          title: '邮箱',
          dataIndex: 'email',
          width: 200
        },
        {
          title: '创建时间',
          dataIndex: 'created_at',
          width: 180
        },
        {
          title: '更新时间',
          dataIndex: 'updated_at',
          width: 180
        },
        {
          title: '操作',
          dataIndex: 'action',
          width: 120
        }
      ],
      drawerVisible: false,
      drawerResId: null,
      drawerTitle: '新增用户'
    })

    const fn = {
      paramsChange: async params => {
        const { code, msg, data: resData } = await $userList(params)
        if (code === 0) {
          data.list = resData.data || []
          data.total = resData.total || 0
        } else {
          message.error(msg || '获取用户列表失败')
        }
      },

      drawerOperate: (resId, type) => {
        data.drawerResId = resId
        data.drawerTitle = resId ? '编辑用户' : '新增用户'
        data.drawerVisible = true
      },

      drawerClose: () => {
        data.drawerVisible = false
        data.drawerResId = null
      },

      delete: async resId => {
        const { code, msg } = await $userDelete(resId)
        if (code === 0) {
          message.success('删除成功')
          fn.paramsChange(data.params)
        } else {
          message.error(msg || '删除失败')
        }
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
}

.search {
  width: 300px;
}

.table {
  margin-bottom: 16px;
}

.color-red {
  color: #ff4d4f;
}
</style>

