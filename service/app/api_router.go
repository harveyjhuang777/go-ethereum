package app

import "github.com/gin-gonic/gin"

func (s *apiService) setSrvAPIRoutes(parentRouteGroup *gin.RouterGroup) {
	s.setBlockAPIRoutes(parentRouteGroup)
}

func (s *apiService) setBlockAPIRoutes(parentRouteGroup *gin.RouterGroup) {
	privateRouteGroup := parentRouteGroup.Group("")

	privateRouteGroup.GET("/blocks", s.in.BlockController.GetBlocks)
	privateRouteGroup.GET("/blocks/:id", s.in.BlockController.GetBlockDetail)
	privateRouteGroup.GET("/transaction/:txHash", s.in.BlockController.GetTransactionDetail)
}
