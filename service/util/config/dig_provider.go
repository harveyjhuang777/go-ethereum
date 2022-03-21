package config

import "sync"

var (
	once sync.Once
	cfg  *setup
)

func NewConfig() IConfig {
	once.Do(func() {
		cfg = newConfig()
	})

	return cfg
}
