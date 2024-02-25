package main

import (
	"IM-Server/core"
	"IM-Server/global"
	"IM-Server/im/browser"
	"IM-Server/im/conn"
	"IM-Server/initialize"
	"context"
	"fmt"
	"sync"
	"time"
)

const ConfigFile = "config.yaml"

func main() {
	//读取配置文件
	initialize.InitViper(ConfigFile)
	//初始化log
	initialize.InitZap()
	//初始化连接mysql
	initialize.InitMysql()
	//初始化连接redis
	initialize.InitRedis()
	//连接mongo
	//initialize.Global_Mongo = initialize.InitMongo()
	//加载配置文件中的jwt信息
	initialize.InitJWT()
	//加载配置文件中的邮箱信息
	initialize.InitEmail()
	//初始化路由
	Host := fmt.Sprintf(":%s", global.Config.Web.Host)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		core.RunWindowsServer(Host)
		wg.Done()
	}()
	//创建长连接
	var server conn.WsServer
	op := &conn.WsServerOptions{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	server = conn.NewWsServer(op)
	browser.SetMessageHandler(browser.Handel)
	cm := browser.NewDefaultManager()
	server.SetConnHandler(func(ctx context.Context, conn conn.Connection) {
		cm.BrowserConnected(ctx, conn)
	})
	browser.SetDefaultManager(cm)

	wg.Add(1)
	go func() {
		addr := global.Config.Ws.Addr
		port := global.Config.Ws.Port
		err := server.Run(addr, port)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
	//browser.SetInterfaceImpl(cm)

}
