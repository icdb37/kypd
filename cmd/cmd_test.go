package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	// 加密
	hw := &KypdHead{}
	hw.preEn("a", "a")
	w := bytes.NewBuffer(nil)
	assert.Nil(t, hw.En(w))
	// 解密
	r := bytes.NewBuffer(w.Bytes())
	hr := &KypdHead{
		password: "a",
	}
	assert.Nil(t, hr.De(r))
	assert.Nil(t, hr.pstDe())

	assert.Equal(t, hw.version, hr.version)
	assert.Equal(t, hw.hpns, hr.hpns)
}

func TestBody(t *testing.T) {
	h := &KypdHead{}
	h.preEn("a", "a")
	except := []byte("abcdefghijklmnopqrstuvwxyz")
	b := NewKypdBody(h)
	r1, w1 := bytes.NewBuffer(except), bytes.NewBuffer(nil)
	if err := b.En(r1, w1); err != nil {
		assert.Nil(t, err)
	}
	r2, w2 := w1, bytes.NewBuffer(nil)
	b = NewKypdBody(h)
	if err := b.De(r2, w2); err != nil {
		assert.Nil(t, err)
	}
	actual, err := io.ReadAll(w2)
	if err != nil {
		assert.Nil(t, err)
	}
	assert.Equal(t, string(except), string(actual))
}
