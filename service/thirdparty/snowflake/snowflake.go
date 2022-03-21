package snowflake

import (
	"context"
	"sync"

	"go.uber.org/dig"

	"github.com/bwmarrin/snowflake"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
)

type IIDGenerator interface {
	GenerateInt64() int64
}

type idGenerator struct {
	node *snowflake.Node
}

func NewIDGenerator(in digIn) digOut {
	once.Do(func() {
		node, err := snowflake.NewNode(1)
		if err != nil {
			ctx := context.Background()
			in.Logger.Error(ctx, err)
			panic(err)
		}
		self = &packet{
			in: in,
			digOut: digOut{
				IDGenerator: &idGenerator{node: node},
			},
		}
	})

	return self.digOut
}

func (id *idGenerator) GenerateInt64() int64 {
	return id.node.Generate().Int64()
}

type digIn struct {
	dig.In

	Logger logger.ILogger
}

type packet struct {
	in digIn

	digOut
}

type digOut struct {
	dig.Out

	IDGenerator IIDGenerator
}
