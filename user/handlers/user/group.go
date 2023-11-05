package user

import (
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
	"github.com/gin-gonic/gin"
)

func InitGroup(g *gin.RouterGroup) {
	userG := g.Group("user")
	{
		userG.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), List)
		userG.POST("loginpwd", LoginByPassword)
		userG.POST("register", Register)
		userG.GET("detail", middlewares.JWTAuth(), Detail)
		userG.PATCH("update", middlewares.JWTAuth(), Update)
	}
}
