package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 是 Logger 接口的具体实现
type ZapLogger struct {
	sugaredLogger *zap.SugaredLogger
	// 缓存 Writer，避免重复创建
	ioWriter io.Writer
}

// 实现接口方法
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

// 简单打印（不带格式化）
func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *ZapLogger) Close() {
	l.sugaredLogger.Sync()
}

func (l *ZapLogger) Init() {

	// 1. 配置编码器（这里使用 JSON 格式，生产环境推荐）
	encoder := getEncoder()

	// 2. 配置日志写入位置（控制台 + 文件）
	writeSyncer := getLogWriter()

	// 3. 配置日志级别（例如 DebugLevel 及以上）
	level := zap.NewAtomicLevelAt(zap.DebugLevel)

	// 4. 创建 core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 5. 创建 Logger
	logger := zap.New(core, zap.AddCaller()) // zap.AddCaller() 用于记录调用位置

	// 6. 使用 SugaredLogger（更易用的 API）
	l.sugaredLogger = logger.Sugar()

	stdLog := zap.NewStdLog(logger.WithOptions(zap.AddCallerSkip(1)))

	l.ioWriter = stdLog.Writer()

}

// 将 Gin 框架的输出重定向到 Zap
func (l *ZapLogger) GetIoWriter() io.Writer {
	return l.ioWriter
}

// JSON 编码配置
func getEncoder() zapcore.Encoder {
	// 日志输出为控制台格式
	// 自定义控制台编码配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
	}

	return zapcore.NewConsoleEncoder(encoderConfig)

	// 日志输出为json格式
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// 日志写入配置（同时写入文件和控制台）
func getLogWriter() zapcore.WriteSyncer {
	// 配置 lumberjack 进行日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/gin.log", // 日志文件路径
		MaxSize:    10,               // 每个文件最大 10MB
		MaxBackups: 5,                // 最多保留 5 个备份
		MaxAge:     30,               // 文件最多保存 30 天
		Compress:   true,             // 是否压缩
	}

	// 同时输出到文件和控制台
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
