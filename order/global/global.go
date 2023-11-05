package global

import (
	"github.com/Yifangmo/micro-shop-services/order/proto"
	"github.com/Yifangmo/micro-shop-web/order/configs"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfigFileName string
	IsDebug             bool
	Locale              string
	GinContext          struct{}
	Translators         ut.Translator
	ServerConfig        configs.ServerConfig
	NacosConfig         configs.NacosConfig
	GoodsSrvClient      proto.GoodsClient
	OrderSrvClient      proto.OrderClient
	InventorySrvClient  proto.InventoryClient
	AlipayClient        *alipay.Client
)

func init() {
	viper.AutomaticEnv()
	IsDebug = viper.GetBool("MICRO_SHOP_DEBUG")
	Locale = "zh"
}
