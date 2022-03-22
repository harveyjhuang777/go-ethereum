package dbcli

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
	cli  *DBClient
)

type IMySQLClient interface {
	Session() *gorm.DB
}

type DBClient struct {
	client *gorm.DB
}

func NewDBClient(in digIn) digOut {
	once.Do(func() {
		self = &packet{in: in}
		opsCfg := self.in.Config.GetOpsConfig().GetOpsMySQLConfig()
		self.digOut.DBClient = initWithConfig(opsCfg)
	})

	return self.digOut
}

func initWithConfig(opsCfg model.MySQLOps) IMySQLClient {
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		opsCfg.Username,
		opsCfg.Password,
		opsCfg.Address,
		opsCfg.Database,
	)

	db, err := gorm.Open(mysql.Open(connect))
	if err != nil {
		panic(err)
	}

	self.in.Logger.Info(context.Background(), fmt.Sprintf("Database [%s] Connect success", opsCfg.Database))

	cli = &DBClient{db}

	return cli
}

// Session creates an original gorm.DB session.
func (*DBClient) Session() *gorm.DB {
	return cli.client.Session(&gorm.Session{})
}

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

	DBClient IMySQLClient
}
