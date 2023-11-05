package fav

import (
	"github.com/Yifangmo/micro-shop-web/user/middlewares"
	"github.com/gin-gonic/gin"
)

func InitGroup(g *gin.RouterGroup) {
	UserFavRouter := g.Group("fav")
	{
		UserFavRouter.GET("", middlewares.JWTAuth(), List)          // 获取当前用户的收藏
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), Detail)    // 获取收藏记录
		UserFavRouter.POST("", middlewares.JWTAuth(), New)          // 新建收藏记录
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), Delete) // 删除收藏记录
	}
}
