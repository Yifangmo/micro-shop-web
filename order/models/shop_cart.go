package models

type ShoppingCartItemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Number  int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartItemUpdateForm struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}
