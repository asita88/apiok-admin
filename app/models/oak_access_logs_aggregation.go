package models

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/validators"
	"errors"
	"fmt"
)

type AccessLogAggregation struct {
	TotalRequests   int64                 `json:"total_requests"`
	AvgResponseTime float64               `json:"avg_response_time"`
	MaxResponseTime float64               `json:"max_response_time"`
	MinResponseTime float64               `json:"min_response_time"`
	TotalBytesSent  int64                 `json:"total_bytes_sent"`
	ErrorCount      int64                 `json:"error_count"`
	ErrorRate       float64               `json:"error_rate"`
	StatusStats     []StatusStat          `json:"status_stats"`
	MethodStats     []MethodStat          `json:"method_stats"`
	ServiceStats    []ServiceStat         `json:"service_stats"`
	HostStats       []HostStat            `json:"host_stats"`
	PathStats       []PathStat            `json:"path_stats"`
	PathBytesStats  []PathBytesStat       `json:"path_bytes_stats"`
	TimeSeries      []TimeSeriesItem      `json:"time_series"`
	BytesTimeSeries []BytesTimeSeriesItem `json:"bytes_time_series"`
}

type StatusStat struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type MethodStat struct {
	Method string `json:"method"`
	Count  int64  `json:"count"`
}

type ServiceStat struct {
	ServiceName string `json:"service_name"`
	Count       int64  `json:"count"`
}

type HostStat struct {
	Host  string `json:"host"`
	Count int64  `json:"count"`
}

type PathStat struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

type PathBytesStat struct {
	Path  string `json:"path"`
	Bytes int64  `json:"bytes"`
}

type TimeSeriesItem struct {
	Time  int64 `json:"time"`
	Count int64 `json:"count"`
}

type BytesTimeSeriesItem struct {
	Time  int64 `json:"time"`
	Bytes int64 `json:"bytes"`
}

