package admin

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var userLoginValidator = validators.UserLogin{}
	if msg, err := packages.ParseRequestParams(c, &userLoginValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	checkUserAndPasswordErr := services.CheckUserAndPassword(userLoginValidator.Email, userLoginValidator.Password)
	if checkUserAndPasswordErr != nil {
		utils.Error(c, checkUserAndPasswordErr.Error())
		return
	}

	token, tokenErr := services.UserLogin(userLoginValidator.Email)
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
