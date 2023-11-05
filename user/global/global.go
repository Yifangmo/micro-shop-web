package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"

	goodsProto "github.com/Yifangmo/micro-shop-services/goods/proto"
	userProto "github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/configs"
)

var (
	NacosConfigFileName string
	IsDebug             bool
	Locale              string
	Translators         ut.Translator
	ServerConfig        configs.ServerConfig
	NacosConfig         configs.NacosConfig
	UserSrvClient       userProto.UserClient
	GoodsSrvClient      goodsProto.GoodsClient
)

func init() {
	viper.AutomaticEnv()
	IsDebug = viper.GetBool("MICRO_SHOP_DEBUG")
	Locale = "zh"
}
