package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

func provideUtil(binder *dig.Container) {
	if err := binder.Provide(config.NewConfig); err != nil {
		panic(err)
	}

	if err := binder.Provide(logger.NewSysLog); err != nil {
		panic(err)
	}
}
