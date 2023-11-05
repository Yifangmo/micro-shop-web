package banner

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Yifangmo/micro-shop-web/goods/utils"
	"github.com/Yifangmo/micro-shop-web/goods/models"
	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-services/goods/proto"
)

func List(ctx *gin.Context) {
	rsp, err := global.GoodsSrvClient.BannerList(ctx.Request.Context(), &empty.Empty{})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["index"] = value.Index
		reMap["image"] = value.Image
		reMap["url"] = value.Url

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func New(ctx *gin.Context) {
	bannerForm := models.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBanner(ctx.Request.Context(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
		Image: bannerForm.Image,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id
	response["index"] = rsp.Index
	response["url"] = rsp.Url
	response["image"] = rsp.Image

	ctx.JSON(http.StatusOK, response)
}

func Update(ctx *gin.Context) {
	bannerForm := models.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBanner(ctx.Request.Context(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBanner(ctx.Request.Context(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "")
}
