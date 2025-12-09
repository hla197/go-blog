package logger

import (
	"io"
)

var Log Logger

// Logger 是一个通用日志接口，类似于 fmt 接口风格

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	// 也可以包含一些对象打印
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Close()

	// 将 Gin 框架的输出重定向到 Zap
	GetIoWriter() io.Writer

	// 初始化
	Init()
}

func InitLogger() {
	// 使用Zap日志
	Log = &ZapLogger{}
	Log.Init()
}
