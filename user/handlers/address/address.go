package address

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/global"
	"github.com/Yifangmo/micro-shop-web/user/models"
	"github.com/Yifangmo/micro-shop-web/user/utils"
)

func List(ctx *gin.Context) {
	request := &proto.ConsigneeAddressRequest{}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.JWTClaims)

	if currentUser.AuthorityId != 2 {
		userId, _ := ctx.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.UserSrvClient.GetConsigneeAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["province"] = value.Province
		reMap["city"] = value.City
		reMap["district"] = value.District
		reMap["address"] = value.ConsigneeAddress
		reMap["signer_name"] = value.ConsigneeName
		reMap["signer_mobile"] = value.ConsigneeMobile

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	addressForm := models.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.UserSrvClient.CreateConsigneeAddress(context.Background(), &proto.ConsigneeAddressRequest{
		UserId:           int32(userId.(uint)),
		Province:         addressForm.Province,
		City:             addressForm.City,
		District:         addressForm.District,
		ConsigneeAddress: addressForm.Address,
		ConsigneeName:    addressForm.SignerName,
		ConsigneeMobile:  addressForm.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("新建地址失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.UserSrvClient.DeleteConsigneeAddress(context.Background(), &proto.ConsigneeAddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Update(ctx *gin.Context) {
	addressForm := models.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.UserSrvClient.UpdateConsigneeAddress(context.Background(), &proto.ConsigneeAddressRequest{
		Id:               int32(i),
		Province:         addressForm.Province,
		City:             addressForm.City,
		District:         addressForm.District,
		ConsigneeAddress: addressForm.Address,
		ConsigneeName:    addressForm.SignerName,
		ConsigneeMobile:  addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
