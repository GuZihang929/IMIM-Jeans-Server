package initialize

import (
	"IM-Server/config/confDetail"
	"IM-Server/global"
)

func InitWeb() *confDetail.Web {
	return &confDetail.Web{
		Host: global.Config.Web.Host,
		Port: global.Config.Web.Port,
		Env:  global.Config.Web.Env,
	}
}
