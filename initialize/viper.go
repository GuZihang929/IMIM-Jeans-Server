package initialize

import (
	"IM-Server/config"
	"IM-Server/global"
	"fmt"
	"github.com/spf13/viper"
)

func InitViper(configFile string) {
	v := viper.New()
	configFilePath := fmt.Sprintf("./config/%s", configFile)
	v.SetConfigFile(configFilePath)
	c := config.Config{}
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = v.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("read config file to struct err: %s \n", err))
	}
	global.Config = &c
}
