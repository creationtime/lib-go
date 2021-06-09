package zerolog

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestName(t *testing.T) {
	Logger.Info().Msg("hello word")
	Logger.Info().Int("userId", 100).Msg("create user")
	Logger.Info().
		Dict("userInfo", zerolog.Dict().Str("username", "tom").Int("age", 10)).
		Float64("money", 11.11).
		Bool("exist", true).
		Msg("hello")
	Logger.Info().Caller(0).Send()
}
