package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var trafficTagValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required": "[%s] is a required field,expected type: %s",
	},
	utils.LocalZh: {
		"required": "[%s]为必填字段，期望类型:%s",
	},
}

var httpMethodList = []string{
	"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "TRACE",
}

type PluginTrafficTagConfig struct{}

type MatchRules struct {
	Path    string            `json:"path"`
	Method  interface{}       `json:"method"`
	Headers map[string]string `json:"headers"`
}

type PluginTrafficTag struct {
	MatchRules MatchRules          `json:"match_rules"`
	Tags       map[string]string   `json:"tags"`
}

func NewTrafficTag() PluginTrafficTagConfig {
	return PluginTrafficTagConfig{}
}

func (trafficTagConfig PluginTrafficTagConfig) PluginConfigDefault() interface{} {
	pluginTrafficTag := PluginTrafficTag{
		MatchRules: MatchRules{
			Path:    "",
			Method:  nil,
			Headers: map[string]string{},
		},
		Tags: map[string]string{},
	}

	return pluginTrafficTag
}

func (trafficTagConfig PluginTrafficTagConfig) PluginConfigParse(configInfo interface{}) (pluginTrafficTagConfig interface{}, err error) {
	pluginTrafficTag := PluginTrafficTag{
		MatchRules: MatchRules{
			Path:    "",
			Method:  nil,
			Headers: map[string]string{},
		},
		Tags: map[string]string{},
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

	err = json.Unmarshal(configInfoJson, &pluginTrafficTag)
	if err != nil {
		return
	}

	if pluginTrafficTag.MatchRules.Headers == nil {
		pluginTrafficTag.MatchRules.Headers = map[string]string{}
	}
	if pluginTrafficTag.Tags == nil {
		pluginTrafficTag.Tags = map[string]string{}
	}

	pluginTrafficTagConfig = pluginTrafficTag

	return
}

func (trafficTagConfig PluginTrafficTagConfig) PluginConfigCheck(configInfo interface{}) error {
	trafficTag, err := trafficTagConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginTrafficTag := trafficTag.(PluginTrafficTag)

	return trafficTagConfig.configValidator(pluginTrafficTag)
}

func (trafficTagConfig PluginTrafficTagConfig) configValidator(config PluginTrafficTag) error {
	if len(config.Tags) == 0 {
		return errors.New(fmt.Sprintf(
			trafficTagValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.tags", "object"))
	}

	if config.MatchRules.Method != nil {
		methodListMap := make(map[string]byte)
		for _, method := range httpMethodList {
			methodListMap[method] = 0
		}

		switch v := config.MatchRules.Method.(type) {
		case string:
			_, exist := methodListMap[v]
			if !exist {
				return errors.New(fmt.Sprintf(
					trafficTagValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					"config.match_rules.method", "one of "+strings.Join(httpMethodList, ",")))
			}
		case []interface{}:
			for _, method := range v {
				methodStr, ok := method.(string)
				if !ok {
					continue
				}
				_, exist := methodListMap[methodStr]
				if !exist {
					return errors.New(fmt.Sprintf(
						trafficTagValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
						"config.match_rules.method", "one of "+strings.Join(httpMethodList, ",")))
				}
			}
		}
	}

	return nil
}

