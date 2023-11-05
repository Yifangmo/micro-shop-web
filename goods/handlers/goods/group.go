package goods

import (
	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/middlewares"
)

func InitGroup(group *gin.RouterGroup) {
	goodsG := group.Group("goods").Use(middlewares.Trace())
	{
		goodsG.GET("", List)                                                            		// 获取商品列表
		goodsG.GET("/:id", Detail)                                                      		// 获取商品的详情
		goodsG.GET("/:id/stocks", Stocks)                                               		// 获取商品的库存
		goodsG.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), New)          		// 新建商品
		goodsG.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), Delete) 		// 删除商品
		goodsG.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), Update)			// 更新商品详情
		goodsG.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), UpdateStatus)	// 更新商品状态
	}
}
