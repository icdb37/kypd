package replace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFI(t *testing.T) {
	except := make([]byte, 256)
	for i := 0; i < 256; i++ {
		except[i] = byte(i)
	}
	actual := make([]byte, len(except))
	copy(actual, except)
	fi := Fixed{
		interval: 10,
	}
	fi.En(actual)
	fi.De(actual)
	assert.Equal(t, actual, except)
	fi.interval = 127
	fi.En(actual)
	fi.De(actual)
	assert.Equal(t, actual, except)
}

func TestAI(t *testing.T) {
	except := make([]byte, 256)
	for i := 0; i < 256; i++ {
		except[i] = byte(i)
	}
	actual := make([]byte, len(except))
	copy(actual, except)
	ai := Accum{
		interval: 10,
	}
	ai.En(actual)
	ai.interval = 10
	ai.De(actual)
	assert.Equal(t, actual, except)
	ai.interval = 127
	ai.En(actual)
	ai.interval = 127
	ai.De(actual)
	assert.Equal(t, actual, except)
}
