package logger

import (
	"net/http"
	"os"

	"github.com/anditakaesar/uwa-server-checker/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const NoticeLevel zapcore.Level = zapcore.DebugLevel - 1

var logInstance *Logger

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

func GetLogInstance() Interface {
	if logInstance != nil {
		return logInstance
	}

	newClient := http.Client{}
	logglyCore := NewLogglyZapCore(NewLogglyLogWriter(
		LogglyLogWriterDependency{
			HttpClient:    &newClient,
			BaseUrl:       env.LogglyBaseUrl(),
			CustomerToken: env.LogglyToken(),
			Tag:           env.LogglyTag(),
		},
	))

	return BuildNewLogger(logglyCore)
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

	return &Logger{
		zap: newZap,
	}
}
