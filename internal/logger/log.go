package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const NoticeLevel zapcore.Level = zapcore.DebugLevel - 1

type Interface interface {
	Info(string, ...zapcore.Field)
	Error(string, error, ...zapcore.Field)
	Warning(string, ...zapcore.Field)
	Flush()
}

type LoggerDependency struct {
	Zap *zap.Logger
}

type Logger struct {
	zap *zap.Logger
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	go l.zap.Info(msg, fields...)
}

func (l *Logger) Error(msg string, err error, fields ...zapcore.Field) {
	fields = append(fields, zap.String("internalError", err.Error()))
	go l.zap.Error(msg, fields...)
}

func (l *Logger) Warning(msg string, fields ...zapcore.Field) {
	go l.zap.Warn(msg, fields...)
}

func (l *Logger) Flush() {
	go l.zap.Sync()
}

func NewLogger(ld *LoggerDependency) Interface {
	defer ld.Zap.Sync()

	return &Logger{
		zap: ld.Zap,
	}
}

func BuildNewLogger(cores ...zapcore.Core) Interface {
	logPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.DebugLevel
	})

	cores = append(cores,
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(
				zap.NewDevelopmentEncoderConfig(),
			), zapcore.Lock(os.Stderr), logPriority))
	core := zapcore.NewTee(cores...)

	newZap := zap.New(core, zap.AddCallerSkip(1))

	defer newZap.Sync()
	return NewLogger(&LoggerDependency{
		Zap: newZap,
	})
}
