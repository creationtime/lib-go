package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	//l := New()
	//l.Logf(logger.DebugLevel, "debug log")
	//l.Logf(logger.InfoLevel, "info log")
	//l.Logf(logger.WarnLevel, "warning log")
	//l.Logf(logger.ErrorLevel, "error log")
	//l.Logf(logger.FatalLevel, "fatal log")

	Infof("user:%v,test:%s", 111, "aseraera")
	Warn("111")
}
