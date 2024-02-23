package config

import "IM-Server/config/confDetail"

type Config struct {
	Web   confDetail.Web   `json:"web" yaml:"web"`
	Ws    confDetail.Ws    `json:"ws" yaml:"ws"`
	Zap   confDetail.Zap   `json:"Logger" yaml:"zap"`
	Mysql confDetail.Mysql `json:"mysql" yaml:"mysql"`
	Redis confDetail.Redis `json:"redis" yaml:"redis"`
	Jwt   confDetail.Jwt   `json:"jwt"   yaml:"jwt"`
	Email confDetail.Email `json:"email" yaml:"email"`
}
