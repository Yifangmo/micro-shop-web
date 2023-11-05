package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"

	"github.com/Yifangmo/micro-shop-web/oss/configs"
)

var (
	NacosConfigFileName string
	IsDebug             bool
	Locale              string
	Translators         ut.Translator
	ServerConfig        configs.ServerConfig
	NacosConfig         configs.NacosConfig
)

func init() {
	viper.AutomaticEnv()
	IsDebug = viper.GetBool("MICRO_SHOP_DEBUG")
	Locale = "zh"
}
