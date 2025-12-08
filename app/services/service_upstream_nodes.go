package services

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/models"
	"apiok-admin/app/rpc"
	"apiok-admin/app/utils"
	"apiok-admin/app/validators"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type UpstreamNodeItem struct {
	ResID         string            `json:"res_id"`
	UpstreamResID string            `json:"upstream_res_id"`
	NodeIP        string            `json:"node_ip"`
	IPType        int               `json:"ip_type"`
	IPTypeName    string            `json:"ip_type_name"`
	NodePort      int               `json:"node_port"`
	NodeWeight    int               `json:"node_weight"`
	Tags          map[string]string `json:"tags"`
}

func (n UpstreamNodeItem) UpstreamNodeListByUpstreamResIds(upstreamResIds []string) (nodeList []UpstreamNodeItem, err error) {
	nodeList = make([]UpstreamNodeItem, 0)
	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds(upstreamResIds)
	if err != nil || len(upstreamNodeList) == 0 {
		return
	}

	iPTypeNameMap := utils.IpIdNameMap()

	for _, upstreamNodeDetail := range upstreamNodeList {
		tags := make(map[string]string)
		if upstreamNodeDetail.Tags != "" {
			json.Unmarshal([]byte(upstreamNodeDetail.Tags), &tags)
		}

		nodeItem := UpstreamNodeItem{
			ResID:         upstreamNodeDetail.ResID,
			UpstreamResID: upstreamNodeDetail.UpstreamResID,
			NodeIP:        upstreamNodeDetail.NodeIP,
			IPType:        upstreamNodeDetail.IPType,
			IPTypeName:    iPTypeNameMap[upstreamNodeDetail.IPType],
			NodePort:      upstreamNodeDetail.NodePort,
			NodeWeight:    upstreamNodeDetail.NodeWeight,
			Tags:          tags,
		}

		nodeList = append(nodeList, nodeItem)
	}

	return
}

func DiffUpstreamNode(upstreamResID string, paramNodeList []validators.UpstreamNodeAddUpdate) (
	addNodeList []models.UpstreamNodes, updateNodeList []models.UpstreamNodes, delNodeResIds []string) {

	if len(upstreamResID) == 0 {
		return
	}

	paramNodeListMap := make(map[string]validators.UpstreamNodeAddUpdate)
	for _, paramNodeInfo := range paramNodeList {
		paramNodeListMapKey := paramNodeInfo.NodeIp + "-" + strconv.Itoa(paramNodeInfo.NodePort)
		paramNodeListMap[paramNodeListMapKey] = paramNodeInfo
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamResID})
	if err != nil {
		return
	}

	upstreamNodeListMap := make(map[string]models.UpstreamNodes)
	for _, upstreamNodeInfo := range upstreamNodeList {
		upstreamNodeListMapKey := upstreamNodeInfo.NodeIP + "-" + strconv.Itoa(upstreamNodeInfo.NodePort)
		upstreamNodeListMap[upstreamNodeListMapKey] = upstreamNodeInfo

		paramNodeInfo, ok := paramNodeListMap[upstreamNodeListMapKey]

		if ok {
			tagsJSON := ""
			if len(paramNodeInfo.Tags) > 0 {
				tagsBytes, _ := json.Marshal(paramNodeInfo.Tags)
				tagsJSON = string(tagsBytes)
			}

			updateNode := models.UpstreamNodes{
				ResID:      upstreamNodeInfo.ResID,
				NodePort:   paramNodeInfo.NodePort,
				NodeWeight: paramNodeInfo.NodeWeight,
				Tags:       tagsJSON,
			}

			updateNodeList = append(updateNodeList, updateNode)
		} else {
			delNodeResIds = append(delNodeResIds, upstreamNodeInfo.ResID)
		}
	}

	ipNameIdMap := utils.IpNameIdMap()

	for _, paramNodeListInfo := range paramNodeList {
		upstreamNodeListMapKey := paramNodeListInfo.NodeIp + "-" + strconv.Itoa(paramNodeListInfo.NodePort)
		_, ok := upstreamNodeListMap[upstreamNodeListMapKey]
		if !ok {

			resId, resIdErr := upstreamNodeModel.ModelUniqueId()
			if resIdErr != nil {
				continue
			}

			ipType, ipTypeErr := utils.DiscernIP(paramNodeListInfo.NodeIp)
			if ipTypeErr != nil {
				continue
			}

			tagsJSON := ""
			if len(paramNodeListInfo.Tags) > 0 {
				tagsBytes, _ := json.Marshal(paramNodeListInfo.Tags)
				tagsJSON = string(tagsBytes)
			}

			newNode := models.UpstreamNodes{
				ResID:         resId,
				UpstreamResID: upstreamResID,
				NodeIP:        paramNodeListInfo.NodeIp,
				IPType:        ipNameIdMap[ipType],
				NodePort:      paramNodeListInfo.NodePort,
				NodeWeight:    paramNodeListInfo.NodeWeight,
				Tags:          tagsJSON,
			}

			addNodeList = append(addNodeList, newNode)
		}
	}

	return
}

