package shoppingcart

import (
	"github.com/Yifangmo/micro-shop-web/order/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGroup(Router *gin.RouterGroup) {
	shoppingCartG := Router.Group("shoppingCart").Use(middlewares.JWTAuth())
	{
		shoppingCartG.GET("", List)          //购物车列表
		shoppingCartG.DELETE("/:id", Delete) //删除条目
		shoppingCartG.POST("", New)          //添加商品到购物车
		shoppingCartG.PATCH("/:id", Update)  //修改条目
	}
}
