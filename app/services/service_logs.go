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

