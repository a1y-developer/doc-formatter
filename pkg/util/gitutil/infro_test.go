package gitutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet_ReturnsInfo(t *testing.T) {
	info := Get(".")

	assert.IsType(t, Info{}, info)
}
