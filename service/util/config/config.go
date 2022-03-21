package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/harveyjhuang777/go-ethereum/service/model"
)

type IConfig interface {
	GetAppConfig() model.IAppConfig
	GetOpsConfig() model.IOpsConfig
}

type setup struct {
	AppConfig model.AppConfig `mapstructure:"app_config"`
	OpsConfig model.OpsConfig `mapstructure:"ops_config"`
}

func (c *setup) GetAppConfig() model.IAppConfig {
	return cfg
}

func (c *setup) GetOpsConfig() model.IOpsConfig {
	return cfg
}

func newConfig() *setup {
	if cfg != nil {
		return cfg
	}
	cfg = new(setup)

	if err := loadConfig("conf.d/config.yaml"); err != nil {
		panic(err)
	}

	if err := loadConfig("conf.d/app.yaml"); err != nil {
		panic(err)
	}

	return cfg
}

// loadConfig ...
func loadConfig(file string) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.Unmarshal(&cfg)

	return nil
}
