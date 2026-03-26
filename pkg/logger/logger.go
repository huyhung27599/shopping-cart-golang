package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	TraceIDContextKey contextKey = "trace_id"
)

var Log *zerolog.Logger

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	IsDev string
}

func InitLogger(config LoggerConfig) {
	Log = NewLogger(config)
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)



	var writer io.Writer

	if config.IsDev == "development" {
		writer = PrettyJSONWriter{Writer: os.Stdout}
	} else {
		writer = &lumberjack.Logger{
			Filename: config.Filename,
			MaxSize: config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge: config.MaxAge,
			Compress: config.Compress,
		}
	}

	
	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, p, "", "  "); err != nil {
		return w.Writer.Write(p)
	}
	return w.Writer.Write(prettyJSON.Bytes())
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDContextKey).(string); ok {
		return traceID
	}
	return ""
}