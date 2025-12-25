package models

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/validators"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AccessLog struct {
	ID                    int64     `gorm:"column:id;primary_key"`
	Timestamp             int64     `gorm:"column:timestamp"`
	RequestMethod         string    `gorm:"column:request_method"`
	RequestURI            string    `gorm:"column:request_uri"`
	RequestPath           string    `gorm:"column:request_path"`
	RequestQueryString    string    `gorm:"column:request_query_string;type:text"`
	RequestProtocol       string    `gorm:"column:request_protocol"`
	RemoteAddr            string    `gorm:"column:remote_addr"`
	RemotePort            int       `gorm:"column:remote_port"`
	ServerAddr            string    `gorm:"column:server_addr"`
	ServerPort            int       `gorm:"column:server_port"`
	RequestHost           string    `gorm:"column:request_host"`
	RequestHeaders        string    `gorm:"column:request_headers;type:text"`
	RequestArgs           string    `gorm:"column:request_args;type:text"`
	RequestBody           string    `gorm:"column:request_body;type:longtext"`
	ResponseStatus        int       `gorm:"column:response_status"`
	ResponseHeaders       string    `gorm:"column:response_headers;type:text"`
	ResponseBody          string    `gorm:"column:response_body;type:longtext"`
	UpstreamResponseTime  string    `gorm:"column:upstream_response_time"`
	UpstreamConnectTime   string    `gorm:"column:upstream_connect_time"`
	RequestTime           float64   `gorm:"column:request_time"`
	BytesSent             int64     `gorm:"column:bytes_sent"`
	ServiceName           string    `gorm:"column:service_name"`
	RouterName            string    `gorm:"column:router_name"`
	CreatedAt             time.Time `gorm:"column:created_at"`
}

func (a *AccessLog) TableName() string {
	return "apiok_access_log"
}

