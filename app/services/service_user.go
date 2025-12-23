package services

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/models"
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"errors"
	"fmt"
	"reflect"
	"time"
)

func CheckUserEmailExist(email string, filterIds []string) error {
	userModel := models.Users{}
	userList := userModel.UserInfosByEmailFilterIds(email, filterIds)
	if len(userList) != 0 {
		return errors.New(enums.CodeMessages(enums.UserEmailExist))
	}

	return nil
}

func CheckUserAndPassword(username string, password string) error {
	conf := packages.GetConfig()
	if conf != nil {
		confValue := reflect.ValueOf(conf).Elem()
		ldapField := confValue.FieldByName("Ldap")
		if ldapField.IsValid() && ldapField.FieldByName("Enabled").Bool() {
			return CheckUserAndPasswordWithLdap(username, password)
		}
	}

	userModel := models.Users{}
	userInfo := userModel.UserInfoByName(username)
	if userInfo.Name != username {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	if utils.Md5(utils.Md5(password)) != userInfo.Password {
		return errors.New(enums.CodeMessages(enums.UserPasswordError))
	}

	return nil
}

func UserLogin(username string) (string, error) {
	username = GetUserEmailWithLdap(username)

	token, tokenErr := utils.GenToken(username)
	if tokenErr != nil {
		return "", errors.New(enums.CodeMessages(enums.UserLoggingInError))
	}

	emailExpires, _ := time.ParseDuration(fmt.Sprintf("+%dm", packages.Token.TokenExpire))

	userTokensModel := models.UserTokens{}
	setErr := userTokensModel.SetTokenExpire(username, token, time.Now().Add(emailExpires))
	if setErr != nil {
		return "", setErr
	}

	return token, nil
}

func UserLogout(token string) (bool, error) {
	username, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	userTokensModel := models.UserTokens{}
	userTokenExpire := userTokensModel.GetTokenExpireByEmail(username)

	if len(userTokenExpire.UserEmail) == 0 || userTokenExpire.UserEmail != username {
		return false, errors.New(enums.CodeMessages(enums.UserNoLoggingIn))
	}

	delTokenExpireByEmailErr := userTokensModel.DelTokenExpireByEmail(username)
	if delTokenExpireByEmailErr != nil {
		return false, delTokenExpireByEmailErr
	}

	return true, nil
}

func UserLoginRefresh(token string) (bool, error) {
	username, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	emailExpires, _ := time.ParseDuration(fmt.Sprintf("+%dm", packages.Token.TokenExpire))

	userTokensModel := models.UserTokens{}
	setErr := userTokensModel.SetTokenExpire(username, token, time.Now().Add(emailExpires))
	if setErr != nil {
		return false, setErr
	}

	return true, nil
}

func CheckUserLoginStatus(token string) (bool, error) {
	username, err := utils.ParseToken(token)
	if err != nil {
		return false, errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	userTokensModel := models.UserTokens{}
	userTokenExpire := userTokensModel.GetTokenExpireByEmail(username)

	if len(userTokenExpire.UserEmail) == 0 || userTokenExpire.UserEmail != username {

		return false, errors.New(enums.CodeMessages(enums.UserNoLoggingIn))

	} else {
		if userTokenExpire.Token != token {
			return false, errors.New(enums.CodeMessages(enums.UserTokenError))
		}

		if userTokenExpire.ExpiredAt.Unix() < time.Now().Unix() {
			return false, errors.New(enums.CodeMessages(enums.UserLoggingInExpire))
		}
	}

	return true, nil
}

func UserChangePassword(token string, request *validators.UserChangePassword) error {
	username, err := utils.ParseToken(token)
	if err != nil {
		return errors.New(enums.CodeMessages(enums.UserTokenError))
	}

	err = CheckUserAndPassword(username, request.OldPassword)
	if err != nil {
		return err
	}

	userModel := models.Users{}
	userInfo := userModel.UserInfoByName(username)
	if userInfo.Name != username {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	err = userModel.UserUpdatePassword(userInfo.Email, request.Password)
	if err != nil {
		return errors.New(enums.CodeMessages(enums.Error))
	}

	return nil
}

type UserItem struct {
	ID      int    `json:"id"`
	ResID   string `json:"res_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func UserList(request *validators.UserList) ([]UserItem, int, error) {
	userModel := models.Users{}
	list, total, err := userModel.UserListPage(request)
	if err != nil {
		return []UserItem{}, 0, err
	}

	userList := make([]UserItem, 0)
	for _, v := range list {
		userItem := UserItem{
			ID:      v.ID,
			ResID:   v.ResID,
			Name:    v.Name,
			Email:   v.Email,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		userList = append(userList, userItem)
	}

	return userList, total, nil
}

func UserInfo(resId string) (UserItem, error) {
	userModel := models.Users{}
	user, err := userModel.UserInfoByResId(resId)
	if err != nil {
		return UserItem{}, err
	}

	userItem := UserItem{
		ID:      user.ID,
		ResID:   user.ResID,
		Name:    user.Name,
		Email:   user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return userItem, nil
}

func UserCreate(userData *validators.UserAddUpdate) error {
	err := CheckUserEmailExist(userData.Email, []string{})
	if err != nil {
		return err
	}

	userModel := &models.Users{
		Name:     userData.Name,
		Email:    userData.Email,
		Password: userData.Password,
	}

	addErr := userModel.UserAdd(userModel)
	return addErr
}

func UserUpdate(resId string, userData *validators.UserAddUpdate) error {
	userModel := models.Users{}
	userInfo, err := userModel.UserInfoByResId(resId)
	if err != nil {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	err = CheckUserEmailExist(userData.Email, []string{userInfo.ResID})
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"name":  userData.Name,
		"email": userData.Email,
	}

	if userData.Password != "" {
		updateData["password"] = userData.Password
	}

	err = userModel.UserUpdate(resId, updateData)
	return err
}

func CheckUserExist(resId string) error {
	userModel := models.Users{}
	userInfo, err := userModel.UserInfoByResId(resId)
	if err != nil {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	if userInfo.ResID != resId {
		return errors.New(enums.CodeMessages(enums.UserNull))
	}

	return nil
}

func UserDelete(resId string) error {
	err := CheckUserExist(resId)
	if err != nil {
		return err
	}

	userModel := models.Users{}
	err = userModel.UserDelete(resId)
	return err
}
