package initialize

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/Yifangmo/micro-shop-web/goods/global"
)

func init() {
	flag.StringVar(&global.NacosConfigFileName, "c", "config", "Nacos配置文件")
	flag.Parse()
}

func InitConfig() {
	//从配置文件中读取出 Nacos 配置
	v := viper.New()
	v.SetConfigFile(global.NacosConfigFileName)
	if err := v.ReadInConfig(); err != nil {
		log.Panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		log.Panic(err)
	}
	log.Printf("Nacos server config: %#v", global.NacosConfig)

	//从 nacos 服务读取微服务配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              ".nacos-tmp/log",
		CacheDir:            ".nacos-tmp/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		log.Panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataID,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		log.Panicf("failed to unmarshal server config: %v", err)
	}
}
