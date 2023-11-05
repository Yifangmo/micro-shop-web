package configs

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataID    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type ServerConfig struct {
	Name         string           `mapstructure:"name" json:"name"`
	Host         string           `mapstructure:"host" json:"host"`
	Tags         []string         `mapstructure:"tags" json:"tags"`
	Port         int              `mapstructure:"port" json:"port"`
	Consul       ConsulConfig     `mapstructure:"consul" json:"consul"`
	GoodsRPC     RPCServiceConfig `mapstructure:"goods_srv" json:"goods_srv"`
	OrderRPC     RPCServiceConfig `mapstructure:"order_srv" json:"order_srv"`
	InventoryRPC RPCServiceConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	Alipay       AlipayConfig     `mapstructure:"alipay" json:"alipay"`
	Jaeger       JaegerConfig     `mapstructure:"jaeger" json:"jaeger"`
	JWT          JWTConfig        `mapstructure:"jwt" json:"jwt"`
}

type RPCServiceConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	Secrect string `mapstructure:"secrect" json:"secrect"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
