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
	"sync"

	"gorm.io/gorm"
)

type ServicesService struct {
}

var (
	servicesService *ServicesService
	servicesOnce    sync.Once
)

func NewServicesService() *ServicesService {

	servicesOnce.Do(func() {
		servicesService = &ServicesService{}
	})

	return servicesService
}

func CheckServiceExist(serviceResId string) error {
	serviceModel := &models.Services{}
	_, err := serviceModel.ServiceInfoById(serviceResId)

	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) ServiceCreate(request *validators.ServiceAddUpdate) error {

	createServiceData := &models.Services{
		Name:     request.Name,
		Protocol: request.Protocol,
		Enable:   request.Enable,
		Release:  utils.ReleaseStatusU,
	}

	// 处理新字段
	createServiceData.ClientMaxBodySize = request.ClientMaxBodySize
	if request.ChunkedTransferEncoding != nil {
		if *request.ChunkedTransferEncoding {
			enable := 1
			createServiceData.ChunkedTransferEncoding = &enable
		} else {
			disable := 2
			createServiceData.ChunkedTransferEncoding = &disable
		}
	}
	if request.ProxyBuffering != nil {
		if *request.ProxyBuffering {
			enable := 1
			createServiceData.ProxyBuffering = &enable
		} else {
			disable := 2
			createServiceData.ProxyBuffering = &disable
		}
	}
	if request.ProxyCache != nil {
		proxyCacheJson, jsonErr := json.Marshal(request.ProxyCache)
		if jsonErr == nil {
			proxyCacheStr := string(proxyCacheJson)
			createServiceData.ProxyCache = &proxyCacheStr
		}
	}
	if request.ProxySetHeader != nil {
		proxySetHeaderJson, jsonErr := json.Marshal(request.ProxySetHeader)
		if jsonErr == nil {
			proxySetHeaderStr := string(proxySetHeaderJson)
			createServiceData.ProxySetHeader = &proxySetHeaderStr
		}
	}

	_, err := (&models.Services{}).ServiceAdd(createServiceData, request.ServiceDomains)

	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) ServiceUpdate(serviceId string, request *validators.ServiceAddUpdate) error {
	serviceModel := models.Services{}
	serviceInfo, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	updateServiceData := models.Services{
		Name:     request.Name,
		Protocol: request.Protocol,
		Enable:   request.Enable,
		Release:  serviceInfo.Release,
	}
	if serviceInfo.Release == utils.ReleaseStatusY {
		updateServiceData.Release = utils.ReleaseStatusT
	}

	// 处理新字段
	if request.ClientMaxBodySize != nil {
		updateServiceData.ClientMaxBodySize = request.ClientMaxBodySize
	}
	if request.ChunkedTransferEncoding != nil {
		if *request.ChunkedTransferEncoding {
			enable := 1
			updateServiceData.ChunkedTransferEncoding = &enable
		} else {
			disable := 2
			updateServiceData.ChunkedTransferEncoding = &disable
		}
	}
	if request.ProxyBuffering != nil {
		if *request.ProxyBuffering {
			enable := 1
			updateServiceData.ProxyBuffering = &enable
		} else {
			disable := 2
			updateServiceData.ProxyBuffering = &disable
		}
	}
	if request.ProxyCache != nil {
		proxyCacheJson, jsonErr := json.Marshal(request.ProxyCache)
		if jsonErr == nil {
			proxyCacheStr := string(proxyCacheJson)
			updateServiceData.ProxyCache = &proxyCacheStr
		}
	}
	if request.ProxySetHeader != nil {
		proxySetHeaderJson, jsonErr := json.Marshal(request.ProxySetHeader)
		if jsonErr == nil {
			proxySetHeaderStr := string(proxySetHeaderJson)
			updateServiceData.ProxySetHeader = &proxySetHeaderStr
		}
	}

	return serviceModel.ServiceUpdate(serviceId, &updateServiceData, request.ServiceDomains)
}

func (s *ServicesService) ServiceUpdateName(serviceId string, request *validators.ServiceUpdateName) error {
	serviceModel := models.Services{}
	service, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	updateParam := map[string]interface{}{
		"name": request.Name,
	}
	if service.Release == utils.ReleaseStatusY {
		updateParam["release"] = utils.ReleaseStatusT
	}

	return serviceModel.ServiceUpdateColumns(serviceId, updateParam)
}

