package initialize

import (
	"IM-Server/config/confDetail"
	"IM-Server/global"
	"fmt"
)

func InitEmail() *confDetail.Email {
	EmailConfig := &confDetail.Email{
		Host:             global.Config.Email.Host,
		Port:             global.Config.Email.Port,
		User:             global.Config.Email.User,
		Password:         global.Config.Email.Password,
		DefaultFromEmail: global.Config.Email.DefaultFromEmail,
		UseSSL:           global.Config.Email.UseSSL,
		UserTls:          global.Config.Email.UserTls,
	}
	// 打印 emailConfig 的信息
	fmt.Printf("Email Config: %+v\n", EmailConfig)

	return EmailConfig
}
