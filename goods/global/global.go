package global

import (
	"github.com/Yifangmo/micro-shop-services/goods/proto"
	"github.com/Yifangmo/micro-shop-web/goods/configs"
	"github.com/spf13/viper"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfigFileName string
	IsDebug             bool
	Locale              string
	Translators         ut.Translator
	ServerConfig        configs.ServerConfig
	NacosConfig         configs.NacosConfig
	GoodsSrvClient      proto.GoodsClient
	GinContext          struct{}
)

func init() {
	viper.AutomaticEnv()
	IsDebug = viper.GetBool("MICRO_SHOP_DEBUG")
	Locale = "zh"
}