func checkServiceEnableChange(serviceId string, enable int) error {
	serviceModel := &models.Services{}
	serviceInfo, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	if serviceInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func (s *ServicesService) ServiceSwitchEnable(serviceId string, enable int) error {
	serviceModel := models.Services{}
	service, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	err = checkServiceEnableChange(serviceId, enable)

	updateParam := map[string]interface{}{
		"enable": enable,
	}
	if service.Release == utils.ReleaseStatusY {
		updateParam["release"] = utils.ReleaseStatusT
	}

	return serviceModel.ServiceUpdateColumns(serviceId, updateParam)
}

type StructServiceInfo struct {
	ID                int64    `json:"id"`              // Service id
	ResID             string   `json:"res_id"`          // Service res id
	Name              string   `json:"name"`            // Service name
	Protocol          int      `json:"protocol"`        // Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable            int      `json:"enable"`          // Service enable  1:on  2:off
	Release           int      `json:"release"`         // Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains    []string `json:"service_domains"` // Service Domains
	ClientMaxBodySize *string  `json:"client_max_body_size,omitempty"`
}

func (s *ServicesService) ServiceInfoById(serviceId string) (StructServiceInfo, error) {
	serviceInfo := StructServiceInfo{}

	service, err := (&models.Services{}).ServiceInfoById(serviceId)
	if err != nil {
		return serviceInfo, err
	}

	serviceInfo = StructServiceInfo{
		ID:                service.ID,
		ResID:             service.ResID,
		Name:              service.Name,
		Protocol:          service.Protocol,
		Enable:            service.Enable,
		Release:           service.Release,
		ClientMaxBodySize: service.ClientMaxBodySize,
	}
	serviceDomain, err := (&models.ServiceDomains{}).DomainInfosByServiceIds([]string{serviceId})

	domain := []string{}
	if err == nil {
		for _, v := range serviceDomain {
			domain = append(domain, v.Domain)
		}
	}

	serviceInfo.ServiceDomains = domain

	return serviceInfo, nil
}

func (s *ServicesService) ServiceDelete(serviceId string) error {

	routeModel := models.Routers{}
	routerList := routeModel.RouterInfosByServiceId(serviceId)

	if len(routerList) > 0 {
		return errors.New(enums.CodeMessages(enums.ServiceBindingRouter))
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeService, []string{serviceId})
	if err != nil {
		return err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		// 删除service 信息
		err := (&models.Services{}).ServiceDelete(serviceId)
		if err != nil {
			return errors.New(err.Error())
		}

		// 删除pluginConfig 信息
		if len(pluginConfigList) > 0 {
			for _, v := range pluginConfigList {
				err = tx.Model(&models.PluginConfigs{}).Where("res_id = ?", v.ResID).Delete(&models.PluginConfigs{}).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 删除数据面 service 数据
	apiokDataModel := models.ApiokData{}
	err = apiokDataModel.Delete("services", serviceId)
	if err != nil {
		return err
	}

	// 删除数据面 plugin 数据
	for _, v := range pluginConfigList {
		_ = apiokDataModel.Delete("plugins", v.ResID)
	}

	return nil
}

type pluginConfig struct {
	ResID  string `json:"res_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Icon   string `json:"icon"`
	Type   int    `json:"type"`
	Enable int    `json:"enable"`
}

type ServiceItem struct {
	ID                int64          `json:"id"`
	ResID             string         `json:"res_id"`          // Service id
	Name              string         `json:"name"`            // Service name
	Protocol          int            `json:"protocol"`        // Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable            int            `json:"enable"`          // Service enable  1:on  2:off
	Release           int            `json:"release"`         // Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains    []string       `json:"service_domains"` // Domain name
	ClientMaxBodySize *string        `json:"client_max_body_size,omitempty"`
	PluginList        []pluginConfig `json:"plugin_list"`
}

func (s *ServicesService) ServiceList(request *validators.ServiceList) ([]ServiceItem, int, error) {
	serviceModel := models.Services{}
	searchContent := strings.TrimSpace(request.Search)

	serviceIds := []string{}
	if searchContent != "" {
		services, err := serviceModel.ServiceInfosLikeResIdName(searchContent)
		if err == nil {
			for _, v := range services {
				serviceIds = append(serviceIds, v.ResID)
			}
		}

		serviceDomainModel := models.ServiceDomains{}
		serviceDomains, err := serviceDomainModel.ServiceDomainInfosLikeDomain(searchContent)
		if err == nil {
			for _, v := range serviceDomains {
				serviceIds = append(serviceIds, v.ServiceResID)
			}
		}

		if len(serviceIds) == 0 {
			return []ServiceItem{}, 0, nil
		}
	}
	list, total, err := serviceModel.ServiceList(serviceIds, request)

	if err != nil && err != gorm.ErrRecordNotFound {
		return []ServiceItem{}, 0, err
	}

	serviceList := []ServiceItem{}

	if len(list) == 0 {
		return []ServiceItem{}, 0, nil
	}

	listServiceId := []string{}
	for _, v := range list {
		listServiceId = append(listServiceId, v.ResID)
	}

	domains, err := (&models.ServiceDomains{}).DomainInfosByServiceIds(listServiceId)

	serviceDomainMap := map[string][]models.ServiceDomains{}
	if err == nil {
		for _, v := range domains {
			if _, ok := serviceDomainMap[v.ServiceResID]; !ok {
				serviceDomainMap[v.ServiceResID] = []models.ServiceDomains{}
			}
			serviceDomainMap[v.ServiceResID] = append(serviceDomainMap[v.ServiceResID], v)
		}
	}

	pluginConfigModel := models.PluginConfigs{}
	pluginConfigList, err := pluginConfigModel.PluginConfigListByTargetResIds(models.PluginConfigsTypeService, listServiceId)
	if err != nil {
		return []ServiceItem{}, 0, err
	}

	pluginConfigListMap := make(map[string][]pluginConfig)
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
			return []ServiceItem{}, 0, err
		}

		pluginListMap := make(map[string]models.Plugins)
		for _, pluginInfo := range pluginList {
			pluginListMap[pluginInfo.ResID] = pluginInfo
		}

		for _, pluginConfigInfo := range pluginConfigList {
			_, ok := pluginConfigListMap[pluginConfigInfo.TargetID]
			if !ok {
				pluginConfigListMap[pluginConfigInfo.TargetID] = make([]pluginConfig, 0)
			}

			pluginConfigInfos := pluginConfig{
				ResID:  pluginConfigInfo.ResID,
				Name:   pluginConfigInfo.Name,
				Key:    pluginConfigInfo.PluginKey,
				Enable: pluginConfigInfo.Enable,
				Icon:   pluginListMap[pluginConfigInfo.PluginResID].Icon,
				Type:   pluginListMap[pluginConfigInfo.PluginResID].Type,
			}

			pluginConfigListMap[pluginConfigInfo.TargetID] = append(pluginConfigListMap[pluginConfigInfo.TargetID], pluginConfigInfos)
		}
	}

	for _, v := range list {

		domain := []string{}
		if tmp, ok := serviceDomainMap[v.ResID]; ok {
			for _, vd := range tmp {
				domain = append(domain, vd.Domain)
			}
		}

		serviceItem := ServiceItem{
			ID:                v.ID,
			ResID:             v.ResID,
			Name:              v.Name,
			Protocol:          v.Protocol,
			Enable:            v.Enable,
			Release:           v.Release,
			ServiceDomains:    domain,
			ClientMaxBodySize: v.ClientMaxBodySize,
			PluginList:        make([]pluginConfig, 0),
		}

		if _, ok := pluginConfigListMap[v.ResID]; ok {
			serviceItem.PluginList = pluginConfigListMap[v.ResID]
		}

		serviceList = append(serviceList, serviceItem)
	}

	return serviceList, total, nil
}

func genServiceReleaseSyncRequest(service models.Services, serviceDomains []models.ServiceDomains, pluginConfigs []models.PluginConfigs) rpc.ServicePutRequest {
	protocols := []string{}
	ports := []int{}
	if service.Protocol == 1 {
		protocols = []string{"http"}
		ports = []int{80}
	} else if service.Protocol == 2 {
		protocols = []string{"https"}
		ports = []int{443}
	} else {
		protocols = []string{"http", "https"}
		ports = []int{80, 443}
	}

	domains := []string{}
	for _, v := range serviceDomains {
		domains = append(domains, v.Domain)
	}

	pluginsList := []rpc.ConfigObjectName{}
	for _, v := range pluginConfigs {
		pluginsList = append(pluginsList, rpc.ConfigObjectName{
			Name: v.ResID,
		})
	}
	enable := false
	if service.Enable == utils.EnableOn {
		enable = true
	}
	servicePutRequest := rpc.ServicePutRequest{
		Name:      service.ResID,
		Protocols: protocols,
		Hosts:     domains,
		Ports:     ports,
		Plugins:   pluginsList,
		Enabled:   enable,
	}

	// 处理新字段
	if service.ClientMaxBodySize != nil {
		sizeBytes, err := utils.ParseSizeToBytes(service.ClientMaxBodySize)
		if err == nil {
			servicePutRequest.ClientMaxBodySize = sizeBytes
		}
	}
	if service.ChunkedTransferEncoding != nil {
		enabled := *service.ChunkedTransferEncoding == 1
		servicePutRequest.ChunkedTransferEncoding = &enabled
	}
	if service.ProxyBuffering != nil {
		enabled := *service.ProxyBuffering == 1
		servicePutRequest.ProxyBuffering = &enabled
	}
	if service.ProxyCache != nil && len(*service.ProxyCache) > 0 {
		var proxyCache map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(*service.ProxyCache), &proxyCache); jsonErr == nil {
			servicePutRequest.ProxyCache = proxyCache
		}
	}
	if service.ProxySetHeader != nil && len(*service.ProxySetHeader) > 0 {
		var proxySetHeader map[string]string
		if jsonErr := json.Unmarshal([]byte(*service.ProxySetHeader), &proxySetHeader); jsonErr == nil {
			servicePutRequest.ProxySetHeader = proxySetHeader
		}
	}

	return servicePutRequest
}

func (s *ServicesService) ServiceRelease(serviceId string) error {

	serviceInfo, err := (&models.Services{}).ServiceInfoById(serviceId)

	if err != nil {
		return err
	}
	if serviceInfo.Release == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {

		updateParam := map[string]interface{}{
			"release": utils.ReleaseStatusY,
		}

		err = (&models.Services{}).ServiceUpdateColumnsWithDB(tx, serviceId, updateParam)
		if err != nil {
			packages.Log.Error("service release update mysql data error", err.Error())
			return err
		}

		serviceDomain, err := (&models.ServiceDomains{}).DomainInfosByServiceIds([]string{serviceId})
		if err != nil {
			packages.Log.Error("service release get domains data error", err.Error())
			return err
		}

		successPluginConfig, err := SyncPluginToDataSide(tx, models.PluginConfigsTypeService, serviceId)

		if err != nil {
			packages.Log.Error("service release sync plugin data error", err.Error())
			return err
		}

		request := genServiceReleaseSyncRequest(serviceInfo, serviceDomain, successPluginConfig)

		// 更新数据面 service 数据
		apiokDataModel := models.ApiokData{}
		err = apiokDataModel.Upsert("services", serviceId, request)

		if err != nil {
			packages.Log.Error("service release sync error", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *ServicesService) CheckExistDomain(domains []string, filterServiceIds []string) error {
	serviceDomainInfo := models.ServiceDomains{}
	serviceDomains, err := serviceDomainInfo.DomainInfosByDomain(domains, filterServiceIds)
	if err != nil {
		return nil
	}

	if len(serviceDomains) == 0 {
		return nil
	}

	existDomains := make([]string, 0)
	tmpExistDomainsMap := make(map[string]byte, 0)
	for _, serviceDomain := range serviceDomains {
		_, exist := tmpExistDomainsMap[serviceDomain.Domain]
		if exist {
			continue
		}

		existDomains = append(existDomains, serviceDomain.Domain)
		tmpExistDomainsMap[serviceDomain.Domain] = 0
	}

	if len(existDomains) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainExist), strings.Join(existDomains, ",")))
	}

	return nil
}

func CheckDomainCertificate(protocol int, domains []string) error {
	if (protocol != utils.ProtocolHTTPS) && (protocol != utils.ProtocolHTTPAndHTTPS) {
		return nil
	}

	domainSniInfos, err := utils.InterceptSni(domains)
	if err != nil {
		return err
	}

	certificatesModel := models.Certificates{}
	domainCertificateInfos := certificatesModel.CertificateInfoByDomainSniInfos(domainSniInfos)
	if len(domainCertificateInfos) == len(domainSniInfos) {
		return nil
	}

	nullCertificateDomains := make([]string, 0)
	for _, domainInfo := range domains {
		if len(domainCertificateInfos) == 0 {

			nullCertificateDomains = append(nullCertificateDomains, domainInfo)
		} else {
			for _, domainCertificateInfo := range domainCertificateInfos {

				domainSni := strings.ReplaceAll(domainCertificateInfo.Sni, "*", "")
				if domainInfo[len(domainInfo)-len(domainSni):] != domainSni {
					nullCertificateDomains = append(nullCertificateDomains, domainInfo)
				}
			}
		}
	}

	if len(nullCertificateDomains) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainSslNull), strings.Join(nullCertificateDomains, ",")))
	}

	return nil
}
