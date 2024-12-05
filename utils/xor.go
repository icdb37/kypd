// Package utils 工具集合
package utils

import (
	"io"

	"github.com/icdb37/kypd/utils/logx"
)

// Xor0 异或
func Xor0(ks, data []byte, j uint8) uint8 {
	ks = Copy(ks)
	dSize := len(data)
	kSize := uint8(len(ks))
	for i := 0; i < dSize; i, j = i+1, j+1 {
		data[i] = data[i] ^ ks[j%kSize]
	}
	return j
}

// Xor1 异或
func Xor1(ks, data []byte) {
	ks = Copy(ks)
	dSize := len(data)
	kSize := len(ks)
	for i, j := 0, 0; i < dSize; i, j = i+1, j+1 {
		data[i] = data[i] ^ ks[j%kSize]
	}
}

// Xor2 异或
func Xor2(ks, data []byte, nseq int) {
	ks = Copy(ks)
	dSize := len(data)
	kSize := len(ks)
	skipSize := nseq % kSize
	if skipSize == 0 {
		skipSize = kSize - 1
	}
	skipIndex := 0
	useKeys := make([]byte, 0)
	for i, j := 0, 1; i < dSize; i = i + 1 {
		q := ks[j%kSize]
		useKeys = append(useKeys, q)
		data[i] = data[i] ^ q
		j += skipIndex + (i / kSize) + int(q%7) + 1
		skipIndex++
		if skipIndex > skipSize {
			skipIndex = 0
		}
	}
	logx.Debug(string(useKeys))
}

// Xor3 异或
func Xor3(ks []string, data []byte, nseq int) {
	for _, k := range ks {
		Xor1([]byte(k), data)
		tk := []byte{}
		for _, x := range []int{3, 5, 7} {
			for j, size := x, len(k); j < size; j = j + x {
				tk = append(tk, k[j])
			}
		}
		if len(tk) == 0 {
			continue
		}
		Xor1(tk, data)
	}
	for _, k := range ks {
		Xor2([]byte(k), data, nseq)
	}
}

// NewXorReader 创建读异或
func NewXorReader(r io.ReadCloser, ks []byte) *XorReader {
	x := &XorReader{
		r:  r,
		ks: ks,
	}
	return x
}

// XorReader 读流异或
type XorReader struct {
	i  int
	ks []byte
	r  io.ReadCloser
}

// Read 读数据
func (x *XorReader) Read(p []byte) (int, error) {
	n, err := x.r.Read(p)
	if n > 0 {
		for i := 0; i < n; i++ {
			p[i] ^= x.ks[x.i]
			x.i++
			x.i = x.i % len(x.ks)
		}
	}
	return n, err
}

// Close 关闭
func (x *XorReader) Close() error {
	return x.r.Close()
}

// NewXorWriter 创建写异或
func NewXorWriter(w io.WriteCloser, ks []byte) *XorWriter {
	x := &XorWriter{
		w:  w,
		ks: ks,
	}
	return x
}

// XorWriter 写流异或
type XorWriter struct {
	i  int
	ks []byte
	w  io.WriteCloser
}

// Write 写数据
func (x *XorWriter) Write(p []byte) (int, error) {
	size := len(p)
	for i := 0; i < size; i++ {
		p[i] ^= x.ks[x.i]
		x.i++
		x.i = x.i % len(x.ks)
	}
	return x.w.Write(p)
}

// Close 关闭
func (x *XorWriter) Close() error {
	return x.w.Close()
}
