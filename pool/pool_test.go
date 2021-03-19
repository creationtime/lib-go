package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := New(10, time.Second)
	assert.NotNil(t, p)
}
