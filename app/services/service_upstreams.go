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
	"strings"
	"sync"

	"gorm.io/gorm"
)

type ServiceUpstream struct {
}

var (
	serviceUpstream *ServiceUpstream
	upstreamOnce    sync.Once
)

func NewServiceUpstream() *ServiceUpstream {

	upstreamOnce.Do(func() {
		serviceUpstream = &ServiceUpstream{}
	})

	return serviceUpstream
}

type UpstreamItem struct {
	ResID          string               `json:"res_id"`
	Name           string               `json:"name"`
	Algorithm      int                  `json:"algorithm"`
	ConnectTimeout int                  `json:"connect_timeout"`
	WriteTimeout   int                  `json:"write_timeout"`
	ReadTimeout    int                  `json:"read_timeout"`
	Enable         int                  `json:"enable"`
	Release        int                  `json:"release"`
	Check          *UpstreamHealthCheck `json:"check"`
}

type UpstreamHealthCheck struct {
	Enabled  bool   `json:"enabled"`
	Tcp      bool   `json:"tcp"`
	Method   string `json:"method"`
	Host     string `json:"host"`
	Uri      string `json:"uri"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type UpstreamListItem struct {
	UpstreamItem
	NodeList []UpstreamNodeItem `json:"node_list"`
}

func (u *ServiceUpstream) UpstreamListPage(request *validators.UpstreamList) (list []UpstreamListItem, total int, err error) {
	list = make([]UpstreamListItem, 0)
	upstreamModel := models.Upstreams{}
	upstreamNodeModel := models.UpstreamNodes{}
	request.Search = strings.TrimSpace(request.Search)

	upstreamResIds := make([]string, 0)
	upstreamResIdsMap := make(map[string]byte)
	if request.Search != "" {

		nodeList := make([]models.UpstreamNodes, 0)
		nodeList, err = upstreamNodeModel.NodesListBySearch(request.Search)

		if err != nil {
			return
		}

		if len(nodeList) > 0 {
			for _, nodeInfo := range nodeList {

				if _, ok := upstreamResIdsMap[nodeInfo.UpstreamResID]; ok {
					continue
				}

				upstreamResIds = append(upstreamResIds, nodeInfo.UpstreamResID)
				upstreamResIdsMap[nodeInfo.UpstreamResID] = 0
			}
		}
	}

	upstreamList := make([]models.Upstreams, 0)
	upstreamList, total, err = upstreamModel.UpstreamListPage(upstreamResIds, request)

	upstreamResIds = make([]string, 0)
	if len(upstreamList) != 0 {
		for _, upstreamInfo := range upstreamList {
			upstreamResIds = append(upstreamResIds, upstreamInfo.ResID)
			checkUri := upstreamInfo.CheckUri
			if checkUri == "" {
				checkUri = "/"
			}
			checkInterval := upstreamInfo.CheckInterval
			if checkInterval == 0 {
				checkInterval = 1
			}
			checkTimeout := upstreamInfo.CheckTimeout
			if checkTimeout == 0 {
				checkTimeout = 1
			}

			upstreamItem := UpstreamItem{
				ResID:          upstreamInfo.ResID,
				Name:           upstreamInfo.Name,
				Algorithm:      upstreamInfo.Algorithm,
				ConnectTimeout: upstreamInfo.ConnectTimeout,
				WriteTimeout:   upstreamInfo.WriteTimeout,
				ReadTimeout:    upstreamInfo.ReadTimeout,
				Enable:         upstreamInfo.Enable,
				Release:        upstreamInfo.Release,
				Check: &UpstreamHealthCheck{
					Enabled:  upstreamInfo.CheckEnabled == 1,
					Tcp:      upstreamInfo.CheckTcp == 1,
					Method:   upstreamInfo.CheckMethod,
					Host:     upstreamInfo.CheckHost,
					Uri:      checkUri,
					Interval: checkInterval,
					Timeout:  checkTimeout,
				},
			}

			upstreamListItem := UpstreamListItem{
				UpstreamItem: upstreamItem,
				NodeList:     make([]UpstreamNodeItem, 0),
			}
			list = append(list, upstreamListItem)
		}
	}

	upstreamNodeItem := UpstreamNodeItem{}

	nodeList := make([]UpstreamNodeItem, 0)
	nodeList, err = upstreamNodeItem.UpstreamNodeListByUpstreamResIds(upstreamResIds)
	if err != nil {
		return
	}

	if len(nodeList) != 0 {
		nodeListMap := make(map[string][]UpstreamNodeItem)
		for _, nodeInfo := range nodeList {
			nodeListMap[nodeInfo.UpstreamResID] = append(nodeListMap[nodeInfo.UpstreamResID], nodeInfo)
		}

		for key, info := range list {
			if _, ok := nodeListMap[info.ResID]; ok {
				list[key].NodeList = nodeListMap[info.ResID]
			}
		}
	}

	return
}

func (u *ServiceUpstream) CheckExistName(names []string, filterResIds []string) (err error) {
	upstreamModel := models.Upstreams{}

	upstreamInfos := make([]models.Upstreams, 0)
	upstreamInfos, err = upstreamModel.UpstreamInfosByNames(names, filterResIds)
	if err != nil {
		return
	}

	if len(upstreamInfos) != 0 {
		err = errors.New(enums.CodeMessages(enums.NameExist))
	}

	return
}

func (u *ServiceUpstream) CheckUpstreamExist(resId string) (err error) {
	upstreamModel := models.Upstreams{}
	upstreamInfo, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if upstreamInfo.ResID != resId {
		err = errors.New(enums.CodeMessages(enums.UpstreamNull))
		return
	}

	return
}

func (u *ServiceUpstream) CheckUpstreamUse(resId string) (err error) {
	if resId == "" {
		return
	}

	routerModel := models.Routers{}
	routerList := make([]models.Routers, 0)
	routerList, err = routerModel.RouterListByUpstreamResIds([]string{resId})
	if err != nil {
		return
	}

	if len(routerList) == 0 {
		return
	}

	err = errors.New(enums.CodeMessages(enums.UpstreamRouterExist))

	return
}

func (u *ServiceUpstream) UpstreamCreate(request *validators.UpstreamAddUpdate) (err error) {
	upstreamModel := models.Upstreams{}

	createUpstreamData := models.Upstreams{
		Name:           request.Name,
		Algorithm:      request.LoadBalance,
		ConnectTimeout: request.ConnectTimeout,
		WriteTimeout:   request.WriteTimeout,
		ReadTimeout:    request.ReadTimeout,
		Enable:         request.Enable,
		Release:        utils.ReleaseStatusU,
	}

	if request.Check != nil {
		if request.Check.Enabled {
			createUpstreamData.CheckEnabled = 1
		} else {
			createUpstreamData.CheckEnabled = 0
		}
		if request.Check.Tcp {
			createUpstreamData.CheckTcp = 1
		} else {
			createUpstreamData.CheckTcp = 0
		}
		createUpstreamData.CheckMethod = request.Check.Method
		createUpstreamData.CheckHost = request.Check.Host
		createUpstreamData.CheckUri = request.Check.Uri
		if createUpstreamData.CheckUri == "" {
			createUpstreamData.CheckUri = "/"
		}
		createUpstreamData.CheckInterval = request.Check.Interval
		if createUpstreamData.CheckInterval == 0 {
			createUpstreamData.CheckInterval = 1
		}
		createUpstreamData.CheckTimeout = request.Check.Timeout
		if createUpstreamData.CheckTimeout == 0 {
			createUpstreamData.CheckTimeout = 1
		}
	} else {
		createUpstreamData.CheckEnabled = 0
		createUpstreamData.CheckTcp = 1
		createUpstreamData.CheckMethod = ""
		createUpstreamData.CheckHost = ""
		createUpstreamData.CheckUri = "/"
		createUpstreamData.CheckInterval = 1
		createUpstreamData.CheckTimeout = 1
	}

	createUpstreamNodesData := make([]models.UpstreamNodes, 0)
	if len(request.UpstreamNodes) != 0 {
		ipNameIdMap := utils.IpNameIdMap()
		for _, reqNodeInfo := range request.UpstreamNodes {
			var ipType string
			ipType, err = utils.DiscernIP(reqNodeInfo.NodeIp)
			if err != nil {
				return
			}

			tagsJSON := ""
			if len(reqNodeInfo.Tags) > 0 {
				tagsBytes, _ := json.Marshal(reqNodeInfo.Tags)
				tagsJSON = string(tagsBytes)
			}

			nodeData := models.UpstreamNodes{
				NodeIP:     reqNodeInfo.NodeIp,
				IPType:     ipNameIdMap[ipType],
				NodePort:   reqNodeInfo.NodePort,
				NodeWeight: reqNodeInfo.NodeWeight,
				Tags:       tagsJSON,
			}

			createUpstreamNodesData = append(createUpstreamNodesData, nodeData)
		}
	}

	_, err = upstreamModel.UpstreamAdd(createUpstreamData, createUpstreamNodesData)

	return
}

func (u *ServiceUpstream) UpstreamUpdate(resId string, request *validators.UpstreamAddUpdate) (err error) {
	upstreamModel := models.Upstreams{}
	var upstreamInfo models.Upstreams
	upstreamInfo, err = upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {

		updateUpstreamData := map[string]interface{}{
			"algorithm":       request.LoadBalance,
			"read_timeout":    request.ReadTimeout,
			"write_timeout":   request.WriteTimeout,
			"connect_timeout": request.ConnectTimeout,
		}
		if upstreamInfo.Release == utils.ReleaseStatusY {
			updateUpstreamData["release"] = utils.ReleaseStatusT
		}
		if request.Name != "---" {
			name := request.Name
			if name == "" {
				name = upstreamInfo.ResID
			}
			updateUpstreamData["name"] = name
		}
		if request.Check != nil {
			if request.Check.Enabled {
				updateUpstreamData["check_enabled"] = 1
			} else {
				updateUpstreamData["check_enabled"] = 0
			}
			if request.Check.Tcp {
				updateUpstreamData["check_tcp"] = 1
			} else {
				updateUpstreamData["check_tcp"] = 0
			}
			updateUpstreamData["check_method"] = request.Check.Method
			updateUpstreamData["check_host"] = request.Check.Host
			checkUri := request.Check.Uri
			if checkUri == "" {
				checkUri = "/"
			}
			updateUpstreamData["check_uri"] = checkUri
			checkInterval := request.Check.Interval
			if checkInterval == 0 {
				checkInterval = 1
			}
			updateUpstreamData["check_interval"] = checkInterval
			checkTimeout := request.Check.Timeout
			if checkTimeout == 0 {
				checkTimeout = 1
			}
			updateUpstreamData["check_timeout"] = checkTimeout
		}

		if err = tx.Table(upstreamModel.TableName()).
			Where("res_id = ?", resId).
			Updates(updateUpstreamData).Error; err != nil {
			return
		}

		addNodeList, updateNodeList, delNodeResIds := DiffUpstreamNode(resId, request.UpstreamNodes)

		upstreamNodeModel := models.UpstreamNodes{}
		if len(addNodeList) > 0 {
			if err = tx.Create(&addNodeList).Error; err != nil {
				return
			}
		}

		if len(updateNodeList) > 0 {
			for _, updateNodeInfo := range updateNodeList {
				if err = tx.Table(upstreamNodeModel.TableName()).
					Where("res_id = ?", updateNodeInfo.ResID).
					Updates(&updateNodeInfo).Error; err != nil {
					return
				}
			}
		}

		if len(delNodeResIds) > 0 {
			if err = tx.Table(upstreamNodeModel.TableName()).
				Where("res_id in ?", delNodeResIds).
				Delete(&upstreamNodeModel).Error; err != nil {
				return
			}
		}

		return
	})

	return
}

func (u *ServiceUpstream) UpstreamDelete(resId string) (err error) {
	upstreamModel := models.Upstreams{}
	upstreamInfo, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return err
	}

	if upstreamInfo.ResID != resId {
		return
	}

	err = u.CheckUpstreamUse(resId)
	if err != nil {
		return
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{resId})
	if err != nil {
		return err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(upstreamModel.TableName()).
			Where("res_id = ?", upstreamInfo.ResID).
			Delete(&upstreamModel).Error; err != nil {
			return
		}

		if err = tx.Table(upstreamNodeModel.TableName()).
			Where("upstream_res_id = ?", upstreamInfo.ResID).
			Delete(&upstreamNodeModel).Error; err != nil {
			return
		}

		return
	})

	if err != nil {
		return
	}

	apiokDataModel := models.ApiokData{}
	for _, nodeInfo := range upstreamNodeList {
		err = apiokDataModel.Delete("upstream_nodes", nodeInfo.ResID)
		if err != nil {
			return
		}
	}

	err = apiokDataModel.Delete("upstreams", resId)
	if err != nil {
		return
	}

	return
}

func (u *ServiceUpstream) UpstreamSwitchEnable(resId string, enable int) (err error) {
	upstreamModel := models.Upstreams{}
	upstreamInfo, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if upstreamInfo.Enable == enable {
		err = errors.New(enums.CodeMessages(enums.SwitchNoChange))
		return
	}

	updateData := map[string]interface{}{
		"enable": enable,
	}
	if upstreamInfo.Release == utils.ReleaseStatusY {
		updateData["release"] = utils.ReleaseStatusT
	}

	err = upstreamModel.UpstreamUpdateColumns(resId, updateData)

	return
}

func (u *ServiceUpstream) UpstreamSwitchRelease(resId string) (err error) {
	upstreamModel := models.Upstreams{}
	upstreamInfo, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if upstreamInfo.Release == utils.ReleaseStatusY {
		err = errors.New(enums.CodeMessages(enums.SwitchPublished))
		return
	}

	err = UpstreamRelease([]string{resId}, utils.ReleaseTypePush)
	if err != nil {
		return
	}

	releaseStatus := map[string]interface{}{
		"release": utils.ReleaseStatusY,
	}
	err = upstreamModel.UpstreamUpdateColumns(resId, releaseStatus)

	return
}

func (u *ServiceUpstream) UpstreamInfoByResId(resId string) (info UpstreamListItem, err error) {
	upstreamModel := models.Upstreams{}

	upstreamInfo := models.Upstreams{}
	upstreamInfo, err = upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if upstreamInfo.ResID != resId {
		err = errors.New(enums.CodeMessages(enums.UpstreamNull))
		return
	}

	upstreamNodeItem := UpstreamNodeItem{}

	nodeList := make([]UpstreamNodeItem, 0)
	nodeList, err = upstreamNodeItem.UpstreamNodeListByUpstreamResIds([]string{resId})
	if err != nil {
		return
	}

	checkUri := upstreamInfo.CheckUri
	if checkUri == "" {
		checkUri = "/"
	}
	checkInterval := upstreamInfo.CheckInterval
	if checkInterval == 0 {
		checkInterval = 1
	}
	checkTimeout := upstreamInfo.CheckTimeout
	if checkTimeout == 0 {
		checkTimeout = 1
	}

	info.ResID = upstreamInfo.ResID
	info.Name = upstreamInfo.Name
	info.Algorithm = upstreamInfo.Algorithm
	info.ConnectTimeout = upstreamInfo.ConnectTimeout
	info.WriteTimeout = upstreamInfo.WriteTimeout
	info.ReadTimeout = upstreamInfo.ReadTimeout
	info.Enable = upstreamInfo.Enable
	info.Release = upstreamInfo.Release
	info.Check = &UpstreamHealthCheck{
		Enabled:  upstreamInfo.CheckEnabled == 1,
		Tcp:      upstreamInfo.CheckTcp == 1,
		Method:   upstreamInfo.CheckMethod,
		Host:     upstreamInfo.CheckHost,
		Uri:      checkUri,
		Interval: checkInterval,
		Timeout:  checkTimeout,
	}
	info.NodeList = nodeList

	return
}

func (u UpstreamItem) UpstreamDetailByResId(resId string) (upstreamItem UpstreamItem, err error) {

	upstreamModel := models.Upstreams{}
	upstreamDetail, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if len(upstreamDetail.ResID) == 0 {
		return
	}

	checkUri := upstreamDetail.CheckUri
	if checkUri == "" {
		checkUri = "/"
	}
	checkInterval := upstreamDetail.CheckInterval
	if checkInterval == 0 {
		checkInterval = 1
	}
	checkTimeout := upstreamDetail.CheckTimeout
	if checkTimeout == 0 {
		checkTimeout = 1
	}

	upstreamItem.ResID = upstreamDetail.ResID
	upstreamItem.Name = upstreamDetail.Name
	upstreamItem.Algorithm = upstreamDetail.Algorithm
	upstreamItem.ConnectTimeout = upstreamDetail.ConnectTimeout
	upstreamItem.WriteTimeout = upstreamDetail.WriteTimeout
	upstreamItem.ReadTimeout = upstreamDetail.ReadTimeout
	upstreamItem.Check = &UpstreamHealthCheck{
		Enabled:  upstreamDetail.CheckEnabled == 1,
		Tcp:      upstreamDetail.CheckTcp == 1,
		Method:   upstreamDetail.CheckMethod,
		Host:     upstreamDetail.CheckHost,
		Uri:      checkUri,
		Interval: checkInterval,
		Timeout:  checkTimeout,
	}

	return
}

func UpstreamRelease(upstreamResIds []string, releaseType string) (err error) {
	if len(upstreamResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	if releaseType == utils.ReleaseTypePush {
		upstreamModel := models.Upstreams{}
		upstreamList, errList := upstreamModel.UpstreamListByResIds(upstreamResIds)
		if errList != nil {
			return errList
		}

		if len(upstreamList) == 0 {
			return
		}

		upstreamNodeModel := models.UpstreamNodes{}
		apiokDataModel := models.ApiokData{}

		for _, upstreamInfo := range upstreamList {
			upstreamNodeList, errList := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamInfo.ResID})
			if errList != nil {
				return errList
			}

			upstreamConfig, errList := generateUpstreamConfig(upstreamInfo)
			if errList != nil {
				return errList
			}

			errList = apiokDataModel.Upsert("upstreams", upstreamInfo.ResID, upstreamConfig)
			if errList != nil {
				return errList
			}

			for _, nodeInfo := range upstreamNodeList {
				nodeConfig, errList := generateUpstreamNodeConfigForData(nodeInfo, upstreamInfo)
				if errList != nil {
					return errList
				}

				errList = apiokDataModel.Upsert("upstream_nodes", nodeInfo.ResID, nodeConfig)
				if errList != nil {
					return errList
				}
			}
		}

	} else {
		apiokDataModel := models.ApiokData{}

		for _, resId := range upstreamResIds {
			upstreamNodeModel := models.UpstreamNodes{}
			upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{resId})
			if err != nil {
				return err
			}

			for _, nodeInfo := range upstreamNodeList {
				err = apiokDataModel.Delete("upstream_nodes", nodeInfo.ResID)
				if err != nil {
					return err
				}
			}

			err = apiokDataModel.Delete("upstreams", resId)
			if err != nil {
				return err
			}
		}
	}

	return
}

func generateUpstreamConfig(upstreamInfo models.Upstreams) (config rpc.UpstreamConfig, err error) {

	configBalanceList := utils.ConfigBalanceList()
	configBalanceMap := make(map[int]string)
	for _, configBalanceInfo := range configBalanceList {
		configBalanceMap[configBalanceInfo.Id] = configBalanceInfo.Name
	}

	config.Algorithm = utils.ConfigBalanceNameRoundRobin
	configBalance, ok := configBalanceMap[upstreamInfo.Algorithm]
	if ok {
		config.Algorithm = configBalance
	}

	config.Name = upstreamInfo.ResID
	config.ConnectTimeout = upstreamInfo.ConnectTimeout
	config.WriteTimeout = upstreamInfo.WriteTimeout
	config.ReadTimeout = upstreamInfo.ReadTimeout
	config.Nodes = make([]rpc.ConfigObjectName, 0)
	config.Enabled = false
	if upstreamInfo.Enable == utils.EnableOn {
		config.Enabled = true
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList := make([]models.UpstreamNodes, 0)
	upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamInfo.ResID})
	if err != nil {
		return
	}

	if len(upstreamNodeList) != 0 {
		for _, upstreamNodeInfo := range upstreamNodeList {
			config.Nodes = append(config.Nodes, rpc.ConfigObjectName{
				Name: upstreamNodeInfo.ResID,
			})
		}
	}

	return
}
