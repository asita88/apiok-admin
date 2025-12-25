package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var wafValidatorErrorMessages = map[string]map[string]string{
	utils.LocalEn: {
		"required":   "[%s] is a required field,expected type: %s",
		"max_length": "[%s] length must be less than or equal to %d",
		"min_length": "[%s] length must be greater than or equal to %d",
		"max_number": "[%s] must be %d or less",
		"min_number": "[%s] must be %d or greater",
		"oneOf":      "[%s] must be a value that exists in [%s]",
	},
	utils.LocalZh: {
		"required":   "[%s]为必填字段，期望类型:%s",
		"max_length": "[%s]长度必须小于或等于%d",
		"min_length": "[%s]长度必须大于或等于%d",
		"max_number": "[%s]必须小于或等于%d",
		"min_number": "[%s]必须大于或等于%d",
		"oneOf":      "[%s]必须是存在于[%s]中的值",
	},
}

var matchTypeList = []string{
	"uri", "args", "header", "body", "all", "method", "request_size",
}

var operatorList = []string{
	"match", "not_match",
}

var actionList = []string{
	"block", "log",
}

type PluginWafConfig struct{}

type IPWhitelist struct {
	Enabled bool     `json:"enabled"`
	IPList  []string `json:"ip_list"`
}

type IPBlacklist struct {
	Enabled bool     `json:"enabled"`
	IPList  []string `json:"ip_list"`
}

type RuleCondition struct {
	Patterns  []string `json:"patterns"`
	MatchType string   `json:"match_type"`
	Operator  string   `json:"operator"`
}

type Rule struct {
	Name       string          `json:"name"`
	Conditions []RuleCondition `json:"conditions"`
	Action     string          `json:"action"`
}

type Rules struct {
	RuleList []Rule `json:"rule_list"`
}

type PluginWaf struct {
	Enabled      bool        `json:"enabled"`
	IPWhitelist  IPWhitelist `json:"ip_whitelist"`
	IPBlacklist  IPBlacklist `json:"ip_blacklist"`
	Rules        Rules       `json:"rules"`
}

func NewWaf() PluginWafConfig {
	return PluginWafConfig{}
}

func (wafConfig PluginWafConfig) PluginConfigDefault() interface{} {
	pluginWaf := PluginWaf{
		Enabled: true,
		IPWhitelist: IPWhitelist{
			Enabled: true,
			IPList:  []string{},
		},
		IPBlacklist: IPBlacklist{
			Enabled: true,
			IPList:  []string{},
		},
		Rules: Rules{
			RuleList: []Rule{},
		},
	}

	return pluginWaf
}

func (wafConfig PluginWafConfig) PluginConfigParse(configInfo interface{}) (pluginWafConfig interface{}, err error) {
	pluginWaf := PluginWaf{
		Enabled: true,
		IPWhitelist: IPWhitelist{
			Enabled: true,
			IPList:  []string{},
		},
		IPBlacklist: IPBlacklist{
			Enabled: true,
			IPList:  []string{},
		},
		Rules: Rules{
			RuleList: []Rule{},
		},
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

	err = json.Unmarshal(configInfoJson, &pluginWaf)
	if err != nil {
		return
	}

	if pluginWaf.IPWhitelist.IPList == nil {
		pluginWaf.IPWhitelist.IPList = []string{}
	}
	if pluginWaf.IPBlacklist.IPList == nil {
		pluginWaf.IPBlacklist.IPList = []string{}
	}
	if pluginWaf.Rules.RuleList == nil {
		pluginWaf.Rules.RuleList = []Rule{}
	}

	for i := range pluginWaf.Rules.RuleList {
		if pluginWaf.Rules.RuleList[i].Conditions == nil {
			pluginWaf.Rules.RuleList[i].Conditions = []RuleCondition{}
		}
		for j := range pluginWaf.Rules.RuleList[i].Conditions {
			if pluginWaf.Rules.RuleList[i].Conditions[j].Patterns == nil {
				pluginWaf.Rules.RuleList[i].Conditions[j].Patterns = []string{}
			}
			if pluginWaf.Rules.RuleList[i].Conditions[j].Operator == "" {
				pluginWaf.Rules.RuleList[i].Conditions[j].Operator = "match"
			}
		}
		if pluginWaf.Rules.RuleList[i].Action == "" {
			pluginWaf.Rules.RuleList[i].Action = "block"
		}
	}

	pluginWafConfig = pluginWaf

	return
}

func (wafConfig PluginWafConfig) PluginConfigCheck(configInfo interface{}) error {
	waf, err := wafConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginWaf := waf.(PluginWaf)

	return wafConfig.configValidator(pluginWaf)
}

func (wafConfig PluginWafConfig) configValidator(config PluginWaf) error {
	if config.IPWhitelist.Enabled {
		if len(config.IPWhitelist.IPList) == 0 {
			return errors.New(fmt.Sprintf(
				wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
				"config.ip_whitelist.ip_list", "array"))
		}
	}

	if config.IPBlacklist.Enabled {
		if len(config.IPBlacklist.IPList) == 0 {
			return errors.New(fmt.Sprintf(
				wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
				"config.ip_blacklist.ip_list", "array"))
		}
	}

	for i, rule := range config.Rules.RuleList {
		if rule.Name == "" {
			return errors.New(fmt.Sprintf(
				wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
				fmt.Sprintf("config.rules.rule_list[%d].name", i), "string"))
		}

		if len(rule.Conditions) == 0 {
			return errors.New(fmt.Sprintf(
				wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
				fmt.Sprintf("config.rules.rule_list[%d].conditions", i), "array"))
		}

		matchTypeListMap := make(map[string]byte)
		for _, matchType := range matchTypeList {
			matchTypeListMap[matchType] = 0
		}

		operatorListMap := make(map[string]byte)
		for _, operator := range operatorList {
			operatorListMap[operator] = 0
		}

		actionListMap := make(map[string]byte)
		for _, action := range actionList {
			actionListMap[action] = 0
		}

		if rule.Action != "" {
			_, exist := actionListMap[rule.Action]
			if !exist {
				return errors.New(fmt.Sprintf(
					wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
					fmt.Sprintf("config.rules.rule_list[%d].action", i), strings.Join(actionList, ",")))
			}
		}

		for j, condition := range rule.Conditions {
			if condition.MatchType == "" {
				return errors.New(fmt.Sprintf(
					wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					fmt.Sprintf("config.rules.rule_list[%d].conditions[%d].match_type", i, j), "string"))
			}

			_, exist := matchTypeListMap[condition.MatchType]
			if !exist {
				return errors.New(fmt.Sprintf(
					wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
					fmt.Sprintf("config.rules.rule_list[%d].conditions[%d].match_type", i, j), strings.Join(matchTypeList, ",")))
			}

			if len(condition.Patterns) == 0 {
				return errors.New(fmt.Sprintf(
					wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
					fmt.Sprintf("config.rules.rule_list[%d].conditions[%d].patterns", i, j), "array"))
			}

			if condition.Operator != "" {
				_, exist := operatorListMap[condition.Operator]
				if !exist {
					return errors.New(fmt.Sprintf(
						wafValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
						fmt.Sprintf("config.rules.rule_list[%d].conditions[%d].operator", i, j), strings.Join(operatorList, ",")))
				}
			}
		}
	}

	return nil
}