func AggregateAccessLog(param *validators.AccessLogList) (agg AccessLogAggregation, err error) {
	db := packages.GetDb()
	tables := resolveAccessLogTables(db, param)
	if len(tables) == 0 {
		return agg, nil
	}
	tx := accessLogFromTables(db, tables, param)

	var totalRequests int64
	tx.Count(&totalRequests)
	agg.TotalRequests = totalRequests

	if totalRequests == 0 {
		return agg, nil
	}

	type StatsResult struct {
		AvgTime    float64
		MaxTime    float64
		MinTime    float64
		TotalBytes int64
	}
	var stats StatsResult
	baseTx := accessLogFromTables(db, tables, param)
	baseTx.Select("AVG(request_time) as avg_time, MAX(request_time) as max_time, MIN(request_time) as min_time, SUM(bytes_sent) as total_bytes").Scan(&stats)
	agg.AvgResponseTime = stats.AvgTime
	agg.MaxResponseTime = stats.MaxTime
	agg.MinResponseTime = stats.MinTime
	agg.TotalBytesSent = stats.TotalBytes

	var errorCount int64
	accessLogFromTables(db, tables, param).Where("`status` >= ?", 400).Count(&errorCount)
	agg.ErrorCount = errorCount
	if totalRequests > 0 {
		agg.ErrorRate = float64(errorCount) / float64(totalRequests) * 100
	}

	type StatusRangeResult struct {
		StatusRange int   `gorm:"column:status_range"`
		Count       int64 `gorm:"column:count"`
	}
	var statusRangeResults []StatusRangeResult
	statusTx := accessLogFromTables(db, tables, param)
	statusTx.Select("FLOOR(`status` / 100) * 100 as status_range, COUNT(*) as count").
		Group("FLOOR(`status` / 100) * 100").
		Order("status_range ASC").
		Scan(&statusRangeResults)

	statusStats := make([]StatusStat, 0)
	for _, result := range statusRangeResults {
		div := result.StatusRange / 100
		if div <= 0 {
			div = 1
		}
		statusStats = append(statusStats, StatusStat{
			Status: fmt.Sprintf("%dxx", div),
			Count:  result.Count,
		})
	}
	agg.StatusStats = statusStats

	var methodStats []MethodStat
	methodTx := accessLogFromTables(db, tables, param)
	methodTx.Select("`method` as method, COUNT(*) as count").
		Group("`method`").
		Order("count DESC").
		Scan(&methodStats)
	agg.MethodStats = methodStats

	var serviceStats []ServiceStat
	serviceTx := accessLogFromTables(db, tables, param)
	serviceTx.Select("server_name as service_name, COUNT(*) as count").
		Where("server_name != '' AND server_name IS NOT NULL").
		Group("server_name").
		Order("count DESC").
		Limit(10).
		Scan(&serviceStats)
	agg.ServiceStats = serviceStats

	var hostStats []HostStat
	hostTx := accessLogFromTables(db, tables, param)
	hostTx.Select("server_name as host, COUNT(*) as count").
		Where("server_name IS NOT NULL AND server_name != ''").
		Group("server_name").
		Order("count DESC").
		Limit(5).
		Scan(&hostStats)
	agg.HostStats = hostStats

	pathExpr := "SUBSTRING_INDEX(url, '?', 1)"
	var pathStats []PathStat
	pathTx := accessLogFromTables(db, tables, param)
	pathTx.Select(pathExpr + " as path, COUNT(*) as count").
		Where("url IS NOT NULL AND url != ''").
		Group(pathExpr).
		Order("count DESC").
		Limit(5).
		Scan(&pathStats)
	agg.PathStats = pathStats

	var pathBytesStats []PathBytesStat
	pathBytesTx := accessLogFromTables(db, tables, param)
	pathBytesTx.Select(pathExpr + " as path, COALESCE(SUM(bytes_sent), 0) as bytes").
		Where("url IS NOT NULL AND url != ''").
		Group(pathExpr).
		Order("bytes DESC").
		Limit(5).
		Scan(&pathBytesStats)
	agg.PathBytesStats = pathBytesStats

	if param.StartTime > 0 && param.EndTime > 0 {
		interval := (param.EndTime - param.StartTime) / 20
		if interval < 60 {
			interval = 60
		}

		var timeSeries []TimeSeriesItem
		var bytesTimeSeries []BytesTimeSeriesItem
		timeCast := "CAST(NULLIF(TRIM(`time`),'') AS UNSIGNED)"
		for t := param.StartTime; t <= param.EndTime; t += interval {
			var count int64
			timeTx := accessLogFromTables(db, tables, param)
			timeTx.Where(timeCast+" >= ? AND "+timeCast+" < ?", t, t+interval).Count(&count)
			timeSeries = append(timeSeries, TimeSeriesItem{
				Time:  t,
				Count: count,
			})

			var bsum int64
			bytesTx := accessLogFromTables(db, tables, param)
			bytesTx.Where(timeCast+" >= ? AND "+timeCast+" < ?", t, t+interval).
				Select("COALESCE(SUM(bytes_sent), 0)").Row().Scan(&bsum)
			bytesTimeSeries = append(bytesTimeSeries, BytesTimeSeriesItem{
				Time:  t,
				Bytes: bsum,
			})
		}
		agg.TimeSeries = timeSeries
		agg.BytesTimeSeries = bytesTimeSeries
	}

	return agg, nil
}

type FieldAggregationResult struct {
	FieldName string      `json:"field_name"`
	Type      string      `json:"type"`
	Results   interface{} `json:"results"`
}

type FieldAggregationItem struct {
	Value interface{} `json:"value"`
	Count int64       `json:"count"`
}

