package admin

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"strings"

	"github.com/gin-gonic/gin"
)

func LetsEncryptRequest(c *gin.Context) {
	var request = &validators.LetsEncryptRequest{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	// 默认不启用
	if request.Enable == 0 {
		request.Enable = 2
	}

	resID, err := services.NewLetsEncryptService().RequestCertificate(request.Domain, request.Enable == 1)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, gin.H{
		"res_id":  resID,
		"message": "Certificate requested successfully",
	})
}

func LetsEncryptChallenge(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	keyAuth, ok := services.NewLetsEncryptService().GetChallengeToken(token)
	if !ok {
		utils.Error(c, "Challenge token not found")
		return
	}

	c.String(200, keyAuth)
}
