package compress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZlib(t *testing.T) {
	data := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	z := &Zlib{}
	d1 := z.En(data)
	d2 := z.De(d1)
	assert.NotEqual(t, data, d1)
	assert.Equal(t, data, d2)
}
