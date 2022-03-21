package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/mysqlcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
)

func provideThirdParty(binder *dig.Container) {
	if err := binder.Provide(mysqlcli.NewDBClient); err != nil {
		panic(err)
	}

	if err := binder.Provide(snowflake.NewIDGenerator); err != nil {
		panic(err)
	}
}
