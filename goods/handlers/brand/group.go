package brand

import (
	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-web/goods/middlewares"
)

func InitGroup(group *gin.RouterGroup) {
	brandG := group.Group("brand").Use(middlewares.Trace())
	{
		brandG.GET("", BrandList)          // 品牌列表页
		brandG.DELETE("/:id", DeleteBrand) // 删除品牌
		brandG.POST("", NewBrand)          //新建品牌
		brandG.PUT("/:id", UpdateBrand)    //修改品牌信息
	}

	categoryBrandG := group.Group("categorybrands")
	{
		categoryBrandG.GET("", CategoryBrandList)          // 类别品牌列表页
		categoryBrandG.DELETE("/:id", DeleteCategoryBrand) // 删除类别品牌
		categoryBrandG.POST("", NewCategoryBrand)          //新建类别品牌
		categoryBrandG.PUT("/:id", UpdateCategoryBrand)    //修改类别品牌
		categoryBrandG.GET("/:id", GetCategoryBrandList)   //获取分类的品牌
	}
}