func UpstreamNodeLocalCloudDiff(localNodeList []models.UpstreamNodes, cloudNodeList []rpc.UpstreamNodeConfig) (
	putNodeIds []string, deleteNodeIds []string) {

	localNodeListMap := make(map[string]models.UpstreamNodes)
	for _, localNodeInfo := range localNodeList {
		localNodeListMap[localNodeInfo.ResID] = localNodeInfo
		putNodeIds = append(putNodeIds, localNodeInfo.ResID)
	}

	for _, cloudNodeInfo := range cloudNodeList {
		if _, exits := localNodeListMap[cloudNodeInfo.Name]; !exits {
			deleteNodeIds = append(deleteNodeIds, cloudNodeInfo.Name)
		}
	}

	return
}

func generateUpstreamNodeConfig(upstreamNodeInfo models.UpstreamNodes, upstreamInfo models.Upstreams) (rpc.UpstreamNodeConfig, error) {
	return generateUpstreamNodeConfigForData(upstreamNodeInfo, upstreamInfo)
}

func generateUpstreamNodeConfigForData(upstreamNodeInfo models.UpstreamNodes, upstreamInfo models.Upstreams) (config rpc.UpstreamNodeConfig, err error) {
	config.Name = upstreamNodeInfo.ResID
	config.Address = upstreamNodeInfo.NodeIP
	config.Port = upstreamNodeInfo.NodePort
	config.Weight = upstreamNodeInfo.NodeWeight
	config.Health = utils.ConfigHealthY

	if upstreamNodeInfo.Tags != "" {
		tags := make(map[string]string)
		err = json.Unmarshal([]byte(upstreamNodeInfo.Tags), &tags)
		if err == nil && len(tags) > 0 {
			config.Tags = tags
		}
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

	config.Check.Enabled = upstreamInfo.CheckEnabled == 1
	config.Check.Tcp = upstreamInfo.CheckTcp == 1
	config.Check.Method = upstreamInfo.CheckMethod
	config.Check.Host = upstreamInfo.CheckHost
	config.Check.Uri = checkUri
	config.Check.Interval = checkInterval
	config.Check.Timeout = checkTimeout

	return
}

func NodeRelease(nodeResIds []string, releaseType string) (err error) {
	if len(nodeResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	apiokDataModel := models.ApiokData{}

	if releaseType == utils.ReleaseTypeDelete {
		for _, nodeResId := range nodeResIds {
			err = apiokDataModel.Delete("upstream_nodes", nodeResId)
			if err != nil {
				return
			}
		}
		return
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList := make([]models.UpstreamNodes, 0)
	upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByResIds(nodeResIds)
	if err != nil {
		return
	}

	if len(upstreamNodeList) == 0 {
		return
	}

	upstreamResIdsMap := make(map[string]bool)
	for _, nodeInfo := range upstreamNodeList {
		upstreamResIdsMap[nodeInfo.UpstreamResID] = true
	}

	upstreamResIds := make([]string, 0)
	for resId := range upstreamResIdsMap {
		upstreamResIds = append(upstreamResIds, resId)
	}

	upstreamModel := models.Upstreams{}
	upstreamList := make([]models.Upstreams, 0)
	upstreamList, err = upstreamModel.UpstreamListByResIds(upstreamResIds)
	if err != nil {
		return
	}

	upstreamMap := make(map[string]models.Upstreams)
	for _, upstreamInfo := range upstreamList {
		upstreamMap[upstreamInfo.ResID] = upstreamInfo
	}

	for _, upstreamNodeInfo := range upstreamNodeList {
		upstreamInfo, ok := upstreamMap[upstreamNodeInfo.UpstreamResID]
		if !ok {
			continue
		}

		nodeConfig, err := generateUpstreamNodeConfig(upstreamNodeInfo, upstreamInfo)
		if err != nil {
			return err
		}

		err = apiokDataModel.Upsert("upstream_nodes", upstreamNodeInfo.ResID, nodeConfig)
		if err != nil {
			return err
		}
	}

	return
}
