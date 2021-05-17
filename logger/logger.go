package logger

import (
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/logger/zerolog/v2"
	"os"
)

var (
	l Logger
)

type Logger logger.Logger

func New(opts ...logger.Option) Logger {
	level, _ := logger.GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	opts = append(opts, logger.WithLevel(level), zerolog.UseAsDefault())

	if level == logger.DebugLevel {
		opts = append(opts, zerolog.WithDevelopmentMode())
	}

	return logger.NewHelper(zerolog.NewLogger(opts...))
}

func Debug(args ...interface{}) {
	l.Log(logger.DebugLevel, args)
}

func Debugf(format string, args ...interface{}) {
	l.Logf(logger.DebugLevel, format, args)
}

func Info(args ...interface{}) {
	l.Log(logger.InfoLevel, args)
}

func Infof(format string, args ...interface{}) {
	l.Logf(logger.InfoLevel, format, args)
}

func Warn(args ...interface{}) {
	l.Log(logger.WarnLevel, args)
}

func Warnf(format string, args ...interface{}) {
	l.Logf(logger.WarnLevel, format, args)
}

func Error(args ...interface{}) {
	l.Log(logger.ErrorLevel, args)
}

func Errorf(format string, args ...interface{}) {
	l.Logf(logger.ErrorLevel, format, args)
}

func Fatal(args ...interface{}) {
	l.Log(logger.FatalLevel, args)
}

func Fatalf(format string, args ...interface{}) {
	l.Logf(logger.FatalLevel, format, args)
}
