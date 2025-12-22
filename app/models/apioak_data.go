package models

import (
	"apiok-admin/app/packages"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ApiokData struct {
	ID   uint64 `gorm:"column:id;primary_key"`
	Type string `gorm:"column:type"`
	Name string `gorm:"column:name"`
	Data string `gorm:"column:data;type:text"`
	ModelTime
}

type ApiokSyncHash struct {
	ID        uint64    `gorm:"column:id;primary_key"`
	HashKey   string    `gorm:"column:hash_key"`
	HashValue string    `gorm:"column:hash_value;type:text"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (a ApiokData) TableName() string {
	return "apiok_data"
}

func (a ApiokSyncHash) TableName() string {
	return "apiok_sync_hash"
}

func updateSyncHash(tx *gorm.DB) error {
	var syncHash ApiokSyncHash
	result := tx.Where("hash_key = ?", "sync/update").First(&syncHash)

	hashValue := map[string]string{
		"old": "",
		"new": "",
	}

	if result.Error == nil {
		if err := json.Unmarshal([]byte(syncHash.HashValue), &hashValue); err != nil {
			hashValue["old"] = ""
		}
	}

	allData := []ApiokData{}
	if err := tx.Find(&allData).Error; err != nil {
		return err
	}

	dataBytes, _ := json.Marshal(allData)
	newHash := fmt.Sprintf("%x", md5.Sum(dataBytes))
	hashValue["old"] = hashValue["new"]
	hashValue["new"] = newHash

	hashValueJSON, _ := json.Marshal(hashValue)

	if result.Error == nil {
		return tx.Model(&syncHash).Updates(map[string]interface{}{
			"hash_value": string(hashValueJSON),
		}).Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tx.Create(&ApiokSyncHash{
			HashKey:   "sync/update",
			HashValue: string(hashValueJSON),
		}).Error
	}
	return result.Error
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
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"data": string(dataJSON),
			}).Error; err != nil {
				return err
			}
		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&ApiokData{
				Type: dataType,
				Name: name,
				Data: string(dataJSON),
			}).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}

		return updateSyncHash(tx)
	})
}

func (a *ApiokData) Delete(dataType string, name string) error {
	return packages.GetDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type = ? AND name = ?", dataType, name).Delete(&ApiokData{}).Error; err != nil {
			return err
		}
		return updateSyncHash(tx)
	})
}

func (a *ApiokData) DeleteByType(dataType string) error {
	return packages.GetDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type = ?", dataType).Delete(&ApiokData{}).Error; err != nil {
			return err
		}
		return updateSyncHash(tx)
	})
}
