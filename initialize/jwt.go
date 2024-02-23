package initialize

import (
	"IM-Server/config/confDetail"
	"IM-Server/global"
)

func InitJWT() *confDetail.Jwt {
	return &confDetail.Jwt{
		Secret:  global.Config.Jwt.Secret,
		Expires: global.Config.Jwt.Expires,
		Issuer:  global.Config.Jwt.Issuer,
	}

}
