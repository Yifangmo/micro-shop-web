package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/goods/handlers/banner"
	"github.com/Yifangmo/micro-shop-web/goods/handlers/brand"
	"github.com/Yifangmo/micro-shop-web/goods/handlers/category"
	"github.com/Yifangmo/micro-shop-web/goods/handlers/goods"
	"github.com/Yifangmo/micro-shop-web/goods/middlewares"
)

func NewServer(addr string) *http.Server {
	if !global.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.Use(middlewares.Cors())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	group := engine.Group("/g/v1")
	banner.InitGroup(group)
	brand.InitGroup(group)
	category.InitGroup(group)
	goods.InitGroup(group)

	return &http.Server{
		Addr:    addr,
		Handler: engine,
	}
}
