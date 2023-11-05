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
	Name        string           `mapstructure:"name" json:"name"`
	Host        string           `mapstructure:"host" json:"host"`
	Tags        []string         `mapstructure:"tags" json:"tags"`
	Port        int              `mapstructure:"port" json:"port"`
	UserService RPCServiceConfig `mapstructure:"user_service" json:"user_service"`
	JWT         JWTConfig        `mapstructure:"jwt" json:"jwt"`
	Sms         SmsConfig        `mapstructure:"sms" json:"sms"`
	Redis       RedisConfig      `mapstructure:"redis" json:"redis"`
	Consul      ConsulConfig     `mapstructure:"consul" json:"consul"`
}

type RPCServiceConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	Secrect string `mapstructure:"secrect" json:"secrect"`
}

type SmsConfig struct {
	Key     string `mapstructure:"key" json:"key"`
	Secrect string `mapstructure:"secrect" json:"secrect"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}
