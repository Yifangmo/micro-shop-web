package message

import (
	"context"
	"net/http"

	"github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/global"
	"github.com/Yifangmo/micro-shop-web/user/models"
	"github.com/Yifangmo/micro-shop-web/user/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	request := &proto.UserMessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.JWTClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.UserSrvClient.GetUserMessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["message"] = value.Message
		reMap["file"] = value.File

		result = append(result, reMap)
	}
	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := models.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.UserSrvClient.CreateUserMessage(context.Background(), &proto.UserMessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
