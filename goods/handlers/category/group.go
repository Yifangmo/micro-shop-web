package category

import (
	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/middlewares"
)

func InitGroup(group *gin.RouterGroup) {
	categoryG := group.Group("category").Use(middlewares.Trace())
	{
		categoryG.GET("", List)          // 商品类别列表页
		categoryG.DELETE("/:id", Delete) // 删除分类
		categoryG.GET("/:id", Detail)    // 获取分类详情
		categoryG.POST("", New)          //新建分类
		categoryG.PUT("/:id", Update)    //修改分类信息
	}
}
