package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger 声明
var (
	logger *zap.Logger
	err    error
)

// Init 初始化 logger
func Init() {
	// 定义解码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 输出日志时间的key名
		LevelKey:       "level",                        // 输出日志级别的key名
		NameKey:        "logger",                       // 输出日志名称的key名
		CallerKey:      "caller",                       // 输出日志调用方的key名
		MessageKey:     "message",                      // 输出日志信息的key名
		StacktraceKey:  "stacktrace",                   // 输出日志堆栈的key名
		LineEnding:     zapcore.DefaultLineEnding,      // 每行的分隔符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行消耗的时间转化成浮点型的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.InfoLevel)

	// 定义日志配置
	config := zap.Config{
		Level:             atom,               // 日志级别
		Development:       true,               // 开发模式，堆栈跟踪
		Encoding:          "json",             // 输出格式 console 或 json
		DisableStacktrace: true,               // 禁用堆栈
		EncoderConfig:     encoderConfig,      // 编码器配置
		OutputPaths:       []string{"stdout"}, // stdout（标准输出，正常颜色）
		ErrorOutputPaths:  []string{"stderr"}, // stderr（错误输出，红色）
	}

	// 构建日志
	logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("logger initialize failed: %v", err))
	}
	logger.Info("logger initialize succeeded.")
}

// Logger 定义日志调用函数
func Logger() *zap.Logger {
	return logger
}
