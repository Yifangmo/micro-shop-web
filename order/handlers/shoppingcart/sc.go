package shoppingcart

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Yifangmo/micro-shop-services/order/proto"
	"github.com/Yifangmo/micro-shop-web/order/global"
	"github.com/Yifangmo/micro-shop-web/order/models"
	"github.com/Yifangmo/micro-shop-web/order/utils"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	orderResp, err := global.OrderSrvClient.ShoppingCartList(ctx.Request.Context(), &proto.ShoppingCartRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	goodsIDs := make([]int32, 0, len(orderResp.Data))
	for _, item := range orderResp.Data {
		goodsIDs = append(goodsIDs, item.GoodsId)
	}
	if len(goodsIDs) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	goodsResp, err := global.GoodsSrvClient.GetGoodsByIDs(ctx.Request.Context(), &proto.GoodsIDsRequest{
		Ids: goodsIDs,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	resM := gin.H{
		"total": orderResp.Total,
	}
	goodsList := make([]gin.H, 0, len(goodsIDs))
	for _, order := range orderResp.Data {
		if goods, ok := goodsResp.GoodsMap[order.GoodsId]; ok {
			goodsList = append(goodsList, gin.H{
				"id":          order.Id,
				"goods_id":    order.GoodsId,
				"good_name":   goods.Name,
				"good_image":  goods.PreviewImage,
				"good_price":  goods.ShopPrice,
				"good_number": order.GoodsNumber,
				"checked":     order.Checked,
			})
		}
	}
	resM["data"] = goodsList
	ctx.JSON(http.StatusOK, resM)
}

func New(ctx *gin.Context) {
	form := models.ShoppingCartItemForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	_, err := global.GoodsSrvClient.GetGoodsByID(ctx.Request.Context(), &proto.GoodsIDRequest{
		Id: form.GoodsId,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	inventoryResp, err := global.InventorySrvClient.GetInventory(ctx.Request.Context(), &proto.GoodsInventory{
		GoodsId: form.GoodsId,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if inventoryResp.Num < form.Number {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	resp, err := global.OrderSrvClient.CreateShoppingCartItem(ctx.Request.Context(), &proto.ShoppingCartRequest{
		GoodsId:     form.GoodsId,
		UserId:      int32(userId.(uint)),
		GoodsNumber: form.Number,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": resp.Id,
	})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := models.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	req := &proto.ShoppingCartRequest{
		UserId:      int32(userId.(uint)),
		GoodsId:     int32(i),
		GoodsNumber: itemForm.Nums,
	}
	if itemForm.Checked != nil {
		req.Checked = wrapperspb.Bool(*itemForm.Checked)
	}

	_, err = global.OrderSrvClient.UpdateShoppingCartItem(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteShoppingCartItem(ctx.Request.Context(), &proto.ShoppingCartRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
