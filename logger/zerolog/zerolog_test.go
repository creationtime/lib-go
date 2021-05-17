package zerolog

import (
	"testing"
)

func TestName(t *testing.T) {
	Log().Infof("username:%s, userId:%d", "test", 1)
}
