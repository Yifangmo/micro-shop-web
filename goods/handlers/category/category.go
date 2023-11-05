package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Yifangmo/micro-shop-services/goods/proto"
	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/goods/models"
	"github.com/Yifangmo/micro-shop-web/goods/utils"
)

func List(ctx *gin.Context) {
	r, err := global.GoodsSrvClient.AllCategoryList(ctx.Request.Context(), &empty.Empty{})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, r.JsonData)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	resp, err := global.GoodsSrvClient.GetSubCategory(ctx.Request.Context(), &proto.GetSubCategoryRequest{
		Id: int32(i),
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	resM := map[string]interface{}{
		"id": resp.Info.Id,
		"name": resp.Info.Name,
		"level": resp.Info.Level,
		"parent_category": resp.Info.ParentCategory,
		"is_tab": resp.Info.IsTab,
	}

	subCategories := make([]interface{}, 0)
	for _, sc := range resp.Subcategory {
		subCategories = append(subCategories, map[string]interface{}{
			"id":              sc.Id,
			"name":            sc.Name,
			"level":           sc.Level,
			"parent_category": sc.ParentCategory,
			"is_tab":          sc.IsTab,
		})
	}
	resM["sub_categories"] = subCategories
	ctx.JSON(http.StatusOK, resM)
}

func New(ctx *gin.Context) {
	categoryForm := models.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategory(ctx.Request.Context(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	ctx.JSON(http.StatusOK, request)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteCategory(ctx.Request.Context(), &proto.CategoryIDRequest{Id: int32(i)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	categoryForm := models.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		utils.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvClient.UpdateCategory(ctx.Request.Context(), request)
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
