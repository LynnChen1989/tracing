package lib

import (
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"tracing/conf"
)

func InitLogger() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   conf.LogPath,       // 日志文件路径
		MaxSize:    conf.LogMaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.LogMaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     conf.LogMaxAge,     // 文件最多保存多少天
		Compress:   conf.LogCompress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	var level zapcore.Level
	switch conf.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", conf.SysName))
	// 构造日志
	logger := zap.New(core, caller, development, filed)
	return logger
}

func GetLogger() *zap.Logger {
	return InitLogger()
}
