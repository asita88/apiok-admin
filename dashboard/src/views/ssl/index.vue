<template>
  <div class="main">
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item
        ><i
          style="color: #448ef7; font-size: 30px"
          class="iconfont icon-zhengshu-copy"
        />证书管理</a-breadcrumb-item
      >
    </a-breadcrumb>
    <a-divider style="margin: 10px 0" />

    <!-- 筛选 -->
    <div class="filter">
      <span>启用状态: </span>
      <a-select
        class="select"
        v-model:value="data.params.enable"
        placeholder="请选择"
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option :value="0">全部</a-select-option>
        <a-select-option :value="1">启用</a-select-option>
        <a-select-option :value="2">停用</a-select-option>
      </a-select>

      <span>CA提供商: </span>
      <a-select
        class="select"
        v-model:value="data.params.ca_provider"
        placeholder="请选择"
        allow-clear
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option value="letsencrypt">Let's Encrypt</a-select-option>
        <a-select-option value="manual">手动上传</a-select-option>
      </a-select>

      <span>加密算法: </span>
      <a-select
        class="select"
        v-model:value="data.params.key_algorithm"
        placeholder="请选择"
        allow-clear
        @change="fn.paramsChange(data.params)"
      >
        <a-select-option value="rsa2048">RSA 2048</a-select-option>
        <a-select-option value="rsa3072">RSA 3072</a-select-option>
        <a-select-option value="rsa4096">RSA 4096</a-select-option>
        <a-select-option value="ecdsa_p256">ECDSA P-256</a-select-option>
        <a-select-option value="ecdsa_p384">ECDSA P-384</a-select-option>
        <a-select-option value="ecdsa_p521">ECDSA P-521</a-select-option>
      </a-select>

      <a-input
        v-model:value="data.params.issuer"
        placeholder="颁发者"
        allow-clear
        style="width: 140px; margin-right: 10px"
        @pressEnter="fn.paramsChange(data.params)"
      />

      <a-input-search
        class="search"
        v-model:value="data.params.search"
        placeholder="搜索域名/ID"
        enter-button
        @search="fn.paramsChange(data.params)"
        @pressEnter="fn.paramsChange(data.params)"
      />

      <!-- 新增 -->
      <a-button type="primary" @click="fn.drawerOperate(null, drawer.typeRouter)" style="margin-right: 10px">
        <i class="iconfont icon-addNode" />新增证书</a-button
      >
      <!-- Let's Encrypt 申请 -->
      <a-button type="primary" @click="fn.drawerOperateLetsEncrypt()">
        <i class="iconfont icon-addNode" />申请 Let's Encrypt 证书</a-button
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
      <!-- 数据处理 -->
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'res_id'">
          {{ record.res_id }}
        </template>

        <template v-if="column.dataIndex === 'sni'">
          {{ record.sni }}
        </template>

        <template v-if="column.dataIndex === 'ca_provider'">
          {{ record.ca_provider }}
        </template>

        <template v-if="column.dataIndex === 'key_algorithm'">
          {{ record.key_algorithm }}
        </template>

        <template v-if="column.dataIndex === 'issuer'">
          <a-tooltip v-if="record.issuer && record.issuer.length > 20" :title="record.issuer">
            {{ record.issuer }}
          </a-tooltip>
          <span v-else>{{ record.issuer }}</span>
        </template>

        <template v-if="column.dataIndex === 'expired_at'">
          <span v-if="data.newTimeStamp > record.expired_at" class="color-red">
            {{ record.expired_at }}
          </span>
          <span v-else>
            {{ record.expired_at }}
          </span>
        </template>

        <template v-if="column.dataIndex === 'enable'">
          <a-switch v-model:checked="record.enable" size="small" @click="fn.enableChange(record)" />
        </template>

        <template v-if="column.dataIndex === 'operation'">
          <span>
            <a @click="fn.drawerOperate(record.res_id)">
              <a-tooltip placement="topRight">
                <template #title> 编辑 </template>
                <span>
                  <i class="iconfont icon-xiugai" />
                </span>
              </a-tooltip>
              <a-divider type="vertical" />
            </a>

            <a-popconfirm
              placement="topRight"
              title="确认要删除证书?"
              ok-text="是"
              cancel-text="否"
              @confirm="fn.deleteFunc(record)"
            >
              <a class="color-red a-delete">
                <i class="iconfont icon-shanchu" />
              </a>
            </a-popconfirm>
          </span>
        </template>
      </template>
    </a-table>

    <!-- 分页 -->
    <a-config-provider :locale="zh_CN">
      <a-pagination
        class="page"
        show-quick-jumper
        show-size-changer
        :total="data.listCount"
        v-model:current="data.params.page"
        @showSizeChange="fn.showSizeChange"
        :show-total="(total, range) => `当前${range[0]}-${range[1]}条，共${total}条`"
        @change="fn.pageChange"
      />
    </a-config-provider>
  </div>

  <!-- 抽屉 -->
  <a-drawer
    v-model:visible="drawer.visible"
    class="custom-class"
    placement="right"
    :title="drawer.title"
    :destroyOnClose="drawer.destroyOnClose"
    :width="drawer.width"
    @after-visible-change="fn.afterVisibleChange"
  >
    <!-- 动态组件完成路由和插件抽屉的展示 -->
    <component
      :is="drawer.componentName"
      :currentResId="drawer.currentResId"
      @componentCloseDrawer="fn.componentCloseDrawer"
      @componentRefreshList="fn.componentRefreshList"
    />
  </a-drawer>
