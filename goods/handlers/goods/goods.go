package goods

import (
	"context"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Yifangmo/micro-shop-services/common"
	"github.com/Yifangmo/micro-shop-services/goods/proto"
	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/goods/models"
	"github.com/Yifangmo/micro-shop-web/goods/utils"
)

func List(ctx *gin.Context) {
	var form models.GoodsQueryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	req := &proto.GoodsListQueryRequest{}
	if form.PriceMin != nil {
		req.PriceMin = wrapperspb.Int32(*form.PriceMin)
	}
	if form.PriceMax != nil {
		req.PriceMax = wrapperspb.Int32(*form.PriceMax)
	}
	if form.IsHot != nil {
		req.IsHot = wrapperspb.Bool(*form.IsHot)
	}
	if form.IsNew != nil {
		req.IsNew = wrapperspb.Bool(*form.IsNew)
	}
	if form.IsTab != nil {
		req.IsTab = wrapperspb.Bool(*form.IsTab)
	}
	if form.TopCategory != nil {
		req.TopCategory = wrapperspb.Int32(*form.TopCategory)
	}
	if form.Keywords != nil {
		req.KeyWord = wrapperspb.String(*form.Keywords)
	}
	if form.BrandID != nil {
		req.Brand = wrapperspb.Int32(*form.BrandID)
	}
	req.PageInfo = &common.PageInfo{
		PageNumber: form.PageNumber,
		PageSize:   form.PageSize,
	}

	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}
	resp, err := global.GoodsSrvClient.GoodsListQuery(context.WithValue(ctx.Request.Context(), global.GinContext, ctx), req)
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()
	resM := map[string]interface{}{
		"total": resp.Total,
	}
	data := make([]interface{}, 0)
	for _, value := range resp.Data {
		data = append(data, map[string]interface{}{
			"id":               value.Id,
			"name":             value.Name,
			"brief":            value.Brief,
			"desc":             value.Desc,
			"is_free_shipping": value.IsFreeShipping,
			"images":           value.Images,
			"desc_images":      value.DescImages,
			"preview_image":    value.PreviewImage,
			"shop_price":       value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":     value.IsHot,
			"is_new":     value.IsNew,
			"is_on_sale": value.IsOnSale,
		})
	}
	resM["data"] = data

	ctx.JSON(http.StatusOK, resM)
}

func New(ctx *gin.Context) {
	form := models.GoodsForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	resp, err := goodsClient.CreateGoods(ctx.Request.Context(), &proto.GoodsInfoRequest{
		Name:           form.Name,
		Sn:             form.GoodsSn,
		Stock:          form.Stocks,
		MarketPrice:    form.MarketPrice,
		ShopPrice:      form.ShopPrice,
		Brief:          form.GoodsBrief,
		IsFreeShipping: *form.ShipFree,
		Images:         form.Images,
		DescImages:     form.DescImages,
		PreviewImage:   form.FrontImage,
		CategoryId:     form.CategoryId,
		BrandId:        form.Brand,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	resp, err := global.GoodsSrvClient.GetGoodsByID(context.WithValue(ctx.Request.Context(), global.GinContext, ctx), &proto.GoodsIDRequest{
		Id: int32(i),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":               resp.Id,
		"name":             resp.Name,
		"brief":            resp.Brief,
		"desc":             resp.Desc,
		"is_free_shipping": resp.IsFreeShipping,
		"images":           resp.Images,
		"desc_images":      resp.DescImages,
		"preview_image":    resp.PreviewImage,
		"shop_price":       resp.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   resp.Category.Id,
			"name": resp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   resp.Brand.Id,
			"name": resp.Brand.Name,
			"logo": resp.Brand.Logo,
		},
		"is_hot":     resp.IsHot,
		"is_new":     resp.IsNew,
		"is_on_sale": resp.IsOnSale,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(ctx.Request.Context(), &proto.GoodsIDRequest{Id: int32(i)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Status(http.StatusOK)
}

func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := models.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if _, err = global.GoodsSrvClient.UpdateGoods(ctx.Request.Context(), &proto.GoodsInfoRequest{
		Id:       int32(i),
		IsHot:    *goodsStatusForm.IsHot,
		IsNew:    *goodsStatusForm.IsNew,
		IsOnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(ctx *gin.Context) {
	goodsForm := models.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if _, err = global.GoodsSrvClient.UpdateGoods(ctx.Request.Context(), &proto.GoodsInfoRequest{
		Id:             int32(i),
		Name:           goodsForm.Name,
		Sn:             goodsForm.GoodsSn,
		Stock:          goodsForm.Stocks,
		MarketPrice:    goodsForm.MarketPrice,
		ShopPrice:      goodsForm.ShopPrice,
		Brief:          goodsForm.GoodsBrief,
		IsFreeShipping: *goodsForm.ShipFree,
		Images:         goodsForm.Images,
		DescImages:     goodsForm.DescImages,
		PreviewImage:   goodsForm.FrontImage,
		CategoryId:     goodsForm.CategoryId,
		BrandId:        goodsForm.Brand,
	}); err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
