package app

import (
	"context"
	"net/http"
	"runtime"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/controller"
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	apiSelf *apiService
	apiOnce sync.Once
)

func NewApiService(in apiServiceIn) IApiService {
	apiOnce.Do(func() {
		apiSelf = &apiService{
			in: in,
		}
	})

	return apiSelf
}

type apiServiceIn struct {
	dig.In
	Config config.IConfig
	Logger logger.ILogger

	BlockController controller.IBlockController
}

type IApiService interface {
	Run(ctx context.Context, stop chan error)
}

type apiService struct {
	in apiServiceIn
}

func (s *apiService) Run(ctx context.Context, stop chan error) {
	engine := s.newEngine()
	engine.SetTrustedProxies(nil)
	s.setRoutes(engine)

	if err := engine.Run(s.in.Config.GetAppConfig().GetGinConfig().Address); err != nil {
		s.in.Logger.Error(ctx, err)
		stop <- err
	}
}

func (s *apiService) newEngine() *gin.Engine {
	return gin.New()
}

func (s *apiService) setRoutes(engine *gin.Engine) {
	// 設定 middlewares
	engine.Use(
		gin.Logger(), // log 之後會換成自定義的 log
		gin.Recovery(),
	)

	// 設定路由
	s.setPublicRoutes(engine)  // 如：Deposit API, Partner API, Admin API
	s.setPrivateRoutes(engine) // 如：pprof, health
}

func (s *apiService) setPublicRoutes(engine *gin.Engine) {
	s.setSrvRoutes(engine) // Gateway 自己的功能
}

func (s *apiService) setSrvRoutes(engine *gin.Engine) {
	privateRouteGroup := engine.Group("")

	// 設定路由
	s.setSrvAPIRoutes(privateRouteGroup)
}

func (s *apiService) setPrivateRoutes(engine *gin.Engine) {
	privateRouteGroup := engine.Group("/_")
	_ = privateRouteGroup

	// health check
	privateRouteGroup.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// pprof
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	pprof.Register(engine, "/_/debug")
}
