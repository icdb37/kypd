// Package compress 数据压缩
package compress

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/utils"
	"github.com/icdb37/kypd/utils/logx"
)

// Compressor 压缩/解压
type Compressor interface {
	En([]byte) []byte
	De([]byte) []byte
}

// Zlib 加密处理对象
type Zlib struct{}

// NewZlib 创建加密处理对象
func NewZlib() *Zlib {
	return &Zlib{}
}

// Code 处理流程编号
func (z *Zlib) Code() byte {
	return enum.Compress
}

// GetHpn 头部编码
func (z *Zlib) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Compress<<enum.Bit8 + enum.Fixed,
		Data: []byte{},
	}
}

// SetHpn 头部编码
func (z *Zlib) SetHpn(h *utils.Hpn) {
}

// En 压缩数据
func (z *Zlib) En(data []byte) []byte {
	zr := bytes.NewBuffer(nil)
	zw := zlib.NewWriter(zr)
	if _, err := zw.Write(data); err != nil {
		return data
	}
	if err := zw.Flush(); err != nil {
		return data
	}
	data = zr.Bytes()
	logx.Debugf("--compress:en-- data: %v", data)
	return data
}

// De 解压数据
func (z *Zlib) De(data []byte) []byte {
	logx.Debugf("--compress:de-- data: %v", data)
	zr, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return data
	}
	w := bytes.NewBuffer(nil)
	if _, err := io.Copy(w, zr); err != nil && err.Error() != "unexpected EOF" {
		return data
	}
	return w.Bytes()
}
