package binder

import (
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/app"
)

func provideApp(binder *dig.Container) {
	// App
	if err := binder.Provide(app.NewRestService); err != nil {
		panic(err)
	}
}
