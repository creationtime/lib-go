package jaeger

import "testing"

func TestNewTracer(t *testing.T) {
	if _, _, err := NewTracer("test", "10.4.61.88:30568"); err != nil {
		t.Error(err)
	}
}
