package controller

import (
	"sync"

	"go.uber.org/dig"
)

var (
	once sync.Once
	self *packet
)

func NewRestCtl(in digIn) digOut {
	once.Do(func() {
		self = &packet{
			in:     in,
			digOut: digOut{},
		}

	})

	return self.digOut
}

type packet struct {
	in digIn

	digOut
}

type digIn struct {
	dig.In
}

type digOut struct {
	dig.Out
}
