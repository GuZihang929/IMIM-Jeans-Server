package global

import (
	"IM-Server/config"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config //读取yaml文件信息
	Logger *zap.Logger    //全局日志
	DB     *gorm.DB
	Redis  *redis.Client
	Mongo  *mongo.Client
)
