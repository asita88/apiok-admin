package validators

import (
	"apiok-admin/app/utils"
)

type ServiceAddUpdate struct {
	Name                    string                 `json:"name" zh:"服务名称" en:"Service name" binding:"omitempty,min=1,max=30"`
	Enable                  int                    `json:"enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Release                 int                    `json:"release" zh:"发布开关" en:"Release status enable" binding:"omitempty,oneof=1 2"`
	Protocol                int                    `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceDomains          []string               `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,min=1,CheckServiceDomain"`
	ClientMaxBodySize       *int64                 `json:"client_max_body_size" zh:"请求体大小限制" en:"Maximum request body size" binding:"omitempty"`
	ChunkedTransferEncoding *bool                  `json:"chunked_transfer_encoding" zh:"分块传输编码" en:"Chunked transfer encoding" binding:"omitempty"`
	ProxyBuffering          *bool                  `json:"proxy_buffering" zh:"代理缓冲" en:"Proxy buffering" binding:"omitempty"`
	ProxyCache              map[string]interface{} `json:"proxy_cache" zh:"代理缓存配置" en:"Proxy cache configuration" binding:"omitempty"`
	ProxySetHeader          map[string]string      `json:"proxy_set_header" zh:"代理请求头设置" en:"Proxy set header configuration" binding:"omitempty"`
}

type ServiceList struct {
	Protocol int    `form:"protocol" json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	Enable   int    `form:"enable" json:"enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Release  int    `form:"release" json:"release" zh:"发布状态" en:"Release status" binding:"omitempty,oneof=1 2 3"`
	Search   string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type ServiceUpdateName struct {
	Name string `json:"name" zh:"服务名称" en:"Service name" binding:"required,min=1,max=30"`
}

type ServiceSwitchEnable struct {
	Enable int `json:"enable" zh:"服务开关" en:"Service enable" binding:"required,oneof=1 2"`
}

func CorrectServiceAttributesDefault(serviceAddUpdate *ServiceAddUpdate) {
	if serviceAddUpdate.Protocol == 0 {
		serviceAddUpdate.Protocol = utils.ProtocolHTTP
	}
	if serviceAddUpdate.Enable == 0 {
		serviceAddUpdate.Enable = utils.EnableOff
	}
}
