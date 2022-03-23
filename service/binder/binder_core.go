package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/core/block"
)

func provideCore(binder *dig.Container) {
	// core
	if err := binder.Provide(block.NewBlock); err != nil {
		panic(err)
	}
}
