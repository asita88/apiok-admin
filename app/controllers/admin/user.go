package admin

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var userLoginValidator = validators.UserLogin{}
	if msg, err := packages.ParseRequestParams(c, &userLoginValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	checkUserAndPasswordErr := services.CheckUserAndPassword(userLoginValidator.Username, userLoginValidator.Password)
	if checkUserAndPasswordErr != nil {
		utils.Error(c, checkUserAndPasswordErr.Error())
		return
	}

	token, tokenErr := services.UserLogin(userLoginValidator.Username)
	if tokenErr != nil {
		utils.Error(c, tokenErr.Error())
		return
	}

	type tokenData struct {
		Token string `json:"token"`
	}
	result := tokenData{
		Token: token,
	}

	utils.Ok(c, result)
}

func UserLogout(c *gin.Context) {
	token := c.GetHeader("auth-token")

	loginStatus, loginStatusErr := services.CheckUserLoginStatus(token)
	if (loginStatusErr != nil) || (loginStatus == false) {
		utils.CustomError(c, http.StatusUnauthorized, loginStatusErr.Error())
		return
	}

	_, logoutErr := services.UserLogout(token)
	if logoutErr != nil {
		utils.Error(c, loginStatusErr.Error())
		return
	}

	utils.Ok(c)
}

func UserChangePassword(c *gin.Context) {
	token := c.GetHeader("auth-token")

	var changePasswordValidator = validators.UserChangePassword{}
	if msg, err := packages.ParseRequestParams(c, &changePasswordValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	changePasswordErr := services.UserChangePassword(token, &changePasswordValidator)
	if changePasswordErr != nil {
		utils.Error(c, changePasswordErr.Error())
		return
	}

	utils.Ok(c)
}

func UserList(c *gin.Context) {
	var request = &validators.UserList{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.UserList(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	res := &utils.ResultPage{
		Param:    request,
		Page:     request.Page,
		PageSize: request.PageSize,
		Data:     list,
		Total:    total,
	}

	utils.Ok(c, res)
}

func UserInfo(c *gin.Context) {
	resId := strings.TrimSpace(c.Param("res_id"))

	if resId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	res, err := services.UserInfo(resId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	utils.Ok(c, res)
}

func UserAdd(c *gin.Context) {
	var bindParams = &validators.UserAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	if bindParams.Password == "" {
		utils.Error(c, "密码不能为空")
		return
	}

	err := services.UserCreate(bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func UserUpdate(c *gin.Context) {
	resId := strings.TrimSpace(c.Param("res_id"))

	var bindParams = &validators.UserAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckUserExist(resId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	updateErr := services.UserUpdate(resId, bindParams)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func UserDelete(c *gin.Context) {
	resId := strings.TrimSpace(c.Param("res_id"))

	err := services.CheckUserExist(resId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	err = services.UserDelete(resId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}
