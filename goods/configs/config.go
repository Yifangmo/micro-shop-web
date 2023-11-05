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

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
