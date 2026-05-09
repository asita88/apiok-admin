package validators

type AccessLogList struct {
	StartTime int64  `form:"start_time" json:"start_time" zh:"开始时间戳" en:"Start timestamp" binding:"omitempty"`
	EndTime   int64  `form:"end_time" json:"end_time" zh:"结束时间戳" en:"End timestamp" binding:"omitempty"`
	Query     string `form:"query" json:"query" zh:"查询条件" en:"Query" binding:"omitempty"`
	BaseListPage
}
