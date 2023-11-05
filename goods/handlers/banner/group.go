package banner

import (
	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/middlewares"
)

func InitGroup(group *gin.RouterGroup) {
	bannerG := group.Group("banners").Use(middlewares.Trace())
	{
		bannerG.GET("", List)                                                            // 轮播图列表页
		bannerG.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), Delete) // 删除轮播图
		bannerG.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), New)          //新建轮播图
		bannerG.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), Update)    //修改轮播图信息
	}
}
