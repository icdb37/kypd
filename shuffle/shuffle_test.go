// Package shuffle 打乱顺序
package shuffle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixed0(t *testing.T) {
	// except := []byte{0, 1, 2, 126, 127, 128, 254, 255}
	except := make([]byte, 256)
	for i := 0; i < 256; i++ {
		except[i] = byte(i)
	}
	actual := make([]byte, len(except))
	copy(actual, except)
	fi := Fixed{
		offset:   2,
		interval: 1,
	}
	fi.En(actual)
	fi.De(actual)
	assert.Equal(t, actual, except)
}

func TestFixed1(t *testing.T) {
	except := []byte{127, 248, 32, 150, 60, 209, 140, 30, 88, 120, 39, 65, 150, 247}
	actual := make([]byte, len(except))
	copy(actual, except)
	fi := Fixed{
		offset:   195,
		interval: 1,
	}
	fi.En(actual)
	fi.De(actual)
	assert.Equal(t, actual, except)
}

func TestAccum(t *testing.T) {
	// except := []byte{0, 1, 2, 126, 127, 128, 253, 254, 255}
	except := make([]byte, 256)
	for i := 0; i < 256; i++ {
		except[i] = byte(i)
	}
	actual := make([]byte, len(except))
	copy(actual, except)
	fi := Accum{
		offset:   2,
		interval: 1,
	}
	fi.En(actual)
	fi.interval = 1
	fi.De(actual)
	assert.Equal(t, actual, except)
}
