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
	Name   string       `mapstructure:"name" json:"name"`
	Host   string       `mapstructure:"host" json:"host"`
	Tags   []string     `mapstructure:"tags" json:"tags"`
	Port   int          `mapstructure:"port" json:"port"`
	JWT    JWTConfig    `mapstructure:"jwt" json:"jwt"`
	Consul ConsulConfig `mapstructure:"consul" json:"consul"`
	Oss    OssConfig    `mapstructure:"oss" json:"oss"`
}

type JWTConfig struct {
	Secrect string `mapstructure:"secrect" json:"secrect"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type OssConfig struct {
	Key         string `mapstructure:"key" json:"key"`
	Secrect     string `mapstructure:"secrect" json:"secrect"`
	Host        string `mapstructure:"host" json:"host"`
	CallbackUrl string `mapstructure:"callback_url" json:"callback_url"`
	UploadDir   string `mapstructure:"upload_dir" json:"upload_dir"`
}
