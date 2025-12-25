package validators

type AccessLogList struct {
	StartTime     int64  `form:"start_time" json:"start_time" zh:"开始时间戳" en:"Start timestamp" binding:"omitempty"`
	EndTime       int64  `form:"end_time" json:"end_time" zh:"结束时间戳" en:"End timestamp" binding:"omitempty"`
	RequestMethod string `form:"request_method" json:"request_method" zh:"请求方法" en:"Request method" binding:"omitempty"`
	RemoteAddr    string `form:"remote_addr" json:"remote_addr" zh:"客户端IP" en:"Remote address" binding:"omitempty"`
	ResponseStatus int   `form:"response_status" json:"response_status" zh:"响应状态码" en:"Response status" binding:"omitempty"`
	ServiceName   string `form:"service_name" json:"service_name" zh:"服务名称" en:"Service name" binding:"omitempty"`
	RouterName    string `form:"router_name" json:"router_name" zh:"路由名称" en:"Router name" binding:"omitempty"`
	RequestHost   string `form:"request_host" json:"request_host" zh:"请求Host" en:"Request host" binding:"omitempty"`
	Search        string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

