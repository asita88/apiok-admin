package admin

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/models"
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"strings"

	"github.com/gin-gonic/gin"
)

func PluginTypeList(c *gin.Context) {
	pluginAllTypes := utils.PluginAllTypes()

	utils.Ok(c, pluginAllTypes)
}

type PluginItem struct {
	ResID       string `json:"res_id"`
	PluginKey   string `json:"plugin_key"`
	Icon        string `json:"icon"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

func PluginAddList(c *gin.Context) {
	pluginModel := models.Plugins{}
	pluginList, err := pluginModel.PluginAllList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	list := make([]PluginItem, 0)
	if len(pluginList) > 0 {
		for _, pluginInfo := range pluginList {
			list = append(list, PluginItem{
				ResID:       pluginInfo.ResID,
				PluginKey:   pluginInfo.PluginKey,
				Icon:        pluginInfo.Icon,
				Type:        pluginInfo.Type,
				Description: pluginInfo.Description,
			})
		}
	}

	utils.Ok(c, list)
}

func PluginInfo(c *gin.Context) {
	pluginResId := strings.TrimSpace(c.Param("plugin_res_id"))

	pluginService := services.PluginsService{}
	pluginConfigDefault, err := pluginService.PluginConfigDefault(pluginResId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, pluginConfigDefault)
}

func GlobalPluginConfigList(c *gin.Context) {
	res, err := services.NewPluginsService().PluginConfigList(models.PluginConfigsTypeGlobal, "")

	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	utils.Ok(c, res)
}

func GlobalPluginConfigInfo(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("res_id"))

	if pluginConfigID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	res, err := services.NewPluginsService().PluginConfigInfoByResId(pluginConfigID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, res)
}

func GlobalPluginConfigAdd(c *gin.Context) {
	var request = &validators.ValidatorPluginConfigAdd{
		Type: models.PluginConfigsTypeGlobal,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}
	pluginConfigResId, err := services.NewPluginsService().PluginConfigAdd(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, map[string]string{
		"res_id": pluginConfigResId,
	})
}

func GlobalPluginConfigUpdate(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("res_id"))

	var request = &validators.ValidatorPluginConfigUpdate{
		PluginConfigId: pluginConfigID,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewPluginsService().PluginConfigUpdate(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func GlobalPluginConfigSwitchEnable(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("res_id"))

	var request = &validators.ValidatorPluginConfigSwitchEnable{
		PluginConfigId: pluginConfigID,
	}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewPluginsService().PluginConfigSwitchEnable(pluginConfigID, request.Enable)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func GlobalPluginConfigDelete(c *gin.Context) {
	pluginConfigID := strings.TrimSpace(c.Param("res_id"))

	if pluginConfigID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	err := services.NewPluginsService().PluginConfigDelete(pluginConfigID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}