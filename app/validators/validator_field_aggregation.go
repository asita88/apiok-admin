package validators

type FieldAggregation struct {
	FieldName    string `form:"field_name" json:"field_name" zh:"字段名" en:"Field name" binding:"required"`
	AggregationType string `form:"aggregation_type" json:"aggregation_type" zh:"聚合类型" en:"Aggregation type" binding:"required"`
	Limit        int    `form:"limit" json:"limit" zh:"限制数量" en:"Limit" binding:"omitempty"`
	AccessLogList
}

