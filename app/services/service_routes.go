package services

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/models"
	"apiok-admin/app/packages"
	"apiok-admin/app/rpc"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CheckRouterExist(routerResId string) error {
	routerModel := &models.Routers{}
	routerInfo, err := routerModel.RouterDetailByResId(routerResId)

	if err != nil || len(routerInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterNull))
	}

	return nil
}

func CheckRouterRelease(routerResId string) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerResId, "")

	if len(routerInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterNull))
	}

	if routerInfo.Release == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func CheckServiceRouterPath(path string) error {
	if path == utils.DefaultRouterPath {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathNoPermission))
	}

	if strings.Index(path, utils.DefaultRouterPath) == 0 {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathForbiddenPrefix))
	}

	return nil
}

func CheckExistServiceRouterPath(serviceResId string, path string, filterRouterResIds []string) error {
	routerModel := models.Routers{}
	routerPaths, err := routerModel.RouterInfosByServiceRouterPath(serviceResId, []string{path}, filterRouterResIds)
	if err != nil {
		return err
	}

	if len(routerPaths) == 0 {
		return nil
	}

	existRouterPath := make([]string, 0)
	tmpExistRouterPathMap := make(map[string]byte, 0)
	for _, routerPath := range routerPaths {
		_, exist := tmpExistRouterPathMap[routerPath.RouterPath]
		if exist {
			continue
		}

		existRouterPath = append(existRouterPath, routerPath.RouterPath)
		tmpExistRouterPathMap[routerPath.RouterPath] = 0
	}

	if len(existRouterPath) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.RouterPathExist), strings.Join(existRouterPath, ",")))
	}

	return nil
}

func CheckRouterEnableChange(routerId string, enable int) error {
	routerModel := &models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerId, "")

	if routerInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func RouterCreate(routerData *validators.ValidatorRouterAddUpdate) (routerResId string, err error) {
	createRouterData := models.Routers{
		ServiceResID:   routerData.ServiceResID,
		UpstreamResID:  routerData.UpstreamResID,
		RouterName:     routerData.RouterName,
		RequestMethods: routerData.RequestMethods,
		RouterPath:     routerData.RouterPath,
		Enable:         routerData.Enable,
		Release:        utils.ReleaseStatusU,
	}

	// 处理新字段
	createRouterData.ClientMaxBodySize = routerData.ClientMaxBodySize
	if routerData.ChunkedTransferEncoding != nil {
		if *routerData.ChunkedTransferEncoding {
			enable := 1
			createRouterData.ChunkedTransferEncoding = &enable
		} else {
			disable := 2
			createRouterData.ChunkedTransferEncoding = &disable
		}
	}
	if routerData.ProxyBuffering != nil {
		if *routerData.ProxyBuffering {
			enable := 1
			createRouterData.ProxyBuffering = &enable
		} else {
			disable := 2
			createRouterData.ProxyBuffering = &disable
		}
	}
	if routerData.ProxyCache != nil {
		proxyCacheJson, jsonErr := json.Marshal(routerData.ProxyCache)
		if jsonErr == nil {
			proxyCacheStr := string(proxyCacheJson)
			createRouterData.ProxyCache = &proxyCacheStr
		}
	}
	if routerData.ProxySetHeader != nil {
		proxySetHeaderJson, jsonErr := json.Marshal(routerData.ProxySetHeader)
		if jsonErr == nil {
			proxySetHeaderStr := string(proxySetHeaderJson)
			createRouterData.ProxySetHeader = &proxySetHeaderStr
		}
	}

	routerResId, err = createRouterData.RouterAdd(createRouterData)

	if err != nil {
		return
	}

	return
}

