package binder

import (
	"sync"

	"go.uber.org/dig"
)

var (
	binder *dig.Container
	once   sync.Once
)

func New() *dig.Container {
	once.Do(func() {
		binder = dig.New()

		provideApp(binder)
		provideController(binder)
		provideCore(binder)
		provideRepository(binder)
		provideThirdParty(binder)
		provideUtil(binder)
	})

	return binder
}
