package binder

import (
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/ethcli"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
)

func provideThirdParty(binder *dig.Container) {
	if err := binder.Provide(dbcli.NewDBClient); err != nil {
		panic(err)
	}

	if err := binder.Provide(snowflake.NewIDGenerator); err != nil {
		panic(err)
	}

	if err := binder.Provide(ethcli.NewEthCli); err != nil {
		panic(err)
	}
}
