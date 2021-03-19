package logger

import (
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/logger/zerolog/v2"
	"os"
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
