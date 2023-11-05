package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/user/handlers/address"
	"github.com/Yifangmo/micro-shop-web/user/handlers/base"
	"github.com/Yifangmo/micro-shop-web/user/handlers/fav"
	"github.com/Yifangmo/micro-shop-web/user/handlers/message"
	"github.com/Yifangmo/micro-shop-web/user/handlers/user"
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
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
	
	group := engine.Group("/u/v1")
	base.InitGroup(group)
	user.InitGroup(group)
	address.InitGroup(group)
	message.InitGroup(group)
	fav.InitGroup(group)

	return &http.Server{
		Addr: addr,
		Handler: engine,
	}
}
