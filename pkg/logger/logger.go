package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
)

type Logger interface {
	Debug(ctx context.Context, message string, args ...any)
	Info(ctx context.Context, message string, args ...any)
	Warn(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, message string, args ...any)
	Fatal(ctx context.Context, message string, args ...any)

	LogMode(level gormLogger.LogLevel) gormLogger.Interface
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type logImpl struct {
	logger *logrus.Logger
}

var _ Logger = (*logImpl)(nil)

type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
)

func New(level string) {
	var l logrus.Level

	switch LogLevel(strings.ToLower(level)) {
	case ErrorLevel:
		l = logrus.ErrorLevel
	case WarnLevel:
		l = logrus.WarnLevel
	case InfoLevel:
		l = logrus.InfoLevel
	case DebugLevel:
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	//logrus.SetLevel(l)

	//skipFrameCount := 1
	logger := logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     l,
	}

	log = &logImpl{&logger}
}

// func (l *logImpl) prepare(ctx context.Context) *logrus.Logger {
// 	// reqID, ok := ctx.Value("req-id").(string)
// 	// if !ok || reqID == "" {
// 	// 	reqID = "not-provided"
// 	// }

// 	// logger := l.logger.With().
// 	// 	Str("req-id", reqID).
// 	// 	Logger()

// 	return l.logger
// }

func (l *logImpl) Debug(ctx context.Context, message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

func (l *logImpl) Info(ctx context.Context, message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *logImpl) Warn(ctx context.Context, message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}

func (l *logImpl) Error(ctx context.Context, message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *logImpl) Fatal(ctx context.Context, message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)

	os.Exit(1)
}

func (l *logImpl) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	var logLevel LogLevel

	switch level {
	case gormLogger.Silent:
		logLevel = DebugLevel
	case gormLogger.Error:
		logLevel = ErrorLevel
	case gormLogger.Warn:
		logLevel = WarnLevel
	case gormLogger.Info:
		logLevel = InfoLevel
	}

	New(string(logLevel))
	return log
}

func (l *logImpl) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	event := l.logger.WithFields(logrus.Fields{
		"elapsed": elapsed,
		"rows":    rows,
		"sql":     sql,
	})

	if err != nil {
		event.Error(err)
	}
}

var log *logImpl

// Log Get logger
func Log() Logger {
	if log == nil {
		New("")
	}
	return log
}