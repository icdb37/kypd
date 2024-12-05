package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswor(t *testing.T) {
	p := RandomPassword(32)
	assert.Equal(t, 32, len(p))
}
