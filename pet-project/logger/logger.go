package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"pet-project/settings"
	"os"
	
	// 添加 lumberjack 支持日志轮转
	"github.com/natefinch/lumberjack"
)

var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	var cfg zap.Config
	if settings.Conf.App.Debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// 设置日志级别
	level := zap.NewAtomicLevel()
	switch settings.Conf.App.LogLevel {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case "info":
		level.SetLevel(zap.InfoLevel)
	case "warn":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}
	cfg.Level = level

	// 自定义时间编码格式
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	// 根据环境设置日志输出
	if settings.Conf.App.Env == "production" {
		// 生产环境仅输出到文件
		fileEncoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
		
		// 创建日志目录
		if err := os.MkdirAll("logs", 0755); err != nil {
			panic(err)
		}
		
		// 使用 lumberjack 进行日志轮转，从配置文件读取参数，防止日志文件无限增长
		lumberjackLogger := &lumberjack.Logger{
			Filename:   settings.Conf.Log.Filename,     // 从配置文件读取日志文件路径
			MaxSize:    settings.Conf.Log.MaxSize,     // 从配置文件读取日志文件最大大小(MB)
			MaxBackups: settings.Conf.Log.MaxBackups,  // 从配置文件读取保留旧文件的最大个数
			MaxAge:     settings.Conf.Log.MaxAge,      // 从配置文件读取保留旧文件的最大天数
			Compress:   settings.Conf.Log.Compress,    // 从配置文件读取是否压缩旧文件
		}
		
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), level),
		)
		Logger = zap.New(core, zap.AddCaller())
	} else {
		// 开发环境仅输出到控制台
		var err error
		Logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
	}
}

// Sync 同步日志
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}