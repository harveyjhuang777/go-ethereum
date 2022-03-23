package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/controller"
)

func provideController(binder *dig.Container) {
	// Controller
	if err := binder.Provide(controller.NewRestController); err != nil {
		panic(err)
	}
}
