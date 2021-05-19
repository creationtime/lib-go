package zerolog

import (
	"os"

	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	level, _ := logger.GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if level == logger.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	//zerolog.TimestampFieldName = "t"
	//zerolog.LevelFieldName = "l"
	//zerolog.MessageFieldName = "m"
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}
