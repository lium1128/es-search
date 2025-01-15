package common_logs

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lium-product/es-search/pkg/cfg"
)

type Logger struct {
	Logger *zap.Logger
	Err    error
}

// InitLogger
// 基础日志初始化逻辑
// @param cfg 为 cfg.LoadLogger() 或者 cfg.LoadMergeLogger() 或其他的日志配置对象，具体与实际业务层初始化为主
func InitLogger(cfg cfg.Logger) *Logger {
	encoder := getEncoder()
	l, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return &Logger{Err: err}
	}
	mCores := make([]zapcore.Core, 0)

	// 添加控制台输出
	if cfg.OutPutConsole {
		mCores = append(mCores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), l))
	}
	// 添加文件输出
	if cfg.OutPutFile {
		writeSyncer := getLogWriter(cfg.FileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
		core := zapcore.NewCore(encoder, writeSyncer, l)
		mCores = append(mCores, core)
	}

	// 生成实例
	logger := zap.New(zapcore.NewTee(mCores...), zap.AddCaller())
	return &Logger{Logger: logger}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder("2006-01-02 15:04:05.000")
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func customTimeEncoder(customFormat string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(customFormat))
	}
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   // 每个日志文件的最大大小（单位：MB）
		MaxBackups: maxBackup, // 保留的旧日志文件最大数量
		MaxAge:     maxAge,    // 保留的旧日志文件最大天数
		LocalTime:  true,      // 使用本地时间
		Compress:   true,      // 是否压缩/归档旧日志文件
	}
	fmt.Println("日志文件路径：", lumberJackLogger.Filename)
	return zapcore.AddSync(lumberJackLogger)
}

func (l *Logger) Debug(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Debug(args...)
}

func (l *Logger) Info(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Info(args...)
}

func (l *Logger) Warn(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Warn(args...)
}

func (l *Logger) Error(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Error(args...)
}

func (l *Logger) DPanic(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).DPanic(args...)
}

func (l *Logger) Panic(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Panic(args...)
}

func (l *Logger) Fatal(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Fatal(args...)
}

func (l *Logger) Debugf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Fatalf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Fatalw(msg, keysAndValues...)
}

func (l *Logger) Debugln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Debugln(args...)
}

func (l *Logger) Infoln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Infoln(args...)
}

func (l *Logger) Warnln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Warnln(args...)
}

func (l *Logger) Errorln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Error(args...)
}

func (l *Logger) DPanicln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).DPanicln(args...)
}

func (l *Logger) Panicln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Panicln(args...)
}

func (l *Logger) Fatalln(args ...any) {
	l.Logger.Sugar().WithOptions(zap.AddCallerSkip(1)).Fatalln(args...)
}
func (l *Logger) WithOptions(opts ...zap.Option) *zap.SugaredLogger {
	return l.Logger.WithOptions(opts...).Sugar()
}

func (l *Logger) WithOutCaller() *zap.SugaredLogger {
	return l.Logger.Sugar().WithOptions(zap.WithCaller(false))
}
