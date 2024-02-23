package confDetail

type Email struct {
	Host             string `json:"host" yaml:"host"`
	Port             int    `json:"port" yaml:"port"`
	User             string `json:"user_api" yaml:"user_api"` //发件人邮箱
	Password         string `json:"password" yaml:"password"`
	DefaultFromEmail string `json:"default_from_email" yaml:"default_from_email"` //默认发件人名字
	UseSSL           bool   `json:"use_ssl" yaml:"use_ssl"`                       //是否使用ssl
	UserTls          bool   `json:"user_tls" yaml:"user_tls"`                     //是否开启tls
}
