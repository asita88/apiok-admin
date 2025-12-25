package services

import (
	"apiok-admin/app/models"
	"apiok-admin/app/validators"
)

type LogItem struct {
	ID           int    `json:"id"`
	ResID        string `json:"res_id"`
	Username     string `json:"username"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	IP           string `json:"ip"`
	RequestData  string `json:"request_data"`
	ResponseData string `json:"response_data"`
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
}

func LogList(request *validators.LogList) ([]LogItem, int, error) {
	logModel := models.Logs{}
	list, total, err := logModel.LogListPage(request)
	if err != nil {
		return []LogItem{}, 0, err
	}

	logList := make([]LogItem, 0)
	for _, v := range list {
		logItem := LogItem{
			ID:           v.ID,
			ResID:        v.ResID,
			Username:     v.Username,
			Action:       v.Action,
			ResourceType: v.ResourceType,
			ResourceID:   v.ResourceID,
			Method:       v.Method,
			Path:         v.Path,
			IP:           v.IP,
			RequestData:  v.RequestData,
			ResponseData: v.ResponseData,
			StatusCode:   v.StatusCode,
			ErrorMessage: v.ErrorMessage,
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		logList = append(logList, logItem)
	}

	return logList, total, nil
}

func LogAdd(logData *models.Logs) error {
	logModel := models.Logs{}
	return logModel.LogAdd(logData)
}

type AccessLogItem struct {
	ID                   int64   `json:"id"`
	Timestamp            int64   `json:"timestamp"`
	RequestMethod        string  `json:"request_method"`
	RequestURI           string  `json:"request_uri"`
	RequestPath          string  `json:"request_path"`
	RequestQueryString   string  `json:"request_query_string"`
	RequestProtocol      string  `json:"request_protocol"`
	RemoteAddr           string  `json:"remote_addr"`
	RemotePort           int     `json:"remote_port"`
	ServerAddr           string  `json:"server_addr"`
	ServerPort           int     `json:"server_port"`
	RequestHost          string  `json:"request_host"`
	RequestHeaders       string  `json:"request_headers"`
	RequestArgs          string  `json:"request_args"`
	RequestBody          string  `json:"request_body"`
	ResponseStatus       int     `json:"response_status"`
	ResponseHeaders      string  `json:"response_headers"`
	ResponseBody         string  `json:"response_body"`
	UpstreamResponseTime string  `json:"upstream_response_time"`
	UpstreamConnectTime  string  `json:"upstream_connect_time"`
	RequestTime          float64 `json:"request_time"`
	BytesSent            int64   `json:"bytes_sent"`
	ServiceName          string  `json:"service_name"`
	RouterName           string  `json:"router_name"`
	CreatedAt            string  `json:"created_at"`
}

func AccessLogList(request *validators.AccessLogList) ([]AccessLogItem, int, error) {
	accessLogModel := models.AccessLog{}
	list, total, err := accessLogModel.AccessLogListPage(request)
	if err != nil {
		return []AccessLogItem{}, 0, err
	}

	logList := make([]AccessLogItem, 0)
	for _, v := range list {
		logItem := AccessLogItem{
			ID:                   v.ID,
			Timestamp:            v.Timestamp,
			RequestMethod:        v.RequestMethod,
			RequestURI:           v.RequestURI,
			RequestPath:          v.RequestPath,
			RequestQueryString:   v.RequestQueryString,
			RequestProtocol:      v.RequestProtocol,
			RemoteAddr:           v.RemoteAddr,
			RemotePort:           v.RemotePort,
			ServerAddr:           v.ServerAddr,
			ServerPort:           v.ServerPort,
			RequestHost:          v.RequestHost,
			RequestHeaders:       v.RequestHeaders,
			RequestArgs:          v.RequestArgs,
			RequestBody:          v.RequestBody,
			ResponseStatus:       v.ResponseStatus,
			ResponseHeaders:      v.ResponseHeaders,
			ResponseBody:         v.ResponseBody,
			UpstreamResponseTime: v.UpstreamResponseTime,
			UpstreamConnectTime:  v.UpstreamConnectTime,
			RequestTime:          v.RequestTime,
			BytesSent:            v.BytesSent,
			ServiceName:          v.ServiceName,
			RouterName:           v.RouterName,
			CreatedAt:            v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		logList = append(logList, logItem)
	}

	return logList, total, nil
}

func AccessLogAggregation(request *validators.AccessLogList) (models.AccessLogAggregation, error) {
	accessLogModel := models.AccessLog{}
	return accessLogModel.AccessLogAggregation(request)
}

func FieldAggregation(request *validators.FieldAggregation) (models.FieldAggregationResult, error) {
	accessLogModel := models.AccessLog{}
	return accessLogModel.FieldAggregation(request)
}