var allowedFields = map[string]bool{
	"id":                     true,
	"time":                   true,
	"remote_addr":            true,
	"x_forwarded_for":        true,
	"url":                    true,
	"method":                 true,
	"request":                true,
	"request_length":         true,
	"referer":                true,
	"user_agent":             true,
	"connection_requests":    true,
	"upstream_cache_status":  true,
	"status":                 true,
	"request_time":           true,
	"upstream_response_time": true,
	"bytes_sent":             true,
	"body_bytes_sent":        true,
	"server_name":            true,
	"upstream_addr":          true,
	"upstream_status":        true,
	"request_id":             true,
	"block_reason":           true,
	"block_rule":             true,
	"geo_country":            true,
	"geo_province":           true,
	"geo_city":               true,
	"geo_isp":                true,
	"geo_country_code":       true,
}

func quoteAccessLogColumn(name string) string {
	reserved := map[string]bool{
		"time": true, "request": true, "method": true, "status": true,
	}
	if reserved[name] {
		return "`" + name + "`"
	}
	return name
}

func AggregateAccessLogField(param *validators.FieldAggregation) (result FieldAggregationResult, err error) {
	if !allowedFields[param.FieldName] {
		return result, errors.New("不允许的字段名")
	}

	result.FieldName = param.FieldName
	result.Type = param.AggregationType

	listParam := param.AccessLogList
	db := packages.GetDb()
	tables := resolveAccessLogTables(db, &listParam)
	fieldSQL := quoteAccessLogColumn(param.FieldName)

	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	if len(tables) == 0 {
		switch param.AggregationType {
		case "avg":
			result.Results = map[string]interface{}{"avg": 0.0}
		case "max":
			result.Results = map[string]interface{}{"max": nil}
		case "min":
			result.Results = map[string]interface{}{"min": nil}
		case "sum":
			result.Results = map[string]interface{}{"sum": int64(0)}
		case "stats":
			result.Results = struct {
				Count int64
				Min   float64
				Max   float64
				Avg   float64
				Sum   int64
			}{}
		default:
			result.Results = []FieldAggregationItem{}
		}
		return result, nil
	}

	tx := accessLogFromTables(db, tables, &listParam)

	switch param.AggregationType {
	case "count", "terms":
		var items []FieldAggregationItem
		tx.Select(fieldSQL + " as value, COUNT(*) as count").
			Where(fieldSQL + " IS NOT NULL AND " + fieldSQL + " != ''").
			Group(fieldSQL).
			Order("count DESC").
			Limit(limit).
			Scan(&items)
		result.Results = items

	case "avg":
		var avgValue float64
		tx.Select("AVG(" + fieldSQL + ")").Row().Scan(&avgValue)
		result.Results = map[string]interface{}{
			"avg": avgValue,
		}

	case "max":
		var maxValue interface{}
		tx.Select("MAX(" + fieldSQL + ")").Row().Scan(&maxValue)
		result.Results = map[string]interface{}{
			"max": maxValue,
		}

	case "min":
		var minValue interface{}
		tx.Select("MIN(" + fieldSQL + ")").Row().Scan(&minValue)
		result.Results = map[string]interface{}{
			"min": minValue,
		}

	case "sum":
		var sumValue int64
		tx.Select("SUM(" + fieldSQL + ")").Row().Scan(&sumValue)
		result.Results = map[string]interface{}{
			"sum": sumValue,
		}

	case "stats":
		type StatsResult struct {
			Count int64
			Min   float64
			Max   float64
			Avg   float64
			Sum   int64
		}
		var stats StatsResult
		tx.Select("COUNT(*) as count, MIN(" + fieldSQL + ") as min, MAX(" + fieldSQL + ") as max, AVG(" + fieldSQL + ") as avg, SUM(" + fieldSQL + ") as sum").Scan(&stats)
		result.Results = stats

	default:
		var items []FieldAggregationItem
		tx.Select(fieldSQL + " as value, COUNT(*) as count").
			Where(fieldSQL + " IS NOT NULL AND " + fieldSQL + " != ''").
			Group(fieldSQL).
			Order("count DESC").
			Limit(limit).
			Scan(&items)
		result.Results = items
	}

	return result, nil
}
