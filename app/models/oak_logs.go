package models

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"errors"
	"strings"
)

type Logs struct {
	ID           int    `gorm:"column:id;primary_key"`
	ResID        string `gorm:"column:res_id"`
	Username     string `gorm:"column:username"`
	Action       string `gorm:"column:action"`
	ResourceType string `gorm:"column:resource_type"`
	ResourceID   string `gorm:"column:resource_id"`
	Method       string `gorm:"column:method"`
	Path         string `gorm:"column:path"`
	IP           string `gorm:"column:ip"`
	RequestData  string `gorm:"column:request_data;type:text"`
	ResponseData string `gorm:"column:response_data;type:text"`
	StatusCode   int    `gorm:"column:status_code"`
	ErrorMessage string `gorm:"column:error_message"`
	ModelTime
}

func (l *Logs) TableName() string {
	return "ok_logs"
}

var recursionTimesLogs = 1

func (m *Logs) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeLog)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesLogs = 1
		return generateId, nil
	} else {
		if recursionTimesLogs == utils.IdGenerateMaxTimes {
			recursionTimesLogs = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesLogs++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (l *Logs) LogAdd(logData *Logs) error {
	logId, err := l.ModelUniqueId()
	if err != nil {
		return err
	}
	logData.ResID = logId

	err = packages.GetDb().
		Table(l.TableName()).
		Create(logData).Error

	return err
}

func (l *Logs) LogListPage(param *validators.LogList) (list []Logs, total int, listError error) {
	logsModel := Logs{}
	tx := packages.GetDb().
		Table(logsModel.TableName())

	if len(param.Username) > 0 {
		tx = tx.Where("username = ?", param.Username)
	}

	if len(param.Action) > 0 {
		tx = tx.Where("action = ?", param.Action)
	}

	if len(param.ResourceType) > 0 {
		tx = tx.Where("resource_type = ?", param.ResourceType)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("path LIKE ? OR resource_id LIKE ? OR ip LIKE ?", search, search, search)
	}

	if param.StatusCode > 0 {
		tx = tx.Where("status_code = ?", param.StatusCode)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("created_at desc")

	listError = ListPaginate(tx, &list, &param.BaseListPage)

	if len(list) == 0 {
		return
	}

	return
}

