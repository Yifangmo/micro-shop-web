package address

import (
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
	"github.com/gin-gonic/gin"
)

func InitGroup(g *gin.RouterGroup) {
	addressG := g.Group("address")
	{
		addressG.GET("", middlewares.JWTAuth(), List)
		addressG.DELETE("/:id", middlewares.JWTAuth(), Delete)
		addressG.POST("", middlewares.JWTAuth(), New)
		addressG.PUT("/:id", middlewares.JWTAuth(), Update)
	}
}
