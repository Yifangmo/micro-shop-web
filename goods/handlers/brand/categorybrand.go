package brand

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Yifangmo/micro-shop-services/common"
	"github.com/Yifangmo/micro-shop-services/goods/proto"
	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/goods/models"
	"github.com/Yifangmo/micro-shop-web/goods/utils"
)

func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	resp, err := global.GoodsSrvClient.GetCategoryBrandByBrandID(ctx.Request.Context(), &proto.BrandIDRequest{
		Id: int32(i),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func CategoryBrandList(ctx *gin.Context) {
	resp, err := global.GoodsSrvClient.CategoryBrandList(ctx.Request.Context(), &common.PageInfo{})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	resM := map[string]interface{}{
		"total": resp.Total,
	}
	data := make([]interface{}, 0, len(resp.Data))
	for _, value := range resp.Data {
		data = append(data, gin.H{
			"id": value.Id,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
		})
	}
	resM["data"] = data
	ctx.JSON(http.StatusOK, resM)
}

func NewCategoryBrand(ctx *gin.Context) {
	categoryBrandForm := models.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	resp, err := global.GoodsSrvClient.CreateCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": resp.Id})
}

func UpdateCategoryBrand(ctx *gin.Context) {
	categoryBrandForm := models.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteCategoryBrand(ctx.Request.Context(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
