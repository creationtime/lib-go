package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var l *zerolog.Logger

func New() *zerolog.Logger {
	level, _ := logger.GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if level == logger.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	nl := zerolog.New(os.Stderr).With().Timestamp().Stack().CallerWithSkipFrameCount(3).Logger()
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		return WashPath(file) + ":" + strconv.Itoa(line)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return &nl
}

// WashPath 路径脱敏
func WashPath(s string) string {
	sPath := strings.Split(s, "elysium")
	if len(sPath) >= 2 {
		return sPath[1]
	}

	return s
}

func init() {
	l = New()
}

func Debugf(format string, args ...interface{}) {
	l.Debug().Msgf(format, args...)
}

func Debug(args interface{}) {
	switch args.(type) {
	case string:
		l.Debug().Msg(args.(string))
	case error:
		l.Debug().Msg(args.(error).Error())
	default:
		l.Debug().Msg(fmt.Sprintf("%v", args))
	}
}

func Infof(format string, args ...interface{}) {
	l.Info().Msgf(format, args...)
}

func Info(args interface{}) {
	switch args.(type) {
	case string:
		l.Info().Msg(args.(string))
	case error:
		l.Info().Msg(args.(error).Error())
	default:
		l.Info().Msg(fmt.Sprintf("%v", args))
	}
}

func Warnf(format string, args ...interface{}) {
	l.Warn().Msgf(format, args...)
}

func Warn(args interface{}) {
	switch args.(type) {
	case string:
		l.Warn().Msg(args.(string))
	case error:
		l.Warn().Msg(args.(error).Error())
	default:
		l.Warn().Msg(fmt.Sprintf("%v", args))
	}
}

func Errorf(format string, args ...interface{}) {
	l.Error().Msgf(format, args...)
}

func Error(err interface{}) {
	switch err.(type) {
	case string:
		l.Error().Msg(err.(string))
	case error:
		l.Err(err.(error)).Send()
	default:
		l.Error().Msg(fmt.Sprintf("%v", err))
	}
}

func Fatalf(format string, args ...interface{}) {
	l.Fatal().Msgf(format, args...)
}

func Fatal(args interface{}) {
	switch args.(type) {
	case string:
		l.Fatal().Msg(args.(string))
	case error:
		l.Fatal().Msg(args.(error).Error())
	default:
		l.Fatal().Msg(fmt.Sprintf("%v", args))
	}
}
