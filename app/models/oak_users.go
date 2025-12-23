package models

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Users struct {
	ID       int    `gorm:"column:id;primary_key"` // primary key
	ResID    string `gorm:"column:res_id"`         // User iD
	Name     string `gorm:"column:name"`           // User name
	Password string `gorm:"column:password"`       // Password
	Email    string `gorm:"column:email"`          // Email
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "ok_users"
}

var recursionTimesUsers = 1

func (m *Users) ModelUniqueId() (generateId string, err error) {
	generateId, err = utils.IdGenerate(utils.IdTypeUser)
	if err != nil {
		return
	}

	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err == nil {
		recursionTimesServices = 1
		return
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			err = errors.New(enums.CodeMessages(enums.IdConflict))
			return
		}

		recursionTimesServices++
		generateId, err = m.ModelUniqueId()
		if err != nil {
			return
		}

		return
	}
}

func (u *Users) UserInfosByEmailFilterIds(email string, filterIds []string) []Users {
	userInfos := make([]Users, 0)
	db := packages.GetDb().
		Table(u.TableName()).
		Where("email = ?", email)

	if len(filterIds) != 0 {
		db = db.Where("id NOT IN ?", filterIds)
	}

	db.Find(&userInfos)

	return userInfos
}

func (u *Users) UserAdd(userData *Users) error {
	userId, userIdUniqueErr := u.ModelUniqueId()
	if userIdUniqueErr != nil {
		return userIdUniqueErr
	}
	userData.ResID = userId
	userData.Password = utils.Md5(utils.Md5(userData.Password))

	err := packages.GetDb().
		Table(u.TableName()).
		Create(userData).Error

	return err
}

func (u *Users) UserInfoByEmail(email string) Users {
	userInfo := Users{}
	packages.GetDb().
		Table(u.TableName()).
		Where("email = ?", email).
		First(&userInfo)

	return userInfo
}

func (u *Users) UserInfoByName(name string) Users {
	userInfo := Users{}
	packages.GetDb().
		Table(u.TableName()).
		Where("name = ?", name).
		First(&userInfo)

	return userInfo
}

func (u *Users) UserUpdatePassword(email string, newPassword string) error {
	hashedPassword := utils.Md5(utils.Md5(newPassword))
	err := packages.GetDb().
		Table(u.TableName()).
		Where("email = ?", email).
		Update("password", hashedPassword).Error

	return err
}

func (u *Users) UserListPage(param *validators.UserList) (list []Users, total int, listError error) {
	usersModel := Users{}
	tx := packages.GetDb().
		Table(usersModel.TableName())

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("name LIKE ? OR email LIKE ?", search, search)
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

func (u *Users) UserInfoByResId(resId string) (Users, error) {
	userInfo := Users{}
	err := packages.GetDb().
		Table(u.TableName()).
		Where("res_id = ?", resId).
		First(&userInfo).Error

	return userInfo, err
}

func (u *Users) UserUpdate(resId string, userData map[string]interface{}) error {
	if password, ok := userData["password"].(string); ok && password != "" {
		userData["password"] = utils.Md5(utils.Md5(password))
	}
	err := packages.GetDb().
		Table(u.TableName()).
		Where("res_id = ?", resId).
		Updates(&userData).Error

	return err
}

func (u *Users) UserDelete(resId string) error {
	err := packages.GetDb().
		Table(u.TableName()).
		Where("res_id = ?", resId).
		Delete(&Users{}).Error

	return err
}