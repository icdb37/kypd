package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	x := Xorml{
		PrevMask: "lNO@X-68(r%/7_rq;`?W,&}[w>5@`/CL",
		CurrMask: ")S?d~[9aq\"^f8G6!(Fc\\UDh+s}=?SkjZ",
		Password: "00000000",
	}
	// x.Init()
	data := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	d1 := make([]byte, len(data))
	copy(d1, data)
	x.En(d1)
	x.De(d1)
	assert.Equal(t, d1, data)
}
