package block

import (
	"sync"

	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/mysqlcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
)

var (
	once sync.Once
	self *packet
)

func NewBlock(in digIn) digOut {
	once.Do(func() {
		self = &packet{
			in:     in,
			digOut: digOut{},
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	MySQLCli    dbcli.IMySQLClient
	IdGenerator snowflake.IIDGenerator
}

type packet struct {
	in digIn

	digOut
}

type digOut struct {
	dig.Out
}
