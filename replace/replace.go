// Package replace 字节替换
package replace

import (
	"math/rand"

	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/utils"
	"github.com/icdb37/kypd/utils/logx"
)

// Fixed 固定间隔替换
type Fixed struct {
	interval byte
}

// NewFixed 固定间隔替换
func NewFixed() *Fixed {
	f := &Fixed{
		interval: byte(rand.Int()) / 2,
	}
	if f.interval == 0 {
		f.interval = 1
	}
	return f
}

// Code 处理流程编号
func (f *Fixed) Code() byte {
	return enum.Replace
}

// GetHpn 头部编码
func (f *Fixed) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Replace<<enum.Bit8 + enum.Fixed,
		Data: []byte{f.interval},
	}
}

// SetHpn 头部编码
func (f *Fixed) SetHpn(h *utils.Hpn) {
	if len(h.Data) != 1 {
		return
	}
	f.interval = h.Data[0]
}

// En 固定间隔数据替换
func (f *Fixed) En(datas []byte) []byte {
	logx.Debugf("[beg] --replace:en-- data: %v", datas)
	for i, size := 0, len(datas); i < size; i++ {
		datas[i] = datas[i] + f.interval
	}
	logx.Debugf("[end] --replace:en-- interval: %d, data: %v", f.interval, datas)
	return datas
}

// De 固定间隔数据恢复
func (f *Fixed) De(datas []byte) []byte {
	logx.Debugf("[beg] --replace:de-- data: %v", datas)
	for i, size := 0, len(datas); i < size; i++ {
		datas[i] = datas[i] - f.interval
	}
	logx.Debugf("[end] --replace:de-- interval: %d, data: %v", f.interval, datas)
	return datas
}

// Accum 累加间隔替换
type Accum struct {
	interval byte
}

// NewAccum 创建累加间隔替换
func NewAccum() *Accum {
	a := &Accum{
		interval: byte(rand.Int()) / 2,
	}
	if a.interval == 0 {
		a.interval = 1
	}
	return a
}

// GetHpn 头部编码
func (a *Accum) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Replace<<enum.Bit8 + enum.Accum,
		Data: []byte{a.interval},
	}
}

// SetHpn 头部编码
func (a *Accum) SetHpn(h *utils.Hpn) {
	if len(h.Data) != 1 {
		return
	}
	a.interval = h.Data[0]
}

// En 累加间隔数据替换
func (a *Accum) En(datas []byte) []byte {
	for i, size := 0, len(datas); i < size; i++ {
		datas[i] = datas[i] + a.interval
		a.interval = a.interval + (datas[i] & 0x04)
	}
	return datas
}

// De 累加间隔数据恢复
func (a *Accum) De(datas []byte) []byte {
	for i, size := 0, len(datas); i < size; i++ {
		interval := a.interval
		a.interval = a.interval + (datas[i] & 0x04)
		datas[i] = datas[i] - interval
	}
	return datas
}