</template>

<script>
import zh_CN from 'ant-design-vue/lib/locale-provider/zh_CN'
import SslOperate from './operate'
import SslOperateLetsEncrypt from './operate-letsencrypt'
import { reactive, ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { formatDate } from '@/utils/time'
import { $sslList, $sslEnable, $sslDelete } from '@/api'

export default {
  components: { SslOperate, SslOperateLetsEncrypt },

  setup() {
    // 初始化——路由列表
    onMounted(async () => {
      getList(data.params)
    })

    // 数据变量
    const data = reactive({
      params: reactive({
        enable: 0,
        ca_provider: undefined,
        key_algorithm: undefined,
        issuer: undefined,
        search: '',
        page: 1,
        page_size: 10
      }),
      columns: reactive([]), // 路由列表表头字段
      list: ref([]), // 路由列表数据
      listCount: 0, // 路由列表总条数
      newTimeStamp: ''
    })

    // 抽屉变量
    const drawer = reactive({
      visible: ref(false), // 抽屉开关
      title: null, // 抽屉标题
      width: '45%', // 抽屉宽度
      currentResId: null, // 当前资源ID
      destroyOnClose: false, // 抽屉内组件销毁  false：不销毁（默认）  true：销毁
      componentName: null // 动态组件名称
    })

    data.columns = reactive([
      { title: 'ID', dataIndex: 'res_id', width: 140 },
      { title: '域名', dataIndex: 'sni', width: 150 },
      { title: 'CA提供商', dataIndex: 'ca_provider', width: 100 },
      { title: '加密算法', dataIndex: 'key_algorithm', width: 100 },
      { title: '颁发者', dataIndex: 'issuer', width: 180, ellipsis: true },
      { title: '过期时间', dataIndex: 'expired_at', width: 160 },
      { title: '启用', dataIndex: 'enable', width: 70 },
      { title: '操作', dataIndex: 'operation', width: 80 }
    ])

    const getList = async params => {
      const reqParams = { ...params }
      if (reqParams.enable === 0) delete reqParams.enable
      if (reqParams.ca_provider === undefined || reqParams.ca_provider === '') delete reqParams.ca_provider
      if (reqParams.key_algorithm === undefined || reqParams.key_algorithm === '') delete reqParams.key_algorithm
      if (reqParams.issuer === undefined || reqParams.issuer === '') delete reqParams.issuer
      let { code, data: dataList, msg } = await $sslList(reqParams)

      if (code != 0) {
        message.error(msg)
      } else {
        data.listCount = dataList.total
        data.newTimeStamp = formatDate(Math.floor(new Date().getTime() / 1000))
        let key = 0
        let tmpList = ref([])
        dataList.data.forEach(item => {
          tmpList.value.push({
            key: key++,
            res_id: item.res_id,
            sni: item.sni,
            ca_provider: item.ca_provider === 'letsencrypt' ? "Let's Encrypt" : (item.ca_provider === 'manual' ? '手动上传' : item.ca_provider || '-'),
            key_algorithm: item.key_algorithm || '-',
            issuer: item.issuer || '-',
            expired_at: formatDate(item.expired_at),
            enable: item.enable == 1 ? true : false
          })
        })

        data.list = tmpList
      }
    }
    // 分页——页码变化
    const pageChange = async page => {
      data.params.page = page

      getList(data.params)
    }

    // 分页——每页显示条数变化
    const showSizeChange = async (current, size) => {
      data.params.page = 1
      data.params.page_size = size

      getList(data.params)
    }

    // 参数变化筛选列表
    const paramsChange = async () => {
      data.params.page = 1

      getList(data.params)
    }

    // 开关状态变化
    const enableChange = async record => {
      let enableValue = record.enable == true ? 1 : 2

      let { code, msg } = await $sslEnable(record.res_id, enableValue)

      if (code != 0) {
        message.error(msg)
        if (record.enable == true) {
          record.enable = false
        } else {
          record.enable = true
        }
        return
      } else {
        message.success(msg)
        getList(data.params)
      }
    }

    // 配置删除
    const deleteFunc = async record => {
      // 待发布且开启状态不允许删除
      if (record.release === 2 && record.enable === 1) {
        message.error('请先【关闭启用状态】后删除!')
        return
      }

      let { code, msg } = await $sslDelete(record.res_id)

      if (code != 0) {
        message.error(msg)
        return
      } else {
        message.success(msg)
        getList()
      }
    }

    // 抽屉操作
    const drawerOperate = async resId => {
      drawer.currentResId = resId
      drawer.visible = true

      drawer.componentName = 'SslOperate'
      drawer.width = '50%'
      if (resId === null) {
        drawer.title = '新增证书'
      } else {
        drawer.title = '编辑证书'
      }
    }

    // Let's Encrypt 申请抽屉操作
    const drawerOperateLetsEncrypt = async () => {
      drawer.currentResId = null
      drawer.visible = true
      drawer.componentName = 'SslOperateLetsEncrypt'
      drawer.width = '50%'
      drawer.title = '申请 Let\'s Encrypt 证书'
    }
    // 抽屉状态变化后钩子函数
    const afterVisibleChange = async () => {
      // 抽屉关闭时销毁抽屉内的组件（这里后期可能需要进一步优化）
      if (drawer.visible == false) {
        drawer.destroyOnClose = true
      } else {
        drawer.destroyOnClose = false
      }
    }

    // 动态组件——关闭抽屉
    const componentCloseDrawer = async () => {
      drawer.visible = false
    }

    // 动态组件——刷新列表
    const componentRefreshList = async () => {
      drawer.visible = false
      getList()
    }

    // 定义函数
    const fn = reactive({
      enableChange,
      paramsChange,
      pageChange,
      showSizeChange,
      deleteFunc,
      drawerOperate,
      drawerOperateLetsEncrypt,
      afterVisibleChange,
      componentCloseDrawer,
      componentRefreshList
    })

    return {
      zh_CN,
      data,
      drawer,
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
  font-size: 20px;
}
.select {
  width: 130px;
  margin-right: 10px;
}
.search {
  width: 260px;
  border-radius: 10px;
}
.filter {
  margin: 20px 0;
}
.filter button {
  float: right;
}
.table {
  margin: 20px 0px;
}
.edit-name {
  width: 80%;
  margin-right: 5px;
}
.page {
  margin: 20px 0px 20px 0px;
  display: flex;
  justify-content: right;
}
.iconfont {
  font-size: 18px;
}
.icon-addNode {
  font-size: 14px;
}
.edit-name-cancel:hover {
  color: #ff1744;
}
.icon-zhengshu-copy {
  margin-right: 10px;
}
</style>
