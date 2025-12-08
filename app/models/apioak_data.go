package models

import (
	"apiok-admin/app/packages"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type ApiokData struct {
	ID   uint64 `gorm:"column:id;primary_key"`
	Type string `gorm:"column:type"`
	Name string `gorm:"column:name"`
	Data string `gorm:"column:data;type:text"`
	ModelTime
}

func (a ApiokData) TableName() string {
	return "apiok_data"
}

func (a *ApiokData) Upsert(dataType string, name string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return packages.GetDb().Transaction(func(tx *gorm.DB) error {
		var existing ApiokData
		result := tx.Where("type = ? AND name = ?", dataType, name).First(&existing)

		if result.Error == nil {
			return tx.Model(&existing).Updates(map[string]interface{}{
				"data": string(dataJSON),
			}).Error
		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return tx.Create(&ApiokData{
				Type: dataType,
				Name: name,
				Data: string(dataJSON),
			}).Error
		}
		return result.Error
	})
}

func (a *ApiokData) Delete(dataType string, name string) error {
	return packages.GetDb().Where("type = ? AND name = ?", dataType, name).Delete(&ApiokData{}).Error
}

func (a *ApiokData) DeleteByType(dataType string) error {
	return packages.GetDb().Where("type = ?", dataType).Delete(&ApiokData{}).Error
}
