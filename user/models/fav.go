package models

type UserFavForm struct {
	GoodsId    int32 `form:"goods" json:"goods" binding:"required"`
}