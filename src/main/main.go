package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rbgayoivoye09/keep-online/src/utils/cmd"
	"github.com/rbgayoivoye09/keep-online/src/utils/log"
)

func main() {

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
		Compress:   true,
	}

	pec := zap.NewProductionEncoderConfig()
	pec.EncodeTime = zapcore.ISO8601TimeEncoder
	// Create a zap core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(pec), // Use JSON format for structured logging
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger)),
		zap.InfoLevel, // Log level
	)

	// Create a logger with the core
	log.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// Use the logger
	defer func() {
		if err := log.Logger.Sync(); err != nil {
			log.Logger.Sugar().Info("Failed to sync logger:", err)
		}
	}()

	if err := cmd.TrootCmd.Execute(); err != nil {
		log.Logger.Sugar().Error(err)
		os.Exit(1)
	}
}
