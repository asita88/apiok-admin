package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var logKafkaValidatorErrorMessages = map[string]map[string]string{
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

var logFormatList = []string{
	"json", "text",
}

type PluginLogKafkaConfig struct{}

type PluginLogKafka struct {
	Enabled            bool     `json:"enabled"`
	Brokers            []string `json:"brokers"`
	Topic              string   `json:"topic"`
	Timeout            int      `json:"timeout"`
	KeepaliveTimeout   int      `json:"keepalive_timeout"`
	IncludeRequestBody bool     `json:"include_request_body"`
	IncludeResponseBody bool    `json:"include_response_body"`
	IncludeHeaders     []string `json:"include_headers"`
	ExcludeHeaders     []string `json:"exclude_headers"`
	LogFormat          string   `json:"log_format"`
}

func NewLogKafka() PluginLogKafkaConfig {
	return PluginLogKafkaConfig{}
}

func (logKafkaConfig PluginLogKafkaConfig) PluginConfigDefault() interface{} {
	pluginLogKafka := PluginLogKafka{
		Enabled:            true,
		Brokers:            []string{},
		Topic:              "",
		Timeout:            5000,
		KeepaliveTimeout:   60000,
		IncludeRequestBody: false,
		IncludeResponseBody: false,
		IncludeHeaders:     []string{},
		ExcludeHeaders:     []string{},
		LogFormat:          "json",
	}

	return pluginLogKafka
}

func (logKafkaConfig PluginLogKafkaConfig) PluginConfigParse(configInfo interface{}) (pluginLogKafkaConfig interface{}, err error) {
	pluginLogKafka := PluginLogKafka{
		Enabled:            true,
		Brokers:            []string{},
		Topic:              "",
		Timeout:            5000,
		KeepaliveTimeout:   60000,
		IncludeRequestBody: false,
		IncludeResponseBody: false,
		IncludeHeaders:     []string{},
		ExcludeHeaders:     []string{},
		LogFormat:          "json",
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

	err = json.Unmarshal(configInfoJson, &pluginLogKafka)
	if err != nil {
		return
	}

	if pluginLogKafka.Brokers == nil {
		pluginLogKafka.Brokers = []string{}
	}
	if pluginLogKafka.IncludeHeaders == nil {
		pluginLogKafka.IncludeHeaders = []string{}
	}
	if pluginLogKafka.ExcludeHeaders == nil {
		pluginLogKafka.ExcludeHeaders = []string{}
	}
	if pluginLogKafka.LogFormat == "" {
		pluginLogKafka.LogFormat = "json"
	}

	pluginLogKafkaConfig = pluginLogKafka

	return
}

func (logKafkaConfig PluginLogKafkaConfig) PluginConfigCheck(configInfo interface{}) error {
	logKafka, err := logKafkaConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginLogKafka := logKafka.(PluginLogKafka)

	return logKafkaConfig.configValidator(pluginLogKafka)
}

func (logKafkaConfig PluginLogKafkaConfig) configValidator(config PluginLogKafka) error {
	if len(config.Brokers) == 0 {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.brokers", "array"))
	}

	if config.Topic == "" {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.topic", "string"))
	}

	if config.Timeout < 1000 {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.timeout", 1000))
	}

	if config.Timeout > 60000 {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.timeout", 60000))
	}

	if config.KeepaliveTimeout < 1000 {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.keepalive_timeout", 1000))
	}

	if config.KeepaliveTimeout > 600000 {
		return errors.New(fmt.Sprintf(
			logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.keepalive_timeout", 600000))
	}

	if config.LogFormat != "" {
		logFormatListMap := make(map[string]byte)
		for _, format := range logFormatList {
			logFormatListMap[format] = 0
		}

		_, exist := logFormatListMap[config.LogFormat]
		if !exist {
			return errors.New(fmt.Sprintf(
				logKafkaValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["oneOf"],
				"config.log_format", strings.Join(logFormatList, ",")))
		}
	}

	return nil
}

