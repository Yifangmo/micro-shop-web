package pay

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"

	"github.com/Yifangmo/micro-shop-services/order/proto"
	"github.com/Yifangmo/micro-shop-web/order/global"
)

// 支付宝交易通知的回调
func Notify(ctx *gin.Context) {
	tn, err := global.AlipayClient.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	_, err = global.OrderSrvClient.UpdateOrderStatus(ctx.Request.Context(), &proto.UpdateOrderStatusRequest{
		OrderSn: tn.OutTradeNo,
		Status:  TradeStatusToOrderStatus(tn.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.String(http.StatusOK, "success")
}

func TradeStatusToOrderStatus(ts alipay.TradeStatus) proto.OrderStatus {
	switch ts {
	case alipay.TradeStatusFinished:
		return proto.OrderStatus_ORDER_STATUS_CLOSED
	case alipay.TradeStatusSuccess:
		return proto.OrderStatus_ORDER_STATUS_SUCCESS
	}
	return proto.OrderStatus_ORDER_STATUS_WAITING_PAY
}
