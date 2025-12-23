package admin

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"

	"github.com/gin-gonic/gin"
)

func LogList(c *gin.Context) {
	var request = &validators.LogList{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.LogList(request)
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

