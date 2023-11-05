package handlers

import (
	"net/http"

	"github.com/Yifangmo/micro-shop-web/order/global"
	"github.com/Yifangmo/micro-shop-web/order/handlers/order"
	"github.com/Yifangmo/micro-shop-web/order/handlers/pay"
	"github.com/Yifangmo/micro-shop-web/order/handlers/shoppingcart"
	"github.com/Yifangmo/micro-shop-web/order/middlewares"

	"github.com/gin-gonic/gin"
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

	group := engine.Group("/o/v1")
	order.InitGroup(group)
	pay.InitGroup(group)
	shoppingcart.InitGroup(group)

	return &http.Server{
		Addr:    addr,
		Handler: engine,
	}
}
