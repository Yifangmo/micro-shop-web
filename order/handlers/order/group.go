package order

import (
	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/order/middlewares"
)

func InitGroup(group *gin.RouterGroup) {
	orderG := group.Group("orders").Use(middlewares.JWTAuth()).Use(middlewares.Trace())
	{
		orderG.GET("", List)       // 订单列表
		orderG.POST("", New)       // 新建订单
		orderG.GET("/:id", Detail) // 订单详情
	}
}
