package app

import (
	"context"
	"net/http"
	"runtime"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	self *restService
	once sync.Once
)

func NewRestService(in restServiceIn) IService {
	once.Do(func() {
		self = &restService{
			in: in,
		}
	})

	return self
}

type restServiceIn struct {
	dig.In
	Config config.IConfig
	Logger logger.ILogger
}

type IService interface {
	Run(ctx context.Context, stop chan error)
}

type restService struct {
	in restServiceIn
}

func (s *restService) Run(ctx context.Context, stop chan error) {
	engine := s.newEngine()
	engine.SetTrustedProxies(nil)
	s.setRoutes(engine)

	if err := engine.Run(s.in.Config.GetAppConfig().GetGinConfig().Address); err != nil {
		s.in.Logger.Error(ctx, err)
		stop <- err
	}
}

func (s *restService) newEngine() *gin.Engine {
	return gin.New()
}

func (s *restService) setRoutes(engine *gin.Engine) {
	// 設定 middlewares
	engine.Use(
		gin.Logger(), // log 之後會換成自定義的 log
		gin.Recovery(),

		//s.in.AuthMiddleware.Handle,
	)

	// 設定路由
	s.setPublicRoutes(engine)  // 如：Deposit API, Partner API, Admin API
	s.setPrivateRoutes(engine) // 如：pprof, health
}

func (s *restService) setPublicRoutes(engine *gin.Engine) {
	s.setSrvRoutes(engine) // Gateway 自己的功能
}

func (s *restService) setSrvRoutes(engine *gin.Engine) {
	privateRouteGroup := engine.Group("")

	// 設定路由
	s.setSrvAPIRoutes(privateRouteGroup)
}

func (s *restService) setPrivateRoutes(engine *gin.Engine) {
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
