package ethcli

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
	cli  *ethCli
)

type digIn struct {
	dig.In

	Config config.IConfig
	Logger logger.ILogger
}

type packet struct {
	in digIn

	digOut
}

type digOut struct {
	dig.Out

	EthCli IEthCli
}

func NewEthCli(in digIn) digOut {
	once.Do(func() {
		self = &packet{in: in}
		rawUrl := self.in.Config.GetAppConfig().GetGEthConfig().Endpoint
		client, err := ethclient.Dial(rawUrl)
		if err != nil {
			panic(err)
		}
		self.digOut.EthCli = &ethCli{
			client: client,
		}
	})

	return self.digOut
}

type IEthCli interface {
	GetLatestHeader(ctx context.Context) (*types.Header, error)
	GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
	SubscribeNewHead(ctx context.Context, headers chan *types.Header) (ethereum.Subscription, error)
}

type ethCli struct {
	in digIn

	client *ethclient.Client
}

func (c *ethCli) GetLatestHeader(ctx context.Context) (*types.Header, error) {
	header, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		c.in.Logger.Error(ctx, err)
		return nil, err
	}
	return header, nil
}

func (c *ethCli) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := c.client.BlockByNumber(ctx, number)
	if err != nil {
		c.in.Logger.Error(ctx, err)
		return nil, err
	}
	return block, nil
}

func (c *ethCli) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := c.client.TransactionByHash(ctx, hash)
	if err != nil {
		c.in.Logger.Error(ctx, err)
		return nil, false, err
	}

	return tx, isPending, nil
}

func (c *ethCli) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	rc, err := c.client.TransactionReceipt(ctx, hash)
	if err != nil {
		c.in.Logger.Error(ctx, err)
		return nil, err
	}

	return rc, nil
}

func (c *ethCli) SubscribeNewHead(ctx context.Context, headers chan *types.Header) (ethereum.Subscription, error) {
	sub, err := c.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		c.in.Logger.Error(ctx, err)
		return nil, err
	}
	return sub, nil
}
