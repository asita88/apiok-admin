<template>
  <div class="main">
    <a-form
      :model="data.formData"
      class="form"
      :label-col="{ style: { width: '100px' } }"
      :wrapper-col="{ span: 18 }"
      @finish="fn.onSubmit"
    >
      <a-form-item label="路由名称：" name="router_name" :rules="schemaRouter.name">
        <a-input v-model:value="data.formData.router_name" />
      </a-form-item>

      <a-form-item label="所属服务：" name="service_res_id" :rules="schemaRouter.service_res_id">
        <a-select
          class="select"
          ref="select"
          :field-names="{
            label: 'name',
            value: 'res_id'
          }"
          show-search
          v-model:value="data.formData.service_res_id"
          placeholder="请选择"
          :filter-option="fn.filterOption"
          :options="data.serviceList"
        ></a-select>
      </a-form-item>

      <a-form-item label="路由路径：" name="router_path" :rules="schemaRouter.router_path">
        <a-input v-model:value="data.formData.router_path" />
      </a-form-item>

      <a-form-item name="priority">
        <template #label>
          <span>
            优先级
            <a-tooltip title="同路径多条规则时数值越大越先匹配">
              <i class="iconfont icon-help color-orange" />
            </a-tooltip>
          </span>
        </template>
        <a-input-number
          v-model:value="data.formData.priority"
          :min="0"
          :precision="0"
          style="width: 100%"
        />
      </a-form-item>

      <a-form-item label="请求方法：" name="request_methods" :rules="schemaRouter.request_methods">
        <a-checkbox-group
          v-model:value="data.formData.request_methods"
          :options="data.methodList"
        />
      </a-form-item>

      <a-form-item label="upstream：" name="upstream_res_id">
        <a-select
          class="select"
          ref="select"
          v-model:value="data.formData.upstream_res_id"
          placeholder="请选择"
        >
          <a-select-option v-for="item in data.upstreamList" :value="item.res_id">{{
            item.name
          }}</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="启用：">
        <a-switch v-model:checked="data.formData.enable" />
      </a-form-item>

      <a-divider orientation="left">改写规则</a-divider>

      <a-form-item label="请求改写：">
        <div>
          <a-switch v-model:checked="data.rewriteProxyRewriteEnabled" />
          <div v-if="data.rewriteProxyRewriteEnabled" class="rewrite-nested">
            <a-form-item
              label="协议"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-radio-group v-model:value="data.rewritePrProtocol" button-style="solid" size="small">
                <a-radio-button value="keep">保持原样</a-radio-button>
                <a-radio-button value="http">HTTP</a-radio-button>
                <a-radio-button value="https">HTTPS</a-radio-button>
              </a-radio-group>
            </a-form-item>

            <a-form-item
              label="路径改写"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-radio-group v-model:value="data.rewritePrPath" button-style="solid" size="small">
                <a-radio-button value="keep">保持原样</a-radio-button>
                <a-radio-button value="static">静态改写</a-radio-button>
                <a-radio-button value="regex">正则改写</a-radio-button>
              </a-radio-group>
              <div v-if="data.rewritePrPath === 'static'" style="margin-top: 10px">
                <div
                  v-for="(row, ri) in data.rewritePrPathStaticRows"
                  :key="'rw-static-' + ri"
                  style="margin-bottom: 10px"
                >
                  <a-row :gutter="8" align="middle">
                    <a-col :span="11">
                      <a-input v-model:value="row.from" placeholder="源路径（须与请求 URI 完全一致）" />
                    </a-col>
                    <a-col :span="11">
                      <a-input v-model:value="row.to" placeholder="目标路径" />
                    </a-col>
                    <a-col :span="2">
                      <a-button
                        type="link"
                        danger
                        @click="fn.removeRewriteStaticRow(ri)"
                        style="padding: 0"
                        :disabled="data.rewritePrPathStaticRows.length <= 1"
                      >
                        <i class="iconfont icon-jian" />
                      </a-button>
                    </a-col>
                  </a-row>
                </div>
                <a-button type="dashed" @click="fn.addRewriteStaticRow" style="width: 100%">
                  <i class="iconfont icon-tianjia" /> 添加静态改写
                </a-button>
              </div>
              <div v-if="data.rewritePrPath === 'regex'" style="margin-top: 10px">
                <div
                  v-for="(row, ri) in data.rewritePrRegexRows"
                  :key="'rw-regex-' + ri"
                  style="margin-bottom: 10px"
                >
                  <a-row :gutter="8" align="middle">
                    <a-col :span="11">
                      <a-input v-model:value="row.pattern" placeholder="匹配正则，如 ^/api/(.*)" />
                    </a-col>
                    <a-col :span="11">
                      <a-input v-model:value="row.replacement" placeholder="替换为，如 /$1" />
                    </a-col>
                    <a-col :span="2">
                      <a-button
                        type="link"
                        danger
                        @click="fn.removeRewriteRegexRow(ri)"
                        style="padding: 0"
                        :disabled="data.rewritePrRegexRows.length <= 1"
                      >
                        <i class="iconfont icon-jian" />
                      </a-button>
                    </a-col>
                  </a-row>
                </div>
                <a-button type="dashed" @click="fn.addRewriteRegexRow" style="width: 100%">
                  <i class="iconfont icon-tianjia" /> 添加正则改写
                </a-button>
              </div>
            </a-form-item>

            <a-form-item
              label="域名改写"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-radio-group v-model:value="data.rewritePrHost" button-style="solid" size="small">
                <a-radio-button value="keep">保持原样</a-radio-button>
                <a-radio-button value="static">静态改写</a-radio-button>
              </a-radio-group>
              <div v-if="data.rewritePrHost === 'static'" style="margin-top: 10px">
                <a-input
                  v-model:value="data.rewritePrHostStatic"
                  placeholder="主机名或 IP，勿填协议与路径"
                />
              </div>
            </a-form-item>

            <a-form-item
              label="方法改写"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-select
                v-model:value="data.rewritePrMethod"
                allow-clear
                placeholder="不改写"
                style="width: 100%; max-width: 280px"
              >
                <a-select-option v-for="m in data.rewriteMethodOptions" :key="m" :value="m">{{ m }}</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item
              label="请求头改写"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <div style="width: 100%">
                <div v-for="(h, hi) in data.rewriteProxyRewriteHeaders" :key="'rw-h-' + hi" style="margin-bottom: 10px">
                  <a-row :gutter="8" align="middle">
                    <a-col :span="10">
                      <a-input v-model:value="h.key" placeholder="参数名称" />
                    </a-col>
                    <a-col :span="12">
                      <a-input v-model:value="h.value" placeholder="参数值" />
                    </a-col>
                    <a-col :span="2">
                      <a-button type="link" danger @click="fn.removeRewriteProxyHeader(hi)" style="padding: 0">
                        <i class="iconfont icon-jian" />
                      </a-button>
                    </a-col>
                  </a-row>
                </div>
                <a-button type="dashed" @click="fn.addRewriteProxyHeader" style="width: 100%">
                  <i class="iconfont icon-tianjia" /> 新建
                </a-button>
              </div>
            </a-form-item>
          </div>
        </div>
      </a-form-item>

      <a-form-item label="路径替换：">
        <div>
          <a-switch v-model:checked="data.rewriteReplaceEnabled" />
          <div v-if="data.rewriteReplaceEnabled" class="rewrite-nested">
            <a-form-item
              label="查找"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-input v-model:value="data.rewriteReplaceFrom" placeholder="URI 中的原字符串，如 /old-api" />
            </a-form-item>
            <a-form-item
              label="替换为"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-input v-model:value="data.rewriteReplaceTo" placeholder="可为空" />
            </a-form-item>
          </div>
        </div>
      </a-form-item>

      <a-form-item label="重定向：">
        <div>
          <a-switch v-model:checked="data.rewriteRedirectEnabled" />
          <div v-if="data.rewriteRedirectEnabled" class="rewrite-nested">
            <a-form-item
              label="HTTPS"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <div>
                <a-switch v-model:checked="data.rewriteRedirect.http_to_https" />
                <span style="margin-left: 8px; color: #999; font-size: 12px">HTTP 跳转 HTTPS</span>
              </div>
            </a-form-item>
            <a-form-item
              label="重定向 URI"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-input v-model:value="data.rewriteRedirect.uri" placeholder="可选" />
            </a-form-item>
            <a-form-item
              label="状态码"
              :label-col="{ style: { width: '100px' } }"
              :wrapper-col="{ span: 18 }"
              class="rewrite-inner-item"
            >
              <a-input-number
                v-model:value="data.rewriteRedirect.ret_code"
                :min="300"
                :max="399"
                placeholder="如 302"
                style="width: 100%; max-width: 200px"
              />
            </a-form-item>
          </div>
        </div>
      </a-form-item>

      <a-form-item
        v-if="data.rewriteRulesExtraKeys.length"
        :label-col="{ style: { width: '100px' } }"
        :wrapper-col="{ span: 18 }"
      >
        <template #label></template>
        <div class="rewrite-extra-tip">
          另含其它插件配置：{{ data.rewriteRulesExtraKeys.join('、') }}（已保留，保存时一并提交）
        </div>
      </a-form-item>



      <a-divider orientation="left">高级配置</a-divider>

      <a-form-item label="请求体大小限制：" name="client_max_body_size">
        <a-input
          v-model:value="data.formData.client_max_body_size"
          placeholder="例如：100M、1G、500K（0表示无限制）"
          style="width: 100%"
        />
        <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
          支持单位：k(千字节), m(兆字节), g(千兆字节)，例如：100M、1G、500K，0表示无限制
        </div>
      </a-form-item>

      <a-form-item label="分块传输编码：">
        <div>
          <a-switch v-model:checked="data.formData.chunked_transfer_encoding" />
          <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
            启用后使用 Transfer-Encoding: chunked，禁用则使用 Content-Length
          </div>
        </div>
      </a-form-item>

      <a-form-item label="代理缓冲：">
        <div>
          <a-switch v-model:checked="data.formData.proxy_buffering" />
          <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
            启用缓冲可减少与上游的连接时间，禁用则支持流式传输
          </div>
        </div>
      </a-form-item>

      <a-form-item label="代理缓存：">
        <a-switch v-model:checked="data.proxyCacheEnabled" @change="fn.onProxyCacheChange" />
        <div v-if="data.proxyCacheEnabled" style="margin-top: 16px; padding-left: 0">
          <a-form-item 
            label="缓存键：" 
            :label-col="{ style: { width: '100px' } }"
            :wrapper-col="{ span: 18 }"
            style="margin-bottom: 16px"
          >
            <a-input
              v-model:value="data.proxyCacheConfig.cache_key"
              placeholder="$scheme$proxy_host$request_uri"
            />
          </a-form-item>
          <a-form-item 
            label="缓存有效期：" 
            :label-col="{ style: { width: '100px' } }"
            :wrapper-col="{ span: 18 }"
            style="margin-bottom: 16px"
          >
            <a-input
              v-model:value="data.proxyCacheConfig.cache_valid"
              placeholder="200 302 10m 404 1m any 5m"
            />
            <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
              格式：状态码 时间，如：200 302 10m 表示200和302状态码缓存10分钟
            </div>
          </a-form-item>
          <a-form-item 
            label="绕过缓存条件：" 
            :label-col="{ style: { width: '100px' } }"
            :wrapper-col="{ span: 18 }"
            style="margin-bottom: 16px"
          >
            <a-textarea
              v-model:value="data.proxyCacheBypassText"
              placeholder="$http_pragma, $http_authorization"
              :rows="2"
            />
            <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
              多个条件用逗号分隔，如：$http_pragma, $http_authorization
            </div>
          </a-form-item>
          <a-form-item 
            label="不缓存条件：" 
            :label-col="{ style: { width: '100px' } }"
            :wrapper-col="{ span: 18 }"
            style="margin-bottom: 16px"
          >
            <a-textarea
              v-model:value="data.proxyNoCacheText"
              placeholder="$http_pragma"
              :rows="2"
            />
            <div style="color: #999; font-size: 12px; margin-top: 6px; line-height: 1.5">
              多个条件用逗号分隔，如：$http_pragma
            </div>
          </a-form-item>
        </div>
      </a-form-item>

      <a-form-item label="代理请求头：">
        <div style="width: 100%">
          <div v-for="(header, index) in data.proxySetHeaders" :key="index" style="margin-bottom: 10px">
            <a-row :gutter="8" align="middle">
              <a-col :span="10">
                <a-input v-model:value="header.key" placeholder="请求头名称" />
              </a-col>
              <a-col :span="12">
                <a-input v-model:value="header.value" placeholder="请求头值（可使用$变量）" />
              </a-col>
              <a-col :span="2">
                <a-button type="link" danger @click="fn.removeProxyHeader(index)" style="padding: 0">
                  <i class="iconfont icon-jian" />
                </a-button>
              </a-col>
            </a-row>
          </div>
          <a-button type="dashed" @click="fn.addProxyHeader" style="width: 100%">
            <i class="iconfont icon-tianjia" /> 添加请求头
          </a-button>
          <div style="color: #999; font-size: 12px; margin-top: 8px; line-height: 1.5">
            可设置自定义请求头传递给上游服务器，支持Nginx变量如：$scheme, $host, $remote_addr
          </div>
        </div>
      </a-form-item>



      <a-form-item :wrapper-col="{ span: 10, offset: 16 }">
        <a-button type="primary" html-type="submit">保存</a-button>
        <a-button style="margin-left: 15px" @click="fn.cancel">取消</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { $serviceList, $upstreamNameList, $routerInfo, $routerAdd, $routerUpdate } from '@/api'
