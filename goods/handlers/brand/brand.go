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

func BrandList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	resp, err := global.GoodsSrvClient.BrandList(ctx.Request.Context(), &common.PageInfo{
		PageNumber: uint32(pnInt),
		PageSize:   uint32(pSizeInt),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	resM := make(map[string]interface{})
	resM["total"] = resp.Total
	for _, value := range resp.Data {
		result = append(result, map[string]interface{}{
			"id":   value.Id,
			"name": value.Name,
			"logo": value.Logo,
		})
	}
	resM["data"] = result

	ctx.JSON(http.StatusOK, resM)
}

func NewBrand(ctx *gin.Context) {
	brandForm := models.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	resp, err := global.GoodsSrvClient.CreateBrand(ctx.Request.Context(), &proto.BrandInfoRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = resp.Id
	request["name"] = resp.Name
	request["logo"] = resp.Logo

	ctx.JSON(http.StatusOK, request)
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBrand(ctx.Request.Context(), &proto.BrandIDRequest{Id: int32(i)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateBrand(ctx *gin.Context) {
	brandForm := models.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBrand(ctx.Request.Context(), &proto.BrandInfoRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
