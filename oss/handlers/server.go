package handlers

import (
	"net/http"

	"github.com/Yifangmo/micro-shop-web/oss/global"
	"github.com/Yifangmo/micro-shop-web/oss/middlewares"
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

	engine.LoadHTMLFiles("oss/templates/index.html")
	engine.StaticFS("/static", http.Dir("oss/static"))
	engine.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	group := engine.Group("/oss/v1")
	InitOssRouter(group)

	return &http.Server{
		Addr:    addr,
		Handler: engine,
	}
}
