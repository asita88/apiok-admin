package models

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"errors"
	"time"
)

type AcmeChallenges struct {
	ID           int       `gorm:"column:id;primary_key"`
	Token        string    `gorm:"column:token;uniqueIndex"`
	KeyAuth      string    `gorm:"column:key_auth"`
	ExpiredAt    time.Time `gorm:"column:expired_at"`
	ModelTime
}

func (a *AcmeChallenges) TableName() string {
	return "ok_acme_challenges"
}

var recursionTimesAcmeChallenges = 1

func (m *AcmeChallenges) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeAcmeChallenge)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("id = ?", generateId).
		Select("id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesAcmeChallenges = 1
		return generateId, nil
	} else {
		if recursionTimesAcmeChallenges == utils.IdGenerateMaxTimes {
			recursionTimesAcmeChallenges = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesAcmeChallenges++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (a *AcmeChallenges) ChallengeAdd(token, keyAuth string, expiredAt time.Time) error {
	challenge := &AcmeChallenges{
		Token:     token,
		KeyAuth:   keyAuth,
		ExpiredAt: expiredAt,
	}
	return packages.GetDb().Table(a.TableName()).Create(challenge).Error
}

func (a *AcmeChallenges) ChallengeGet(token string) (string, error) {
	var challenge AcmeChallenges
	err := packages.GetDb().Table(a.TableName()).
		Where("token = ? AND expired_at > ?", token, time.Now()).
		First(&challenge).Error
	if err != nil {
		return "", err
	}
	return challenge.KeyAuth, nil
}

func (a *AcmeChallenges) ChallengeDelete(token string) error {
	return packages.GetDb().Table(a.TableName()).
		Where("token = ?", token).
		Delete(&AcmeChallenges{}).Error
}

func (a *AcmeChallenges) ChallengeDeleteExpired() error {
	return packages.GetDb().Table(a.TableName()).
		Where("expired_at <= ?", time.Now()).
		Delete(&AcmeChallenges{}).Error
}

