package logger

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogglyLogWriter struct {
	HttpClient *http.Client
	Url        string
}

type LogglyLogWriterDependency struct {
	HttpClient    *http.Client
	BaseUrl       string
	CustomerToken string
	Tag           string
}

func NewLogglyLogWriter(d LogglyLogWriterDependency) *LogglyLogWriter {
	return &LogglyLogWriter{
		HttpClient: d.HttpClient,
		Url:        fmt.Sprintf("%s/%s/tag/%s", d.BaseUrl, d.CustomerToken, d.Tag),
	}
}

func NewLogglyZapCore(w *LogglyLogWriter) zapcore.Core {
	logPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.DebugLevel
	})
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(
			zap.NewProductionEncoderConfig(),
		),
		zapcore.AddSync(w),
		logPriority,
	)
}

func (j *LogglyLogWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest(http.MethodPost, j.Url, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := j.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return 0, fmt.Errorf("[LogglyLogWriter] error status %d", res.StatusCode)
	}

	return len(p), nil
}