import { MethodOption } from '@/hooks'
import { schemaRouter } from '@/schema'

export default {
  props: {
    serviceResId: {
      String
    },
    currentResId: {
      String
    }
  },
  emits: ['componentCloseDrawer', 'componentRefreshList'],
  setup(props, { emit }) {
    // 初始化——服务详情数据
    onMounted(async () => {
      if (props.serviceResId !== null) {
        data.formData.service_res_id = props.serviceResId
      }

      if (props.currentResId !== null && props.serviceResId !== null) {
        getInfo(props.currentResId)
      }

      getServiceList(data.serviceParam)

      getUpstreamNameList()
    })

    // 定义变量
    const data = reactive({
      formData: {
        service_res_id: '',
        upstream_res_id: '',
        router_name: '',
        request_methods: [],
        router_path: '',
        priority: 0,
        enable: false,
        client_max_body_size: '',
        chunked_transfer_encoding: undefined,
        proxy_buffering: undefined,
        proxy_cache: undefined,
        proxy_set_header: undefined
      },
      proxyCacheEnabled: false,
      proxyCacheConfig: {
        cache_key: '',
        cache_valid: '',
        cache_bypass: [],
        no_cache: []
      },
      proxyCacheBypassText: '',
      proxyNoCacheText: '',
      proxySetHeaders: [],
      rewriteProxyRewriteEnabled: false,
      rewritePrProtocol: 'keep',
      rewritePrPath: 'keep',
      rewritePrPathStaticRows: [{ from: '', to: '' }],
      rewritePrRegexRows: [{ pattern: '', replacement: '' }],
      rewritePrHost: 'keep',
      rewritePrHostStatic: '',
      rewritePrMethod: undefined,
      rewriteMethodOptions: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS', 'TRACE'],
      rewriteProxyRewriteHeaders: [],
      rewriteRedirectEnabled: false,
      rewriteRedirect: {
        http_to_https: false,
        uri: '',
        ret_code: undefined
      },
      rewriteReplaceEnabled: false,
      rewriteReplaceFrom: '',
      rewriteReplaceTo: '',
      rewriteRulesOriginalExtras: {},
      rewriteRulesExtraKeys: [],
      serviceParam: reactive({
        page: 1,
        page_size: 1000 // 此处暂时不做轮询获取 暂定获取前1000条
      }),
      serviceList: ref([]), // 服务列表
      upstreamParam: reactive({
        enable: 1,
        release: 3,
        page: 1,
        page_size: 1000 // 此处暂时不做轮询获取 暂定获取前1000条
      }),
      upstreamList: reactive({}), // upstream列表
      methodList: MethodOption // 请求方法列表
    })

    // 获取服务列表
    const getServiceList = async params => {
      let { code, data: dataList, msg } = await $serviceList(params)

      if (code != 0) {
        message.error(msg)
      } else {
        let tmpList = ref([])
        dataList.data.forEach(item => {
          tmpList.value.push({
            res_id: item.res_id,
            name: item.name
          })
        })

        data.serviceList = tmpList
      }
    }

    // 下拉选择服务列表搜索
    const filterOption = (input, option) => {
      return option.name.toLowerCase().indexOf(input.toLowerCase()) >= 0
    }

    // 获取upstream列表
    const getUpstreamNameList = async () => {
      let { code, data: dataList, msg } = await $upstreamNameList()

      let tmpList = ref([])

      tmpList.value.push({
        res_id: '',
        name: '请选择'
      })

      if (code != 0) {
        message.error(msg)
      } else {
        dataList.forEach(item => {
          tmpList.value.push({
            res_id: item.res_id,
            name: item.name
          })
        })
      }

      data.upstreamList = tmpList
    }
    // 获取详情
    const getInfo = async (resId) => {
      let { code, data: dataInfo, msg } = await $routerInfo(resId)

      if (code !== 0) {
        message.error(msg)
        emit('componentCloseDrawer')
        return
      } else {
        data.formData.service_res_id = dataInfo.service_res_id
        data.formData.upstream_res_id = dataInfo.upstream_res_id
        data.formData.router_name = dataInfo.router_name
        data.formData.request_methods = dataInfo.request_methods
        data.formData.router_path = dataInfo.router_path
        data.formData.priority =
          dataInfo.priority !== undefined && dataInfo.priority !== null ? Number(dataInfo.priority) : 0
        data.formData.enable = dataInfo.enable == 1 ? true : false
        
        // 处理新字段
        if (dataInfo.client_max_body_size !== undefined && dataInfo.client_max_body_size !== null) {
          data.formData.client_max_body_size = dataInfo.client_max_body_size
        }
        if (dataInfo.chunked_transfer_encoding !== undefined && dataInfo.chunked_transfer_encoding !== null) {
          data.formData.chunked_transfer_encoding = dataInfo.chunked_transfer_encoding == 1
        }
        if (dataInfo.proxy_buffering !== undefined && dataInfo.proxy_buffering !== null) {
          data.formData.proxy_buffering = dataInfo.proxy_buffering == 1
        }
        
        // 处理代理缓存配置
        if (dataInfo.proxy_cache) {
          try {
            const proxyCache = typeof dataInfo.proxy_cache === 'string' 
              ? JSON.parse(dataInfo.proxy_cache) 
              : dataInfo.proxy_cache
            if (proxyCache && proxyCache.enabled) {
              data.proxyCacheEnabled = true
              data.proxyCacheConfig.cache_key = proxyCache.cache_key || ''
              data.proxyCacheConfig.cache_valid = proxyCache.cache_valid || ''
              data.proxyCacheBypassText = Array.isArray(proxyCache.cache_bypass) 
                ? proxyCache.cache_bypass.join(', ') 
                : ''
              data.proxyNoCacheText = Array.isArray(proxyCache.no_cache) 
                ? proxyCache.no_cache.join(', ') 
                : ''
            }
          } catch (e) {
            console.error('解析 proxy_cache 失败:', e)
          }
        }
        
        // 处理代理请求头配置
        if (dataInfo.proxy_set_header) {
          try {
            const proxySetHeader = typeof dataInfo.proxy_set_header === 'string' 
              ? JSON.parse(dataInfo.proxy_set_header) 
              : dataInfo.proxy_set_header
            if (proxySetHeader && typeof proxySetHeader === 'object') {
              data.proxySetHeaders = Object.keys(proxySetHeader).map(key => ({
                key: key,
                value: proxySetHeader[key]
              }))
            }
          } catch (e) {
            console.error('解析 proxy_set_header 失败:', e)
          }
        }

        resetRewriteFormFromApi(dataInfo.rewrite_rules)
      }
    }

    const resetRewriteFormFromApi = rr => {
      data.rewriteProxyRewriteEnabled = false
      data.rewritePrProtocol = 'keep'
      data.rewritePrPath = 'keep'
      data.rewritePrPathStaticRows = [{ from: '', to: '' }]
      data.rewritePrRegexRows = [{ pattern: '', replacement: '' }]
      data.rewritePrHost = 'keep'
      data.rewritePrHostStatic = ''
      data.rewritePrMethod = undefined
      data.rewriteProxyRewriteHeaders = []
      data.rewriteRedirectEnabled = false
      data.rewriteRedirect = {
        http_to_https: false,
        uri: '',
        ret_code: undefined
      }
      data.rewriteReplaceEnabled = false
      data.rewriteReplaceFrom = ''
      data.rewriteReplaceTo = ''
      data.rewriteRulesOriginalExtras = {}
      data.rewriteRulesExtraKeys = []

      if (!rr || typeof rr !== 'object' || Array.isArray(rr)) {
        return
      }

      data.rewriteRulesOriginalExtras = { ...rr }
      const pr = rr['proxy-rewrite']
      if (pr && typeof pr === 'object') {
        delete data.rewriteRulesOriginalExtras['proxy-rewrite']
        data.rewriteProxyRewriteEnabled = true
        const sch = pr.scheme
        if (sch === 'http') data.rewritePrProtocol = 'http'
        else if (sch === 'https') data.rewritePrProtocol = 'https'
        else data.rewritePrProtocol = 'keep'

        const pm = pr.path_mode
        if (pm === 'keep' || pm === 'static' || pm === 'regex') {
          data.rewritePrPath = pm
        } else if (Array.isArray(pr.regex_uris) && pr.regex_uris.length > 0) {
          data.rewritePrPath = 'regex'
        } else if (Array.isArray(pr.uri_pairs) && pr.uri_pairs.length > 0) {
          data.rewritePrPath = 'static'
        } else {
          data.rewritePrPath = 'keep'
        }

        if (data.rewritePrPath === 'static') {
          data.rewritePrPathStaticRows =
            Array.isArray(pr.uri_pairs) && pr.uri_pairs.length > 0
              ? pr.uri_pairs.map(pair => ({
                  from: String(pair[0] ?? ''),
                  to: String(pair[1] ?? '')
                }))
              : [{ from: '', to: '' }]
        } else if (data.rewritePrPath === 'regex') {
          data.rewritePrRegexRows =
            Array.isArray(pr.regex_uris) && pr.regex_uris.length > 0
              ? pr.regex_uris.map(pair => ({
                  pattern: String(pair[0] ?? ''),
                  replacement: String(pair[1] ?? '')
                }))
              : [{ pattern: '', replacement: '' }]
        }

        if (pr.host) {
          data.rewritePrHost = 'static'
          data.rewritePrHostStatic = pr.host
        } else {
          data.rewritePrHost = 'keep'
        }

        if (pr.method) {
          data.rewritePrMethod = String(pr.method).toUpperCase()
        }

        const setObj = pr.headers && pr.headers.set && typeof pr.headers.set === 'object' ? pr.headers.set : null
        if (setObj) {
          data.rewriteProxyRewriteHeaders = Object.keys(setObj).map(k => ({
            key: k,
            value: setObj[k] != null ? String(setObj[k]) : ''
          }))
        }
      }

      const rp = rr.replace
      if (rp && typeof rp === 'object') {
        delete data.rewriteRulesOriginalExtras.replace
        data.rewriteReplaceEnabled = true
        if (rp.from != null) data.rewriteReplaceFrom = String(rp.from)
        if (rp.to != null) data.rewriteReplaceTo = String(rp.to)
      }

      const rd = rr['redirect']
      if (rd && typeof rd === 'object') {
        delete data.rewriteRulesOriginalExtras['redirect']
        data.rewriteRedirectEnabled = true
        if (rd.http_to_https === true) data.rewriteRedirect.http_to_https = true
        if (rd.uri) data.rewriteRedirect.uri = rd.uri
        if (rd.ret_code != null && rd.ret_code !== '') {
          const n = Number(rd.ret_code)
          if (!Number.isNaN(n)) data.rewriteRedirect.ret_code = n
        }
      }

      data.rewriteRulesExtraKeys = Object.keys(data.rewriteRulesOriginalExtras)
    }

    const buildProxyRewritePlugin = () => {
      const o = {
        path_mode: data.rewritePrPath
      }

      if (data.rewritePrProtocol === 'http') o.scheme = 'http'
      else if (data.rewritePrProtocol === 'https') o.scheme = 'https'

      if (data.rewritePrPath === 'static') {
        const pairs = []
        for (const row of data.rewritePrPathStaticRows) {
          const from = String(row.from || '').trim()
          const to = String(row.to || '').trim()
          if (!from && !to) continue
          if ((from && !to) || (!from && to)) {
            message.error('每组静态改写需同时填写「源路径」与「目标路径」')
            return null
          }
          pairs.push([from, to])
        }
        o.uri_pairs = pairs
      } else if (data.rewritePrPath === 'regex') {
        const pairs = []
        for (const row of data.rewritePrRegexRows) {
          const p = String(row.pattern || '').trim()
          const r = String(row.replacement || '').trim()
          if (!p && !r) continue
          if ((p && !r) || (!p && r)) {
            message.error('每组正则改写需同时填写「匹配」与「替换」')
            return null
          }
          pairs.push([p, r])
        }
        o.regex_uris = pairs
      }

      if (data.rewritePrHost === 'static') {
        const h = data.rewritePrHostStatic.trim()
        if (h) o.host = h
      }

      if (data.rewritePrMethod) {
        const m = String(data.rewritePrMethod).trim().toUpperCase()
        if (m) o.method = m
      }

      const set = {}
      data.rewriteProxyRewriteHeaders.forEach(row => {
        const k = (row.key || '').trim()
        if (k && row.value != null && String(row.value).trim() !== '') {
          set[k] = String(row.value).trim()
        }
      })
      if (Object.keys(set).length > 0) {
        o.headers = { set }
      }
      return o
    }

    const buildRedirectPlugin = () => {
      const o = {}
      if (data.rewriteRedirect.http_to_https) {
        o.http_to_https = true
      }
      const uri = (data.rewriteRedirect.uri || '').trim()
      if (uri) o.uri = uri
      if (data.rewriteRedirect.ret_code != null && data.rewriteRedirect.ret_code !== '') {
        const n = Number(data.rewriteRedirect.ret_code)
        if (!Number.isNaN(n)) o.ret_code = n
      }
      return Object.keys(o).length > 0 ? o : null
    }

    const buildRewriteRulesPayload = () => {
      const out = { ...data.rewriteRulesOriginalExtras }

      if (data.rewriteProxyRewriteEnabled) {
        const pr = buildProxyRewritePlugin()
        if (pr === null) return null
        if (Object.keys(pr).length > 0) {
          out['proxy-rewrite'] = pr
        }
      }

      if (data.rewriteReplaceEnabled) {
        const from = (data.rewriteReplaceFrom || '').trim()
        if (!from) {
          message.error('路径替换需填写「查找」')
          return null
        }
        out.replace = {
          from,
          to: (data.rewriteReplaceTo || '').trim()
        }
      }

      if (data.rewriteRedirectEnabled) {
        const rd = buildRedirectPlugin()
        if (rd) {
          out['redirect'] = rd
        }
      }

      return out
    }

    const addRewriteProxyHeader = () => {
      data.rewriteProxyRewriteHeaders.push({ key: '', value: '' })
    }

    const removeRewriteProxyHeader = index => {
      data.rewriteProxyRewriteHeaders.splice(index, 1)
    }

    const addRewriteStaticRow = () => {
      data.rewritePrPathStaticRows.push({ from: '', to: '' })
    }

    const removeRewriteStaticRow = index => {
      if (data.rewritePrPathStaticRows.length <= 1) return
      data.rewritePrPathStaticRows.splice(index, 1)
    }

    const addRewriteRegexRow = () => {
      data.rewritePrRegexRows.push({ pattern: '', replacement: '' })
    }

    const removeRewriteRegexRow = index => {
      if (data.rewritePrRegexRows.length <= 1) return
      data.rewritePrRegexRows.splice(index, 1)
    }

    // 处理代理缓存配置变化
    const onProxyCacheChange = (checked) => {
      if (!checked) {
        data.proxyCacheConfig = {
          cache_key: '',
          cache_valid: '',
          cache_bypass: [],
          no_cache: []
        }
        data.proxyCacheBypassText = ''
        data.proxyNoCacheText = ''
      }
    }

    // 添加代理请求头
    const addProxyHeader = () => {
      data.proxySetHeaders.push({
        key: '',
        value: ''
      })
    }

    // 删除代理请求头
    const removeProxyHeader = (index) => {
      data.proxySetHeaders.splice(index, 1)
    }

    // 点击提交保存当前数据
    const onSubmit = async () => {
      let formData = JSON.parse(JSON.stringify(data.formData))
      formData.enable = formData.enable == true ? 1 : 2
      formData.request_methods = formData.request_methods.join(',')
      {
        let pri = formData.priority
        if (pri === null || pri === undefined || pri === '') {
          pri = 0
        } else {
          pri = Number(pri)
          if (!Number.isFinite(pri) || pri < 0) pri = 0
        }
        formData.priority = pri
      }

      // 处理新字段
      if (formData.client_max_body_size === undefined || formData.client_max_body_size === null || formData.client_max_body_size === '' || formData.client_max_body_size === '0') {
        delete formData.client_max_body_size
      } else {
        formData.client_max_body_size = String(formData.client_max_body_size).trim()
      }
      
      if (formData.chunked_transfer_encoding === undefined || formData.chunked_transfer_encoding === null) {
        delete formData.chunked_transfer_encoding
      } else {
        formData.chunked_transfer_encoding = formData.chunked_transfer_encoding
      }
      
      if (formData.proxy_buffering === undefined || formData.proxy_buffering === null) {
        delete formData.proxy_buffering
      } else {
        formData.proxy_buffering = formData.proxy_buffering
      }

      // 处理代理缓存配置
      if (data.proxyCacheEnabled) {
        const cacheBypass = data.proxyCacheBypassText
          ? data.proxyCacheBypassText.split(',').map(s => s.trim()).filter(s => s)
          : []
        const noCache = data.proxyNoCacheText
          ? data.proxyNoCacheText.split(',').map(s => s.trim()).filter(s => s)
          : []
        
        formData.proxy_cache = {
          enabled: true,
          cache_key: data.proxyCacheConfig.cache_key || undefined,
          cache_valid: data.proxyCacheConfig.cache_valid || undefined,
          cache_bypass: cacheBypass.length > 0 ? cacheBypass : undefined,
          no_cache: noCache.length > 0 ? noCache : undefined
        }
        // 移除空值
        Object.keys(formData.proxy_cache).forEach(key => {
          if (formData.proxy_cache[key] === undefined) {
            delete formData.proxy_cache[key]
          }
        })
      } else {
        delete formData.proxy_cache
      }

      // 处理代理请求头配置
      if (data.proxySetHeaders.length > 0) {
        const proxySetHeader = {}
        data.proxySetHeaders.forEach(header => {
          if (header.key && header.value) {
            proxySetHeader[header.key] = header.value
          }
        })
        if (Object.keys(proxySetHeader).length > 0) {
          formData.proxy_set_header = proxySetHeader
        } else {
          delete formData.proxy_set_header
        }
      } else {
        delete formData.proxy_set_header
      }

      const rrPayload = buildRewriteRulesPayload()
      if (rrPayload === null) {
        return
      }
      const rrKeys = Object.keys(rrPayload)
      if (rrKeys.length === 0) {
        if (props.currentResId != null) {
          formData.rewrite_rules = {}
        } else {
          delete formData.rewrite_rules
        }
      } else {
        formData.rewrite_rules = rrPayload
      }

      // 调用增加/修改接口
      let routerRes
      if (props.currentResId != null) {
        routerRes = await $routerUpdate(props.currentResId, formData)
      } else {
        routerRes = await $routerAdd(formData)
      }

      if (routerRes.code != 0) {
        message.error(routerRes.msg)
        return
      } else {
        message.success(routerRes.msg)
        emit('componentRefreshList')
      }
    }

    // 表单取消，关闭抽屉
    const cancel = async () => {
      emit('componentCloseDrawer')
    }

    // 定义函数
    const fn = reactive({
      onSubmit,
      cancel,
      filterOption,
      onProxyCacheChange,
      addProxyHeader,
      removeProxyHeader,
      addRewriteProxyHeader,
      removeRewriteProxyHeader,
      addRewriteStaticRow,
      removeRewriteStaticRow,
      addRewriteRegexRow,
      removeRewriteRegexRow
    })

    return {
      schemaRouter,
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
.form {
  width: '150px';
}
.rewrite-nested {
  margin-top: 16px;
}
.rewrite-inner-item {
  margin-bottom: 16px;
}
.rewrite-extra-tip {
  font-size: 12px;
  color: #888;
  line-height: 1.5;
}
</style>
