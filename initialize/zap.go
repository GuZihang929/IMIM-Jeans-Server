package initialize

import (
	"IM-Server/global"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

func InitZap() {
	currentTime := time.Now().Format("2006-01-02")
	Filename := fmt.Sprintf("log/log_%s_%s.log", currentTime)
	logFile := &lumberjack.Logger{
		Filename:   Filename,
		MaxSize:    global.Config.Zap.MaxSize,
		MaxBackups: global.Config.Zap.MaxBackups,
		MaxAge:     global.Config.Zap.MaxAge,
		Compress:   global.Config.Zap.Compress,
	}
	// 配置日志级别
	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	// 配置日志编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 配置核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		logLevel,
	)
	// 创建Logger
	logger := zap.New(core, zap.AddCaller())
	//将log对象赋给全局log对象
	global.Logger = logger
}
