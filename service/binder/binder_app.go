package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/app"
)

func provideApp(binder *dig.Container) {
	// App
	if err := binder.Provide(app.NewApiService); err != nil {
		panic(err)
	}

	if err := binder.Provide(app.NewIndexerService); err != nil {
		panic(err)
	}
}
