package zerolog

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

var mainLogger *zerolog.Logger

type FormatLog struct {
	l *zerolog.Logger
}

func NewFormat(l *zerolog.Logger) *FormatLog {
	return &FormatLog{
		l: l,
	}
}

func init() {
	level, _ := logger.GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if level == logger.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	mainLogger = &l
}

func Log() *FormatLog {
	return NewFormat(mainLogger)
}

func (f *FormatLog) Debugf(format string, args ...interface{}) {
	f.l.Debug().Msgf(format, args...)
}

func (f *FormatLog) Debug(args interface{}) {
	switch args.(type) {
	case string:
		f.l.Debug().Msg(args.(string))
	case error:
		f.l.Debug().Msg(args.(error).Error())
	default:
		f.l.Debug().Msg(fmt.Sprintf("%v", args))
	}
}

func (f *FormatLog) Infof(format string, args ...interface{}) {
	f.l.Info().Msgf(format, args...)
}

func (f *FormatLog) Info(args interface{}) {
	switch args.(type) {
	case string:
		f.l.Info().Msg(args.(string))
	case error:
		f.l.Info().Msg(args.(error).Error())
	default:
		f.l.Info().Msg(fmt.Sprintf("%v", args))
	}
}

func (f *FormatLog) Warnf(format string, args ...interface{}) {
	f.l.Warn().Msgf(format, args...)
}

func (f *FormatLog) Warn(args interface{}) {
	switch args.(type) {
	case string:
		f.l.Warn().Msg(args.(string))
	case error:
		f.l.Warn().Msg(args.(error).Error())
	default:
		f.l.Warn().Msg(fmt.Sprintf("%v", args))
	}
}

func (f *FormatLog) Errorf(format string, args ...interface{}) {
	f.l.Error().Msgf(format, args...)
}

func (f *FormatLog) Error(err interface{}) {
	switch err.(type) {
	case string:
		f.l.Error().Msg(err.(string))
	case error:
		f.l.Err(err.(error)).Send()
	default:
		f.l.Error().Msg(fmt.Sprintf("%v", err))
	}

}

func (f *FormatLog) Fatalf(format string, args ...interface{}) {
	f.l.Fatal().Msgf(format, args...)
}

func (f *FormatLog) Fatal(args interface{}) {

	switch args.(type) {
	case string:
		f.l.Fatal().Msg(args.(string))
	case error:
		f.l.Fatal().Msg(args.(error).Error())
	default:
		f.l.Fatal().Msg(fmt.Sprintf("%v", args))
	}
}
