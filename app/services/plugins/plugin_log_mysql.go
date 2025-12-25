package plugins

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var logMysqlValidatorErrorMessages = map[string]map[string]string{
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

type PluginLogMysqlConfig struct{}

type PluginLogMysql struct {
	Enabled            bool     `json:"enabled"`
	Host               string   `json:"host"`
	Port               int      `json:"port"`
	Database           string   `json:"database"`
	User               string   `json:"user"`
	Password           string   `json:"password"`
	TableName          string   `json:"table_name"`
	Timeout            int      `json:"timeout"`
	PoolSize           int      `json:"pool_size"`
	IncludeRequestBody bool     `json:"include_request_body"`
	IncludeResponseBody bool    `json:"include_response_body"`
	IncludeHeaders     []string `json:"include_headers"`
	ExcludeHeaders     []string `json:"exclude_headers"`
	BatchSize          int      `json:"batch_size"`
	BatchTimeout       int      `json:"batch_timeout"`
}

func NewLogMysql() PluginLogMysqlConfig {
	return PluginLogMysqlConfig{}
}

func (logMysqlConfig PluginLogMysqlConfig) PluginConfigDefault() interface{} {
	pluginLogMysql := PluginLogMysql{
		Enabled:            true,
		Host:               "127.0.0.1",
		Port:               3306,
		Database:           "",
		User:               "",
		Password:           "",
		TableName:          "apiok_access_log",
		Timeout:            5000,
		PoolSize:           100,
		IncludeRequestBody: false,
		IncludeResponseBody: false,
		IncludeHeaders:     []string{},
		ExcludeHeaders:     []string{},
		BatchSize:          100,
		BatchTimeout:       5000,
	}

	return pluginLogMysql
}

func (logMysqlConfig PluginLogMysqlConfig) PluginConfigParse(configInfo interface{}) (pluginLogMysqlConfig interface{}, err error) {
	pluginLogMysql := PluginLogMysql{
		Enabled:            true,
		Host:               "127.0.0.1",
		Port:               3306,
		Database:           "",
		User:               "",
		Password:           "",
		TableName:          "apiok_access_log",
		Timeout:            5000,
		PoolSize:           100,
		IncludeRequestBody: false,
		IncludeResponseBody: false,
		IncludeHeaders:     []string{},
		ExcludeHeaders:     []string{},
		BatchSize:          100,
		BatchTimeout:       5000,
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

	err = json.Unmarshal(configInfoJson, &pluginLogMysql)
	if err != nil {
		return
	}

	if pluginLogMysql.Host == "" {
		pluginLogMysql.Host = "127.0.0.1"
	}
	if pluginLogMysql.Port == 0 {
		pluginLogMysql.Port = 3306
	}
	if pluginLogMysql.TableName == "" {
		pluginLogMysql.TableName = "apiok_access_log"
	}
	if pluginLogMysql.Timeout == 0 {
		pluginLogMysql.Timeout = 5000
	}
	if pluginLogMysql.PoolSize == 0 {
		pluginLogMysql.PoolSize = 100
	}
	if pluginLogMysql.BatchSize == 0 {
		pluginLogMysql.BatchSize = 100
	}
	if pluginLogMysql.BatchTimeout == 0 {
		pluginLogMysql.BatchTimeout = 5000
	}
	if pluginLogMysql.IncludeHeaders == nil {
		pluginLogMysql.IncludeHeaders = []string{}
	}
	if pluginLogMysql.ExcludeHeaders == nil {
		pluginLogMysql.ExcludeHeaders = []string{}
	}

	pluginLogMysqlConfig = pluginLogMysql

	return
}

func (logMysqlConfig PluginLogMysqlConfig) PluginConfigCheck(configInfo interface{}) error {
	logMysql, err := logMysqlConfig.PluginConfigParse(configInfo)
	if err != nil {
		return err
	}

	pluginLogMysql := logMysql.(PluginLogMysql)

	return logMysqlConfig.configValidator(pluginLogMysql)
}

func (logMysqlConfig PluginLogMysqlConfig) configValidator(config PluginLogMysql) error {
	if config.Host == "" {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.host", "string"))
	}

	if config.Port < 1 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.port", 1))
	}

	if config.Port > 65535 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.port", 65535))
	}

	if config.Database == "" {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.database", "string"))
	}

	if config.User == "" {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.user", "string"))
	}

	if config.TableName == "" {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["required"],
			"config.table_name", "string"))
	}

	if config.Timeout < 1000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.timeout", 1000))
	}

	if config.Timeout > 60000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.timeout", 60000))
	}

	if config.PoolSize < 1 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.pool_size", 1))
	}

	if config.PoolSize > 1000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.pool_size", 1000))
	}

	if config.BatchSize < 1 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.batch_size", 1))
	}

	if config.BatchSize > 1000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.batch_size", 1000))
	}

	if config.BatchTimeout < 1000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["min_number"],
			"config.batch_timeout", 1000))
	}

	if config.BatchTimeout > 60000 {
		return errors.New(fmt.Sprintf(
			logMysqlValidatorErrorMessages[strings.ToLower(packages.GetValidatorLocale())]["max_number"],
			"config.batch_timeout", 60000))
	}

	return nil
}

