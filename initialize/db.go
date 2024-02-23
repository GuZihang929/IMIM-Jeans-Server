package initialize

import (
	"IM-Server/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitMysql() {
	dsn := global.Config.Mysql.Dsn()
	fmt.Println(dsn)
	var mysqlLogger logger.Interface
	if global.Config.Web.Env == "debug" {
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) //仅仅打印错误日志
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		global.Logger.Error("redis client error.")
		log.Fatalf("%s mysql client error.", dsn)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)  //最大连接数
	sqlDB.SetMaxOpenConns(100) //
	global.DB = db
	fmt.Println("MySQL连接成功")
}
