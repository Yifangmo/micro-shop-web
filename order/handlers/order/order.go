package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"

	"github.com/Yifangmo/micro-shop-services/common"
	"github.com/Yifangmo/micro-shop-services/order/proto"
	"github.com/Yifangmo/micro-shop-web/order/global"
	"github.com/Yifangmo/micro-shop-web/order/models"
	"github.com/Yifangmo/micro-shop-web/order/utils"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	req := proto.OrderListRequest{}
	model := claims.(*models.JWTClaims)
	if model.AuthorityId == 1 {
		req.UserId = int32(userId.(uint))
	}
	pages, perNums := ctx.DefaultQuery("p", "0"), ctx.DefaultQuery("pnum", "0")
	pagesInt, _ := strconv.Atoi(pages)
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PageInfo = &common.PageInfo{
		PageNumber: uint32(pagesInt),
		PageSize:   uint32(perNumsInt),
	}

	resp, err := global.OrderSrvClient.OrderList(ctx.Request.Context(), &req)
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	resMap := gin.H{
		"total": resp.Total,
	}
	data := make([]gin.H, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, gin.H{
			"id":                item.Id,
			"user":              item.UserId,
			"status":            item.Status,
			"created_at":        item.CreateAt,
			"pay_type":          item.PayType,
			"amount":            item.Amount,
			"order_sn":          item.OrderSn,
			"consignee_name":    item.ConsigneeName,
			"consignee_mobile":  item.ConsigneeMobile,
			"consignee_address": item.ConsigneeAddress,
			"remark":            item.Remark,
		})
	}
	resMap["data"] = data
	ctx.JSON(http.StatusOK, resMap)
}

func New(ctx *gin.Context) {
	orderForm := models.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")
	resp, err := global.OrderSrvClient.CreateOrder(context.WithValue(ctx.Request.Context(), global.GinContext, ctx), &proto.OrderRequest{
		UserId:           int32(userId.(uint)),
		ConsigneeName:    orderForm.Name,
		ConsigneeMobile:  orderForm.Mobile,
		ConsigneeAddress: orderForm.Address,
		Remark:           orderForm.Post,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	url, err := global.AlipayClient.TradePagePay(
		alipay.TradePagePay{
			Trade: alipay.Trade{
				Subject:     "微商城订单-" + resp.OrderSn,
				OutTradeNo:  resp.OrderSn,
				ProductCode: "FAST_INSTANT_TRADE_PAY",
				TotalAmount: strconv.FormatFloat(resp.Amount, 'f', 2, 64),
				NotifyURL:   global.ServerConfig.Alipay.NotifyURL,
				ReturnURL:   global.ServerConfig.Alipay.ReturnURL,
			},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         resp.Id,
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	req := &proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.JWTClaims)
	if model.AuthorityId == 1 {
		req.UserId = int32(userId.(uint))
	}
	resp, err := global.OrderSrvClient.OrderDetail(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	resM := gin.H{
		"id":                resp.OrderInfo.Id,
		"user":              resp.OrderInfo.UserId,
		"status":            resp.OrderInfo.Status,
		"created_at":        resp.OrderInfo.CreateAt,
		"pay_type":          resp.OrderInfo.PayType,
		"amount":            resp.OrderInfo.Amount,
		"order_sn":          resp.OrderInfo.OrderSn,
		"consignee_name":    resp.OrderInfo.ConsigneeName,
		"consignee_mobile":  resp.OrderInfo.ConsigneeMobile,
		"consignee_address": resp.OrderInfo.ConsigneeAddress,
		"remark":            resp.OrderInfo.Remark,
	}

	goodsList := make([]gin.H, 0, len(resp.OrderGoods))
	for _, item := range resp.OrderGoods {
		goodsList = append(goodsList, gin.H{
			"id":     item.GoodsId,
			"name":   item.GoodsName,
			"image":  item.GoodsImage,
			"price":  item.GoodsPrice,
			"number": item.GoodsNumber,
		})
	}
	resM["goods"] = goodsList

	ctx.JSON(http.StatusOK, resM)
}
