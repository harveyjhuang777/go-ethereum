package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/repository"
)

func provideRepository(binder *dig.Container) {
	if err := binder.Provide(repository.NewRepository); err != nil {
		panic(err)
	}
}
