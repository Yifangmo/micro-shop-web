package initialize

import (
	"github.com/Yifangmo/micro-shop-web/order/global"
	"github.com/smartwalle/alipay/v3"
)

func InitAlipayClient() {
	var err error
	global.AlipayClient, err = alipay.New(global.ServerConfig.Alipay.AppID, global.ServerConfig.Alipay.PrivateKey, !global.IsDebug)
	if err != nil {
		panic(err)
	}
	err = global.AlipayClient.LoadAliPayPublicKey((global.ServerConfig.Alipay.AliPublicKey))
	if err != nil {
		panic(err)
	}
}
