package message

import (
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
	"github.com/gin-gonic/gin"
)

func InitGroup(g *gin.RouterGroup) {
	messageG := g.Group("message").Use(middlewares.JWTAuth())
	{
		messageG.GET("", List)
		messageG.POST("", New)
	}
}