func (a *AccessLog) AccessLogListPage(param *validators.AccessLogList) (list []AccessLog, total int, listError error) {
	tx := packages.GetDb().Table(a.TableName())

	if param.StartTime > 0 {
		tx = tx.Where("timestamp >= ?", param.StartTime)
	}

	if param.EndTime > 0 {
		tx = tx.Where("timestamp <= ?", param.EndTime)
	}

	if len(param.RequestMethod) > 0 {
		tx = tx.Where("request_method = ?", param.RequestMethod)
	}

	if len(param.RemoteAddr) > 0 {
		tx = tx.Where("remote_addr = ?", param.RemoteAddr)
	}

	if param.ResponseStatus > 0 {
		tx = tx.Where("response_status = ?", param.ResponseStatus)
	}

	if len(param.ServiceName) > 0 {
		tx = tx.Where("service_name = ?", param.ServiceName)
	}

	if len(param.RouterName) > 0 {
		tx = tx.Where("router_name = ?", param.RouterName)
	}

	if len(param.RequestHost) > 0 {
		tx = tx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("timestamp desc")

	listError = ListPaginate(tx, &list, &param.BaseListPage)

	return
}

type AccessLogAggregation struct {
	TotalRequests    int64   `json:"total_requests"`
	AvgResponseTime  float64 `json:"avg_response_time"`
	MaxResponseTime  float64 `json:"max_response_time"`
	MinResponseTime  float64 `json:"min_response_time"`
	TotalBytesSent   int64   `json:"total_bytes_sent"`
	ErrorCount       int64   `json:"error_count"`
	ErrorRate        float64 `json:"error_rate"`
	StatusStats      []StatusStat `json:"status_stats"`
	MethodStats      []MethodStat `json:"method_stats"`
	ServiceStats     []ServiceStat `json:"service_stats"`
	HostStats        []HostStat `json:"host_stats"`
	PathStats        []PathStat `json:"path_stats"`
	PathBytesStats   []PathBytesStat `json:"path_bytes_stats"`
	TimeSeries       []TimeSeriesItem `json:"time_series"`
	BytesTimeSeries  []BytesTimeSeriesItem `json:"bytes_time_series"`
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

func (a *AccessLog) AccessLogAggregation(param *validators.AccessLogList) (agg AccessLogAggregation, err error) {
	tx := packages.GetDb().Table(a.TableName())

	if param.StartTime > 0 {
		tx = tx.Where("timestamp >= ?", param.StartTime)
	}

	if param.EndTime > 0 {
		tx = tx.Where("timestamp <= ?", param.EndTime)
	}

	if len(param.RequestMethod) > 0 {
		tx = tx.Where("request_method = ?", param.RequestMethod)
	}

	if len(param.RemoteAddr) > 0 {
		tx = tx.Where("remote_addr = ?", param.RemoteAddr)
	}

	if param.ResponseStatus > 0 {
		tx = tx.Where("response_status = ?", param.ResponseStatus)
	}

	if len(param.ServiceName) > 0 {
		tx = tx.Where("service_name = ?", param.ServiceName)
	}

	if len(param.RouterName) > 0 {
		tx = tx.Where("router_name = ?", param.RouterName)
	}

	if len(param.RequestHost) > 0 {
		tx = tx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}

	var totalRequests int64
	tx.Count(&totalRequests)
	agg.TotalRequests = totalRequests

	if totalRequests == 0 {
		return agg, nil
	}

	type StatsResult struct {
		AvgTime float64
		MaxTime float64
		MinTime float64
		TotalBytes int64
	}
	var stats StatsResult
	baseTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		baseTx = baseTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		baseTx = baseTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		baseTx = baseTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		baseTx = baseTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		baseTx = baseTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		baseTx = baseTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		baseTx = baseTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		baseTx = baseTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		baseTx = baseTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	baseTx.Select("AVG(request_time) as avg_time, MAX(request_time) as max_time, MIN(request_time) as min_time, SUM(bytes_sent) as total_bytes").Scan(&stats)
	agg.AvgResponseTime = stats.AvgTime
	agg.MaxResponseTime = stats.MaxTime
	agg.MinResponseTime = stats.MinTime
	agg.TotalBytesSent = stats.TotalBytes

	var errorCount int64
	baseTx.Where("response_status >= ?", 400).Count(&errorCount)
	agg.ErrorCount = errorCount
	if totalRequests > 0 {
		agg.ErrorRate = float64(errorCount) / float64(totalRequests) * 100
	}

	var statusStats []StatusStat
	statusTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		statusTx = statusTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		statusTx = statusTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		statusTx = statusTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		statusTx = statusTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		statusTx = statusTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		statusTx = statusTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		statusTx = statusTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		statusTx = statusTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		statusTx = statusTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	type StatusRangeResult struct {
		StatusRange int   `gorm:"column:status_range"`
		Count       int64 `gorm:"column:count"`
	}
	var statusRangeResults []StatusRangeResult
	statusTx.Select("FLOOR(response_status / 100) * 100 as status_range, COUNT(*) as count").
		Group("status_range").
		Order("status_range ASC").
		Scan(&statusRangeResults)
	
	statusStats = make([]StatusStat, 0)
	for _, result := range statusRangeResults {
		statusStats = append(statusStats, StatusStat{
			Status: fmt.Sprintf("%dxx", result.StatusRange/100),
			Count:  result.Count,
		})
	}
	agg.StatusStats = statusStats

	var methodStats []MethodStat
	methodTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		methodTx = methodTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		methodTx = methodTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		methodTx = methodTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		methodTx = methodTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		methodTx = methodTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		methodTx = methodTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		methodTx = methodTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		methodTx = methodTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		methodTx = methodTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	methodTx.Select("request_method as method, COUNT(*) as count").
		Group("request_method").
		Order("count DESC").
		Scan(&methodStats)
	agg.MethodStats = methodStats

	var serviceStats []ServiceStat
	serviceTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		serviceTx = serviceTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		serviceTx = serviceTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		serviceTx = serviceTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		serviceTx = serviceTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		serviceTx = serviceTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		serviceTx = serviceTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		serviceTx = serviceTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		serviceTx = serviceTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		serviceTx = serviceTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	serviceTx.Select("service_name, COUNT(*) as count").
		Where("service_name != '' AND service_name IS NOT NULL").
		Group("service_name").
		Order("count DESC").
		Limit(10).
		Scan(&serviceStats)
	agg.ServiceStats = serviceStats

	var hostStats []HostStat
	hostTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		hostTx = hostTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		hostTx = hostTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		hostTx = hostTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		hostTx = hostTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		hostTx = hostTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		hostTx = hostTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		hostTx = hostTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		hostTx = hostTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		hostTx = hostTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	hostTx.Select("request_host as host, COUNT(*) as count").
		Where("request_host IS NOT NULL AND request_host != ''").
		Group("request_host").
		Order("count DESC").
		Limit(5).
		Scan(&hostStats)
	agg.HostStats = hostStats

	var pathStats []PathStat
	pathTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		pathTx = pathTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		pathTx = pathTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		pathTx = pathTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		pathTx = pathTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		pathTx = pathTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		pathTx = pathTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		pathTx = pathTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		pathTx = pathTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		pathTx = pathTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	pathTx.Select("request_path as path, COUNT(*) as count").
		Where("request_path IS NOT NULL AND request_path != ''").
		Group("request_path").
		Order("count DESC").
		Limit(5).
		Scan(&pathStats)
	agg.PathStats = pathStats

	var pathBytesStats []PathBytesStat
	pathBytesTx := packages.GetDb().Table(a.TableName())
	if param.StartTime > 0 {
		pathBytesTx = pathBytesTx.Where("timestamp >= ?", param.StartTime)
	}
	if param.EndTime > 0 {
		pathBytesTx = pathBytesTx.Where("timestamp <= ?", param.EndTime)
	}
	if len(param.RequestMethod) > 0 {
		pathBytesTx = pathBytesTx.Where("request_method = ?", param.RequestMethod)
	}
	if len(param.RemoteAddr) > 0 {
		pathBytesTx = pathBytesTx.Where("remote_addr = ?", param.RemoteAddr)
	}
	if param.ResponseStatus > 0 {
		pathBytesTx = pathBytesTx.Where("response_status = ?", param.ResponseStatus)
	}
	if len(param.ServiceName) > 0 {
		pathBytesTx = pathBytesTx.Where("service_name = ?", param.ServiceName)
	}
	if len(param.RouterName) > 0 {
		pathBytesTx = pathBytesTx.Where("router_name = ?", param.RouterName)
	}
	if len(param.RequestHost) > 0 {
		pathBytesTx = pathBytesTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}
	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		pathBytesTx = pathBytesTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}
	pathBytesTx.Select("request_path as path, COALESCE(SUM(bytes_sent), 0) as bytes").
		Where("request_path IS NOT NULL AND request_path != ''").
		Group("request_path").
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
		baseTx := packages.GetDb().Table(a.TableName())
		if param.StartTime > 0 {
			baseTx = baseTx.Where("timestamp >= ?", param.StartTime)
		}
		if param.EndTime > 0 {
			baseTx = baseTx.Where("timestamp <= ?", param.EndTime)
		}
		if len(param.RequestMethod) > 0 {
			baseTx = baseTx.Where("request_method = ?", param.RequestMethod)
		}
		if len(param.RemoteAddr) > 0 {
			baseTx = baseTx.Where("remote_addr = ?", param.RemoteAddr)
		}
		if param.ResponseStatus > 0 {
			baseTx = baseTx.Where("response_status = ?", param.ResponseStatus)
		}
		if len(param.ServiceName) > 0 {
			baseTx = baseTx.Where("service_name = ?", param.ServiceName)
		}
		if len(param.RouterName) > 0 {
			baseTx = baseTx.Where("router_name = ?", param.RouterName)
		}
		if len(param.RequestHost) > 0 {
			baseTx = baseTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
		}
		param.Search = strings.TrimSpace(param.Search)
		if len(param.Search) != 0 {
			search := "%" + param.Search + "%"
			baseTx = baseTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
		}

		var bytesTimeSeries []BytesTimeSeriesItem
		for t := param.StartTime; t <= param.EndTime; t += interval {
			var count int64
			var bytes int64
			timeTx := packages.GetDb().Table(a.TableName())
			if param.StartTime > 0 {
				timeTx = timeTx.Where("timestamp >= ?", param.StartTime)
			}
			if param.EndTime > 0 {
				timeTx = timeTx.Where("timestamp <= ?", param.EndTime)
			}
			if len(param.RequestMethod) > 0 {
				timeTx = timeTx.Where("request_method = ?", param.RequestMethod)
			}
			if len(param.RemoteAddr) > 0 {
				timeTx = timeTx.Where("remote_addr = ?", param.RemoteAddr)
			}
			if param.ResponseStatus > 0 {
				timeTx = timeTx.Where("response_status = ?", param.ResponseStatus)
			}
			if len(param.ServiceName) > 0 {
				timeTx = timeTx.Where("service_name = ?", param.ServiceName)
			}
			if len(param.RouterName) > 0 {
				timeTx = timeTx.Where("router_name = ?", param.RouterName)
			}
			if len(param.RequestHost) > 0 {
				timeTx = timeTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
			}
			param.Search = strings.TrimSpace(param.Search)
			if len(param.Search) != 0 {
				search := "%" + param.Search + "%"
				timeTx = timeTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
			}
			timeTx.Where("timestamp >= ? AND timestamp < ?", t, t+interval).Count(&count)
			timeSeries = append(timeSeries, TimeSeriesItem{
				Time:  t,
				Count: count,
			})

			bytesTx := packages.GetDb().Table(a.TableName())
			if param.StartTime > 0 {
				bytesTx = bytesTx.Where("timestamp >= ?", param.StartTime)
			}
			if param.EndTime > 0 {
				bytesTx = bytesTx.Where("timestamp <= ?", param.EndTime)
			}
			if len(param.RequestMethod) > 0 {
				bytesTx = bytesTx.Where("request_method = ?", param.RequestMethod)
			}
			if len(param.RemoteAddr) > 0 {
				bytesTx = bytesTx.Where("remote_addr = ?", param.RemoteAddr)
			}
			if param.ResponseStatus > 0 {
				bytesTx = bytesTx.Where("response_status = ?", param.ResponseStatus)
			}
			if len(param.ServiceName) > 0 {
				bytesTx = bytesTx.Where("service_name = ?", param.ServiceName)
			}
			if len(param.RouterName) > 0 {
				bytesTx = bytesTx.Where("router_name = ?", param.RouterName)
			}
			if len(param.RequestHost) > 0 {
				bytesTx = bytesTx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
			}
			param.Search = strings.TrimSpace(param.Search)
			if len(param.Search) != 0 {
				search := "%" + param.Search + "%"
				bytesTx = bytesTx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
			}
			bytesTx.Where("timestamp >= ? AND timestamp < ?", t, t+interval).Select("COALESCE(SUM(bytes_sent), 0)").Row().Scan(&bytes)
			bytesTimeSeries = append(bytesTimeSeries, BytesTimeSeriesItem{
				Time:  t,
				Bytes: bytes,
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
	"timestamp":             true,
	"request_method":        true,
	"request_uri":           true,
	"request_path":          true,
	"request_query_string":  true,
	"request_protocol":      true,
	"remote_addr":           true,
	"remote_port":           true,
	"server_addr":           true,
	"server_port":           true,
	"request_host":          true,
	"request_headers":       true,
	"request_args":          true,
	"request_body":          true,
	"response_status":       true,
	"response_headers":      true,
	"response_body":          true,
	"upstream_response_time": true,
	"upstream_connect_time": true,
	"request_time":          true,
	"bytes_sent":            true,
	"service_name":          true,
	"router_name":           true,
}

func (a *AccessLog) FieldAggregation(param *validators.FieldAggregation) (result FieldAggregationResult, err error) {
	if !allowedFields[param.FieldName] {
		return result, errors.New("不允许的字段名")
	}

	result.FieldName = param.FieldName
	result.Type = param.AggregationType

	tx := packages.GetDb().Table(a.TableName())

	if param.StartTime > 0 {
		tx = tx.Where("timestamp >= ?", param.StartTime)
	}

	if param.EndTime > 0 {
		tx = tx.Where("timestamp <= ?", param.EndTime)
	}

	if len(param.RequestMethod) > 0 {
		tx = tx.Where("request_method = ?", param.RequestMethod)
	}

	if len(param.RemoteAddr) > 0 {
		tx = tx.Where("remote_addr = ?", param.RemoteAddr)
	}

	if param.ResponseStatus > 0 {
		tx = tx.Where("response_status = ?", param.ResponseStatus)
	}

	if len(param.ServiceName) > 0 {
		tx = tx.Where("service_name = ?", param.ServiceName)
	}

	if len(param.RouterName) > 0 {
		tx = tx.Where("router_name = ?", param.RouterName)
	}

	if len(param.RequestHost) > 0 {
		tx = tx.Where("request_host LIKE ?", "%"+param.RequestHost+"%")
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("request_uri LIKE ? OR request_path LIKE ? OR request_host LIKE ?", search, search, search)
	}

	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	switch param.AggregationType {
	case "count", "terms":
		var items []FieldAggregationItem
		fieldName := param.FieldName
		tx.Select(fieldName+" as value, COUNT(*) as count").
			Where(fieldName+" IS NOT NULL AND "+fieldName+" != ''").
			Group(fieldName).
			Order("count DESC").
			Limit(limit).
			Scan(&items)
		result.Results = items

	case "avg":
		var avgValue float64
		tx.Select("AVG(" + param.FieldName + ")").Row().Scan(&avgValue)
		result.Results = map[string]interface{}{
			"avg": avgValue,
		}

	case "max":
		var maxValue interface{}
		tx.Select("MAX(" + param.FieldName + ")").Row().Scan(&maxValue)
		result.Results = map[string]interface{}{
			"max": maxValue,
		}

	case "min":
		var minValue interface{}
		tx.Select("MIN(" + param.FieldName + ")").Row().Scan(&minValue)
		result.Results = map[string]interface{}{
			"min": minValue,
		}

	case "sum":
		var sumValue int64
		tx.Select("SUM(" + param.FieldName + ")").Row().Scan(&sumValue)
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
		tx.Select("COUNT(*) as count, MIN("+param.FieldName+") as min, MAX("+param.FieldName+") as max, AVG("+param.FieldName+") as avg, SUM("+param.FieldName+") as sum").Scan(&stats)
		result.Results = stats

	default:
		var items []FieldAggregationItem
		fieldName := param.FieldName
		tx.Select(fieldName+" as value, COUNT(*) as count").
			Where(fieldName+" IS NOT NULL AND "+fieldName+" != ''").
			Group(fieldName).
			Order("count DESC").
			Limit(limit).
			Scan(&items)
		result.Results = items
	}

	return result, nil
}

