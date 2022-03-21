package config

import "github.com/harveyjhuang777/go-ethereum/service/model"

type IAppConfig interface {
	GetLogConfig() model.LogConfig
	GetGinConfig() model.GinConfig
}

func (c *setup) GetLogConfig() model.LogConfig {
	return c.AppConfig.LogConfig
}

func (c *setup) GetGinConfig() model.GinConfig {
	return c.AppConfig.GinConfig
}