type routerPlugin struct {
	ResID  string `json:"res_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Icon   string `json:"icon"`
	Type   int    `json:"type"`
	Enable int    `json:"enable"`
}

type RouterInfo struct {
	ResId                   string                 `json:"res_id"`
	ServiceResId            string                 `json:"service_res_id"`
	ServiceName             string                 `json:"service_name,omitempty"`
	RouterName              string                 `json:"router_name"`
	RequestMethods          []string               `json:"request_methods"`
	RouterPath              string                 `json:"router_path"`
	Enable                  int                    `json:"enable"`
	Release                 int                    `json:"release"`
	UpstreamResId           string                 `json:"upstream_res_id,omitempty"`
	UpstreamName            string                 `json:"upstream_name,omitempty"`
	ClientMaxBodySize       *string                `json:"client_max_body_size,omitempty"`
	ChunkedTransferEncoding *bool                  `json:"chunked_transfer_encoding,omitempty"`
	ProxyBuffering          *bool                  `json:"proxy_buffering,omitempty"`
	ProxyCache              map[string]interface{} `json:"proxy_cache,omitempty"`
	ProxySetHeader          map[string]string      `json:"proxy_set_header,omitempty"`
	PluginList              []routerPlugin         `json:"plugin_list,omitempty"`
}

func RouterListPage(serviceResId string, param *validators.ValidatorRouterList) (
	routerList []RouterInfo, total int, err error) {

	routerList = make([]RouterInfo, 0)

	routerModel := models.Routers{}
	routerInfos := make([]models.Routers, 0)
	routerInfos, total, err = routerModel.RouterListPage(serviceResId, param)

	routerServiceResIds := make([]string, 0)
	routerServiceResIdsMap := make(map[string]byte)

	routerUpstreamResIds := make([]string, 0)
	routerUpstreamResIdsMap := make(map[string]byte)

	routerResIds := make([]string, 0)
	if len(routerInfos) != 0 {
		for _, routerInfo := range routerInfos {
			routerListItem := RouterInfoFromModel(routerInfo)
			routerList = append(routerList, routerListItem)
			routerResIds = append(routerResIds, routerInfo.ResID)

			if len(routerInfo.ServiceResID) > 0 {
				if _, ok := routerServiceResIdsMap[routerInfo.ServiceResID]; !ok {
					routerServiceResIds = append(routerServiceResIds, routerInfo.ServiceResID)
					routerServiceResIdsMap[routerInfo.ServiceResID] = 0
				}
			}

			if len(routerInfo.UpstreamResID) > 0 {
				if _, ok := routerUpstreamResIdsMap[routerInfo.UpstreamResID]; !ok {
					routerUpstreamResIds = append(routerUpstreamResIds, routerInfo.UpstreamResID)
					routerUpstreamResIdsMap[routerInfo.UpstreamResID] = 0
				}
			}
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, routerResIds)
	if err != nil {
		return
	}

	if len(pluginConfigList) > 0 {

		pluginResIds := make([]string, 0)
		pluginResIdsMap := make(map[string]byte)
		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginResIdsMap[pluginConfigInfo.PluginResID]
			if ok == false {
				pluginResIds = append(pluginResIds, pluginConfigInfo.PluginResID)
			}
		}

		pluginModel := models.Plugins{}
		pluginList := make([]models.Plugins, 0)
		pluginList, err = pluginModel.PluginAllList()
		if err != nil {
			return
		}

		pluginListMap := make(map[string]models.Plugins)
		for _, pluginInfo := range pluginList {
			pluginListMap[pluginInfo.ResID] = pluginInfo
		}

		pluginConfigMapList := make(map[string][]routerPlugin)
		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginConfigMapList[pluginConfigInfo.TargetID]
			if ok == false {
				pluginConfigMapList[pluginConfigInfo.TargetID] = make([]routerPlugin, 0)
			}
			pluginConfigMapList[pluginConfigInfo.TargetID] = append(pluginConfigMapList[pluginConfigInfo.TargetID], routerPlugin{
				ResID:  pluginConfigInfo.ResID,
				Name:   pluginConfigInfo.Name,
				Key:    pluginConfigInfo.PluginKey,
				Enable: pluginConfigInfo.Enable,
				Icon:   pluginListMap[pluginConfigInfo.PluginResID].Icon,
				Type:   pluginListMap[pluginConfigInfo.PluginResID].Type,
			})
		}

		if len(routerList) > 0 {
			for key, routerInfo := range routerList {
				routerPluginList, ok := pluginConfigMapList[routerInfo.ResId]
				if ok {
					routerList[key].PluginList = routerPluginList
				}
			}
		}
	}

	if len(routerServiceResIds) > 0 {

		serviceModel := models.Services{}
		serviceList := make([]models.Services, 0)
		serviceList, err = serviceModel.ServiceListByResIds(routerServiceResIds)
		if err != nil {
			return
		}

		serviceMap := make(map[string]models.Services)
		for _, serviceInfo := range serviceList {
			serviceMap[serviceInfo.ResID] = serviceInfo
		}

		if len(routerList) > 0 {
			for key, routerInfo := range routerList {
				serviceInfo, ok := serviceMap[routerInfo.ServiceResId]
				if ok {
					routerList[key].ServiceName = serviceInfo.Name
				}
			}
		}
	}

	if len(routerUpstreamResIds) > 0 {
		upstreamModel := models.Upstreams{}
		upstreamList := make([]models.Upstreams, 0)
		upstreamList, err = upstreamModel.UpstreamListByResIds(routerUpstreamResIds)
		if err != nil {
			return
		}

		upstreamMap := make(map[string]models.Upstreams)
		for _, upstreamInfo := range upstreamList {
			upstreamMap[upstreamInfo.ResID] = upstreamInfo
		}

		if len(routerList) > 0 {
			for key, routerInfo := range routerList {
				upstreamInfo, ok := upstreamMap[routerInfo.UpstreamResId]
				if ok {
					routerList[key].UpstreamName = upstreamInfo.Name
				}
			}
		}
	}

	return
}

func RouterInfoFromModel(routerModelDetail models.Routers) RouterInfo {
	routerInfo := RouterInfo{
		ResId:             routerModelDetail.ResID,
		ServiceResId:      routerModelDetail.ServiceResID,
		RouterName:        routerModelDetail.RouterName,
		RequestMethods:    strings.Split(routerModelDetail.RequestMethods, ","),
		RouterPath:        routerModelDetail.RouterPath,
		Enable:            routerModelDetail.Enable,
		Release:           routerModelDetail.Release,
		UpstreamResId:     routerModelDetail.UpstreamResID,
		ClientMaxBodySize: routerModelDetail.ClientMaxBodySize,
		PluginList:        make([]routerPlugin, 0),
	}

	if routerModelDetail.ChunkedTransferEncoding != nil {
		enabled := *routerModelDetail.ChunkedTransferEncoding == 1
		routerInfo.ChunkedTransferEncoding = &enabled
	}

	if routerModelDetail.ProxyBuffering != nil {
		enabled := *routerModelDetail.ProxyBuffering == 1
		routerInfo.ProxyBuffering = &enabled
	}

	if routerModelDetail.ProxyCache != nil && len(*routerModelDetail.ProxyCache) > 0 {
		var proxyCache map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(*routerModelDetail.ProxyCache), &proxyCache); jsonErr == nil {
			routerInfo.ProxyCache = proxyCache
		}
	}

	if routerModelDetail.ProxySetHeader != nil && len(*routerModelDetail.ProxySetHeader) > 0 {
		var proxySetHeader map[string]string
		if jsonErr := json.Unmarshal([]byte(*routerModelDetail.ProxySetHeader), &proxySetHeader); jsonErr == nil {
			routerInfo.ProxySetHeader = proxySetHeader
		}
	}

	if len(routerModelDetail.ServiceResID) > 0 {
		serviceModel := models.Services{}
		serviceList, err := serviceModel.ServiceListByResIds([]string{routerModelDetail.ServiceResID})
		if err == nil && len(serviceList) > 0 && len(serviceList[0].ResID) > 0 {
			routerInfo.ServiceName = serviceList[0].Name
		}
	}

	if len(routerModelDetail.UpstreamResID) > 0 {
		upstreamModel := models.Upstreams{}
		upstreamDetail, err := upstreamModel.UpstreamDetailByResId(routerModelDetail.UpstreamResID)
		if err == nil && len(upstreamDetail.ResID) > 0 {
			routerInfo.UpstreamName = upstreamDetail.Name
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, []string{routerModelDetail.ResID})
	if err == nil && len(pluginConfigList) > 0 {
		pluginResIds := make([]string, 0)
		pluginResIdsMap := make(map[string]byte)
		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginResIdsMap[pluginConfigInfo.PluginResID]
			if !ok {
				pluginResIds = append(pluginResIds, pluginConfigInfo.PluginResID)
				pluginResIdsMap[pluginConfigInfo.PluginResID] = 0
			}
		}

		pluginModel := models.Plugins{}
		pluginList, err := pluginModel.PluginAllList()
		if err == nil {
			pluginListMap := make(map[string]models.Plugins)
			for _, pluginInfo := range pluginList {
				pluginListMap[pluginInfo.ResID] = pluginInfo
			}

			pluginListResult := make([]routerPlugin, 0)
			for _, pluginConfigInfo := range pluginConfigList {
				pluginInfo, ok := pluginListMap[pluginConfigInfo.PluginResID]
				if ok {
					pluginListResult = append(pluginListResult, routerPlugin{
						ResID:  pluginConfigInfo.ResID,
						Name:   pluginConfigInfo.Name,
						Key:    pluginConfigInfo.PluginKey,
						Enable: pluginConfigInfo.Enable,
						Icon:   pluginInfo.Icon,
						Type:   pluginInfo.Type,
					})
				}
			}
			routerInfo.PluginList = pluginListResult
		}
	}

	return routerInfo
}

func RouterUpdate(routerResId string, routerData validators.ValidatorRouterAddUpdate) (err error) {
	routerModel := models.Routers{}

	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	updateRouterData := make(map[string]interface{})
	updateRouterData["service_res_id"] = routerData.ServiceResID
	updateRouterData["request_methods"] = routerData.RequestMethods
	updateRouterData["router_path"] = routerData.RouterPath
	updateRouterData["enable"] = routerData.Enable
	updateRouterData["upstream_res_id"] = routerData.UpstreamResID

	if len(routerData.RouterName) != 0 {
		updateRouterData["router_name"] = routerData.RouterName
	}
	if routerDetail.Release == utils.ReleaseStatusY {
		updateRouterData["release"] = utils.ReleaseStatusT
	}

	// 处理新字段
	if routerData.ClientMaxBodySize != nil {
		updateRouterData["client_max_body_size"] = routerData.ClientMaxBodySize
	}
	if routerData.ChunkedTransferEncoding != nil {
		if *routerData.ChunkedTransferEncoding {
			updateRouterData["chunked_transfer_encoding"] = 1
		} else {
			updateRouterData["chunked_transfer_encoding"] = 2
		}
	}
	if routerData.ProxyBuffering != nil {
		if *routerData.ProxyBuffering {
			updateRouterData["proxy_buffering"] = 1
		} else {
			updateRouterData["proxy_buffering"] = 2
		}
	}
	if routerData.ProxyCache != nil {
		proxyCacheJson, jsonErr := json.Marshal(routerData.ProxyCache)
		if jsonErr == nil {
			updateRouterData["proxy_cache"] = string(proxyCacheJson)
		}
	}
	if routerData.ProxySetHeader != nil {
		proxySetHeaderJson, jsonErr := json.Marshal(routerData.ProxySetHeader)
		if jsonErr == nil {
			updateRouterData["proxy_set_header"] = string(proxySetHeaderJson)
		}
	}
	if err = packages.GetDb().Table(routerModel.TableName()).
		Where("res_id = ?", routerResId).
		Updates(&updateRouterData).Error; err != nil {
		return
	}

	return
}

func filterPushedServiceRouterResIds(routerResIds []string) (opRoutersResIds []string, publishedRouterResIds []string, err error) {
	if len(routerResIds) == 0 {
		return
	}

	routerModel := models.Routers{}
	routerList := make([]models.Routers, 0)
	routerList, err = routerModel.RouterListByRouterResIds(routerResIds)
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	serviceResIds := make([]string, 0)

	for _, routerInfo := range routerList {
		if len(routerInfo.ServiceResID) > 0 {
			serviceResIds = append(serviceResIds, routerInfo.ServiceResID)

			if routerInfo.Release != utils.ReleaseStatusU {
				publishedRouterResIds = append(publishedRouterResIds, routerInfo.ResID)
			}
		}
	}

	serviceModel := models.Services{}
	serviceList := make([]models.Services, 0)
	serviceList, err = serviceModel.ServiceListByResIds(serviceResIds)
	if err != nil {
		return
	}

	publishedServiceResIdsMap := make(map[string]byte)
	for _, serviceInfo := range serviceList {
		if serviceInfo.Release != utils.ReleaseStatusU {
			publishedServiceResIdsMap[serviceInfo.ResID] = 0
		}
	}

	for _, routerInfo := range routerList {
		_, ok := publishedServiceResIdsMap[routerInfo.ServiceResID]
		if ok {
			opRoutersResIds = append(opRoutersResIds, routerInfo.ResID)
		}
	}

	return
}

func RouterRelease(routerResIds []string, releaseType string) (err error) {
	if len(routerResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	opRouterResIds := make([]string, 0)
	publishedRouterResIds := make([]string, 0)
	opRouterResIds, publishedRouterResIds, err = filterPushedServiceRouterResIds(routerResIds)

	routerModel := models.Routers{}
	routerList := make([]models.Routers, 0)
	routerList, err = routerModel.RouterListByRouterResIds(opRouterResIds)
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	ApiokDataModel := models.ApiokData{}
	if releaseType == utils.ReleaseTypePush {

		err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {

			for _, routerInfo := range routerList {

				_, err = SyncPluginToDataSide(tx, models.PluginConfigsTypeRouter, routerInfo.ResID)

				if err != nil {
					return
				}

				var routerConfig rpc.RouterConfig
				routerConfig, err = generateRouterConfig(routerInfo)
				if err != nil {
					return
				}

				if len(routerConfig.Name) == 0 {
					continue
				}

				err = routerModel.RouterSwitchRelease(routerInfo.ResID, utils.ReleaseStatusY)
				if err != nil {
					return
				}

				err = ApiokDataModel.Upsert("routers", routerInfo.ResID, routerConfig)
				if err != nil {
					return
				}
			}

			return
		})

	} else {
		for _, resId := range publishedRouterResIds {
			err = ApiokDataModel.Delete("routers", resId)
			if err != nil {
				return
			}
		}
	}

	return
}

func generateRouterConfig(routerInfo models.Routers) (rpc.RouterConfig, error) {
	routerConfig := rpc.RouterConfig{}

	routerConfig.Name = routerInfo.ResID
	routerConfig.Methods = strings.Split(routerInfo.RequestMethods, ",")
	routerConfig.Paths = append(routerConfig.Paths, routerInfo.RouterPath)
	routerConfig.Enabled = false
	if routerInfo.Enable == utils.EnableOn {
		routerConfig.Enabled = true
	}
	routerConfig.Headers = make(map[string]string)
	routerConfig.Service.Name = routerInfo.ServiceResID
	routerConfig.Upstream.Name = routerInfo.UpstreamResID
	routerConfig.Plugins = make([]rpc.ConfigObjectName, 0)

	// 处理新字段
	if routerInfo.ClientMaxBodySize != nil {
		sizeBytes, err := utils.ParseSizeToBytes(routerInfo.ClientMaxBodySize)
		if err == nil {
			routerConfig.ClientMaxBodySize = sizeBytes
		}
	}
	if routerInfo.ChunkedTransferEncoding != nil {
		enabled := *routerInfo.ChunkedTransferEncoding == 1
		routerConfig.ChunkedTransferEncoding = &enabled
	}
	if routerInfo.ProxyBuffering != nil {
		enabled := *routerInfo.ProxyBuffering == 1
		routerConfig.ProxyBuffering = &enabled
	}
	if routerInfo.ProxyCache != nil && len(*routerInfo.ProxyCache) > 0 {
		var proxyCache map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(*routerInfo.ProxyCache), &proxyCache); jsonErr == nil {
			routerConfig.ProxyCache = proxyCache
		}
	}
	if routerInfo.ProxySetHeader != nil && len(*routerInfo.ProxySetHeader) > 0 {
		var proxySetHeader map[string]string
		if jsonErr := json.Unmarshal([]byte(*routerInfo.ProxySetHeader), &proxySetHeader); jsonErr == nil {
			routerConfig.ProxySetHeader = proxySetHeader
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, []string{routerInfo.ResID})
	if err != nil {
		return routerConfig, err
	}

	if len(pluginConfigList) > 0 {
		for _, pluginConfigInfo := range pluginConfigList {
			if pluginConfigInfo.Enable == utils.EnableOff {
				continue
			}

			routerConfig.Plugins = append(routerConfig.Plugins, rpc.ConfigObjectName{
				Name: pluginConfigInfo.ResID,
			})
		}
	}

	return routerConfig, nil
}

func CheckEditDefaultPathRouter(routerId string) error {
	routerModel := models.Routers{}
	routerInfo := routerModel.RouterDetailByResIdServiceResId(routerId, "")
	if routerInfo.RouterPath == utils.DefaultRouterPath {
		return errors.New(enums.CodeMessages(enums.RouterDefaultPathNoPermission))
	}

	return nil
}

func RouterDelete(routerResId string) (err error) {
	routerModel := models.Routers{}

	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	if routerDetail.ResID != routerResId {
		return
	}

	if err = packages.GetDb().Table(routerModel.TableName()).
		Where("res_id = ?", routerResId).
		Delete(&routerModel).Error; err != nil {
		return
	}

	err = RouterRelease([]string{routerResId}, utils.ReleaseTypeDelete)
	if err != nil {
		return
	}

	return
}

func RouterCopy(routerResId string) (err error) {
	routerModel := models.Routers{}
	var routerDetail models.Routers
	routerDetail, err = routerModel.RouterDetailByResId(routerResId)
	if err != nil {
		return
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList := make([]models.PluginConfigs, 0)
	pluginConfigList, err = pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeRouter, []string{routerResId})
	if err != nil {
		return
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {
		newRouterResId, err := routerModel.ModelUniqueId()
		if err != nil {
			return
		}

		randomStr := utils.RandomStrGenerate(4)
		err = tx.Table(routerModel.TableName()).Create(&models.Routers{
			ResID:          newRouterResId,
			ServiceResID:   routerDetail.ServiceResID,
			UpstreamResID:  routerDetail.UpstreamResID,
			RequestMethods: routerDetail.RequestMethods,
			RouterName:     routerDetail.RouterName + "-copy-" + randomStr,
			RouterPath:     routerDetail.RouterPath + "-copy-" + randomStr,
			Enable:         routerDetail.Enable,
			Release:        utils.ReleaseStatusU,
		}).Error
		if err != nil {
			return
		}

		newRouterPluginConfig := make([]models.PluginConfigs, 0)
		if len(pluginConfigList) > 0 {
			for _, pluginConfigInfo := range pluginConfigList {
				var pluginConfigresId string
				pluginConfigresId, err = pluginConfigModel.ModelUniqueId()
				if err != nil {
					return
				}

				newRouterPluginConfig = append(newRouterPluginConfig, models.PluginConfigs{
					ResID:       pluginConfigresId,
					Name:        pluginConfigInfo.Name,
					Type:        models.PluginConfigsTypeRouter,
					TargetID:    newRouterResId,
					PluginResID: pluginConfigInfo.PluginResID,
					PluginKey:   pluginConfigInfo.PluginKey,
					Config:      pluginConfigInfo.Config,
					Enable:      pluginConfigInfo.Enable,
				})
			}

			err = tx.Table(pluginConfigModel.TableName()).Create(&newRouterPluginConfig).Error
			if err != nil {
				return
			}
		}

		return
	})

	return
}

func RouterInfoByResId(resId string) (models.Routers, error) {
	return (&models.Routers{}).RouterDetailByResId(resId)
}
