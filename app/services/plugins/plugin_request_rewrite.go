package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var requestRewriteValidatorErrorMessages = map[string]map[string]string{
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

var uriRewriteTypeList = []string{
	"regex", "replace", "prefix", "suffix",
}

type PluginRequestRewriteConfig struct{}

type UriRewriteValue struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
	Flags       string `json:"flags"`
	From        string `json:"from"`
	To          string `json:"to"`
	Remove      string `json:"remove"`
	Add         string `json:"add"`
}

type UriRewrite struct {
	Type  string          `json:"type"`
	Value UriRewriteValue `json:"value"`
}

type PluginRequestRewrite struct {
	Enabled    bool                 `json:"enabled"`
	UriRewrite *UriRewrite          `json:"uri_rewrite"`
	Headers    map[string]interface{} `json:"headers"`
	QueryArgs  map[string]interface{} `json:"query_args"`
}

func NewRequestRewrite() PluginRequestRewriteConfig {
	return PluginRequestRewriteConfig{}
}

func (requestRewriteConfig PluginRequestRewriteConfig) PluginConfigDefault() interface{} {
	pluginRequestRewrite := PluginRequestRewrite{
		Enabled:    true,
		UriRewrite: nil,
		Headers:    map[string]interface{}{},
		QueryArgs:  map[string]interface{}{},
	}

	return pluginRequestRewrite
}

func (requestRewriteConfig PluginRequestRewriteConfig) PluginConfigParse(configInfo interface{}) (pluginRequestRewriteConfig interface{}, err error) {
	pluginRequestRewrite := PluginRequestRewrite{
		Enabled:    true,
		UriRewrite: nil,
		Headers:    map[string]interface{}{},
		QueryArgs:  map[string]interface{}{},
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

	err = json.Unmarshal(configInfoJson, &pluginRequestRewrite)
	if err != nil {
		return
	}

	if pluginRequestRewrite.Headers == nil {
		pluginRequestRewrite.Headers = map[string]interface{}{}
	}
	if pluginRequestRewrite.QueryArgs == nil {
		pluginRequestRewrite.QueryArgs = map[string]interface{}{}
	}

	pluginRequestRewriteConfig = pluginRequestRewrite

	return
}

func (requestRewriteConfig PluginRequestRewriteConfig) PluginConfigCheck(configInfo interface{}) error {
	requestRewrite, err := requestRewriteConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginRequestRewrite := requestRewrite.(PluginRequestRewrite)

	return requestRewriteConfig.configValidator(pluginRequestRewrite)
}

func (requestRewriteConfig PluginRequestRewriteConfig) configValidator(config PluginRequestRewrite) error {
	if config.UriRewrite != nil {
		uriRewriteTypeListMap := make(map[string]byte)
		for _, rewriteType := range uriRewriteTypeList {
			uriRewriteTypeListMap[rewriteType] = 0
		}

		_, exist := uriRewriteTypeListMap[config.UriRewrite.Type]
		if !exist {
			return errors.New(fmt.Sprintf(
				requestRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
				"config.uri_rewrite.type", strings.Join(uriRewriteTypeList, ",")))
		}

		if config.UriRewrite.Type == "regex" {
			if config.UriRewrite.Value.Pattern == "" {
				return errors.New(fmt.Sprintf(
					requestRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					"config.uri_rewrite.value.pattern", "string"))
			}
		} else if config.UriRewrite.Type == "replace" {
			if config.UriRewrite.Value.From == "" {
				return errors.New(fmt.Sprintf(
					requestRewriteValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					"config.uri_rewrite.value.from", "string"))
			}
		}
	}

	return nil
}

