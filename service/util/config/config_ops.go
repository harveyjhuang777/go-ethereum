package config

import "github.com/harveyjhuang777/go-ethereum/service/model"

type IOpsConfig interface {
	GetOpsMySQLConfig() model.MySQLOps
}

func (c *setup) GetOpsMySQLConfig() model.MySQLOps {
	return c.OpsConfig.MySQLOps
}
