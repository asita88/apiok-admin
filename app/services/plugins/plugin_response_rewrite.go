package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var responseRewriteValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required":   "[%s] is a required field,expected type: %s",
		"min_length": "[%s] length must be greater than or equal to %d",
		"max_number": "[%s] must be %d or less",
		"min_number": "[%s] must be %d or greater",
		"oneOf":      "[%s] must be a value that exists in [%s]",
	},
	utils.LocalZh: {
		"required":   "[%s]为必填字段，期望类型:%s",
		"min_length": "[%s]长度必须大于或等于%d",
		"max_number": "[%s]必须小于或等于%d",
		"min_number": "[%s]必须大于或等于%d",
		"oneOf":      "[%s]必须是存在于[%s]中的值",
	},
}

var bodyRewriteTypeList = []string{
	"regex", "replace", "prefix", "suffix",
}

type PluginResponseRewriteConfig struct{}

type BodyRewriteValue struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
	Flags       string `json:"flags"`
	From        string `json:"from"`
	To          string `json:"to"`
	Remove      string `json:"remove"`
	Add         string `json:"add"`
}

type BodyRewrite struct {
	Type  string          `json:"type"`
	Value BodyRewriteValue `json:"value"`
}

type PluginResponseRewrite struct {
	Enabled     bool                 `json:"enabled"`
	Headers     map[string]interface{} `json:"headers"`
	StatusCode  *int                 `json:"status_code"`
	BodyRewrite *BodyRewrite         `json:"body_rewrite"`
}

func NewResponseRewrite() PluginResponseRewriteConfig {
	return PluginResponseRewriteConfig{}
}

func (responseRewriteConfig PluginResponseRewriteConfig) PluginConfigDefault() interface{} {
	pluginResponseRewrite := PluginResponseRewrite{
		Enabled:     true,
		Headers:     map[string]interface{}{},
		StatusCode:  nil,
		BodyRewrite: nil,
	}

	return pluginResponseRewrite
}

func (responseRewriteConfig PluginResponseRewriteConfig) PluginConfigParse(configInfo interface{}) (pluginResponseRewriteConfig interface{}, err error) {
	pluginResponseRewrite := PluginResponseRewrite{
		Enabled:     true,
		Headers:     map[string]interface{}{},
		StatusCode:  nil,
		BodyRewrite: nil,
	}

	var configInfoJson []byte
	_, ok := configInfo.(string)
	if ok {
		configInfoJson = []byte(fmt.Sprint(configInfo))
	} else {
		configInfoJson, err = json.Marshal(configInfo)
		if err != nil {
			return
		}
	}

	err = json.Unmarshal(configInfoJson, &pluginResponseRewrite)
	if err != nil {
		return
	}

	if pluginResponseRewrite.Headers == nil {
		pluginResponseRewrite.Headers = map[string]interface{}{}
	}

	pluginResponseRewriteConfig = pluginResponseRewrite

	return
}

func (responseRewriteConfig PluginResponseRewriteConfig) PluginConfigCheck(configInfo interface{}) error {
	responseRewrite, err := responseRewriteConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginResponseRewrite := responseRewrite.(PluginResponseRewrite)

	return responseRewriteConfig.configValidator(pluginResponseRewrite)
}

func (responseRewriteConfig PluginResponseRewriteConfig) configValidator(config PluginResponseRewrite) error {
	if config.StatusCode != nil {
		if *config.StatusCode < 100 {
			return errors.New(fmt.Sprintf(
				responseRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
				"config.status_code", 100))
		}

		if *config.StatusCode > 599 {
			return errors.New(fmt.Sprintf(
				responseRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
				"config.status_code", 599))
		}
	}

	if config.BodyRewrite != nil {
		bodyRewriteTypeListMap := make(map[string]byte)
		for _, rewriteType := range bodyRewriteTypeList {
			bodyRewriteTypeListMap[rewriteType] = 0
		}

		_, exist := bodyRewriteTypeListMap[config.BodyRewrite.Type]
		if !exist {
			return errors.New(fmt.Sprintf(
				responseRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
				"config.body_rewrite.type", strings.Join(bodyRewriteTypeList, ",")))
		}

		if config.BodyRewrite.Type == "regex" {
			if config.BodyRewrite.Value.Pattern == "" {
				return errors.New(fmt.Sprintf(
					responseRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					"config.body_rewrite.value.pattern", "string"))
			}
		} else if config.BodyRewrite.Type == "replace" {
			if config.BodyRewrite.Value.From == "" {
				return errors.New(fmt.Sprintf(
					responseRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					"config.body_rewrite.value.from", "string"))
			}
		}
	}

	return nil
}

