package main

import (
	"context"
	"sync"

	"go.uber.org/dig"
	"golang.org/x/xerrors"

	"github.com/harveyjhuang777/go-ethereum/service/app"
	"github.com/harveyjhuang777/go-ethereum/service/binder"
	"github.com/harveyjhuang777/go-ethereum/service/repository"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
)

var (
	srvApp     *SrvApp
	srvSetOnce sync.Once
	stop       = make(chan error, 1)
)

type SrvApp struct {
	dig.In

	RestService app.IService
	MySQLClient dbcli.IMySQLClient
}

func initServer(app SrvApp) {
	srvSetOnce.Do(func() {
		srvApp = &app
	})
}

func main() {
	binder := binder.New()
	if err := binder.Invoke(initServer); err != nil {
		panic(err)
	}

	if srvApp == nil {
		panic(xerrors.New("srvApp is nil"))
	}
	ctx := context.Background()

	// Migration
	db := srvApp.MySQLClient.Session()
	if err := repository.Migration(db); err != nil {
		panic(err)
	}

	srvApp.RestService.Run(ctx, stop)
}
