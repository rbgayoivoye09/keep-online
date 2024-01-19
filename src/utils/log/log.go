package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var err error

func init() {
	// 配置日志文件的路径和其他相关参数
	logDirectory := "./logs/"
	logFile := logDirectory + "app.log"
	maxSize := 10 // MB
	maxBackups := 5
	maxAge := 7 // days

	// 创建日志目录
	if err := os.MkdirAll(logDirectory, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// 创建一个 lumberjack.Logger，用于处理日志轮换
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	// 创建 zap 的配置
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logFile}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 配置 zap logger
	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapcore.AddSync(lumberjackLogger)),
	)
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}

	defer func(Logger *zap.Logger) {
		err = Logger.Sync()
		if err != nil {
			// Handle the error, e.g., log it or take appropriate action
			Logger.Sugar().Error(err.Error())
		}
	}(Logger)
}
