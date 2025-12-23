package validators

type LogList struct {
	Username     string `form:"username" json:"username" zh:"用户名" en:"Username" binding:"omitempty"`
	Action       string `form:"action" json:"action" zh:"操作类型" en:"Action type" binding:"omitempty"`
	ResourceType string `form:"resource_type" json:"resource_type" zh:"资源类型" en:"Resource type" binding:"omitempty"`
	StatusCode   int    `form:"status_code" json:"status_code" zh:"状态码" en:"Status code" binding:"omitempty"`
	Search       string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

