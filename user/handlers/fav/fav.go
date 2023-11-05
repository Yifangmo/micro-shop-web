package fav

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	goodsProto "github.com/Yifangmo/micro-shop-services/goods/proto"
	userProto "github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/global"
	"github.com/Yifangmo/micro-shop-web/user/models"
	"github.com/Yifangmo/micro-shop-web/user/utils"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	userFavResp, err := global.UserSrvClient.GetUserFavList(ctx.Request.Context(), &userProto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range userFavResp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	goodsResp, err := global.GoodsSrvClient.GetGoodsByIDs(ctx.Request.Context(), &goodsProto.GoodsIDsRequest{
		Ids: ids,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": userFavResp.Total,
	}

	data := make([]interface{}, 0)
	for _, item := range userFavResp.Data {
		if goods, ok := goodsResp.GoodsMap[item.GoodsId]; ok {
			data = append(data, gin.H{
				"id":         item.GoodsId,
				"name":       goods.Name,
				"shop_price": goods.ShopPrice,
			})
		}
	}
	reMap["data"] = data
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userFavForm := models.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err := global.UserSrvClient.AddUserFav(ctx.Request.Context(), &userProto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserSrvClient.DeleteUserFav(ctx.Request.Context(), &userProto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Detail(ctx *gin.Context) {
	goodsId := ctx.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.UserSrvClient.GetUserFavDetail(ctx.Request.Context(), &userProto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
