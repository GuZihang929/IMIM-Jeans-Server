package initialize

import (
	"IM-Server/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func InitRedis() {
	// 创建基础Redis客户端
	baseRedis := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr(),   // Redis服务器的地址和端口
		Password: global.Config.Redis.Password, // 密码（如果有的话）
		DB:       0,                            // 使用的数据库，默认为0
	})
	//测试连接
	_, err := baseRedis.Ping(context.Background()).Result()
	fmt.Println(err)
	if err != nil {
		global.Logger.Error("redis client error.")
		log.Fatalf("%s redis client error.", global.Config.Redis.Addr())
		return
	}
	global.Redis = baseRedis
	fmt.Println("Redis连接成功")
}
