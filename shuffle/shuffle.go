// Package shuffle 打乱顺序
package shuffle

import (
	"math/rand"

	"github.com/icdb37/kypd/utils/logx"

	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/utils"
)

// Shuffler 乱序/恢复
type Shuffler interface {
	En([]byte) []byte
	De([]byte) []byte
}

// Fixed 固定间隔打乱
type Fixed struct {
	offset   byte
	interval byte
}

// NewFixed 乱序固定间隔
func NewFixed() *Fixed {
	f := &Fixed{}
	f.Init()
	return f
}

// Init 初始化
func (f *Fixed) Init() {
	f.offset = byte(rand.Int())
	f.interval = byte(rand.Int()) & 0x07
	if f.interval == 0 {
		f.interval = 1
	}
}

// Code 处理流程编号
func (f *Fixed) Code() byte {
	return enum.Shuffle
}

// SetHpn 头部编码
func (f *Fixed) SetHpn(h *utils.Hpn) {
	if len(h.Data) != 2 {
		return
	}
	f.offset = h.Data[0]
	f.interval = h.Data[1]
}

// GetHpn 头部编码
func (f *Fixed) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Shuffle<<enum.Bit8 + enum.Fixed,
		Data: []byte{f.offset, f.interval},
	}
}

// En 固定间隔数据替换
func (f *Fixed) En(datas []byte) []byte {
	logx.Debugf("[beg] --shuffle:en-- data: %v", datas)
	size := len(datas)
	offset := int(f.offset) % size
	i := 0
	for i = 0; i < offset; i++ {
		datas[i], datas[size-offset+i] = datas[size-offset+i], datas[i]
	}
	for i = 0; i < size; i = i + int(f.interval) + 1 {
		j := (i + int(f.interval)) % size
		datas[i], datas[j] = datas[j], datas[i]
	}
	logx.Debugf("[end] --shuffle:en-- interval: %d, offset: %d, i: %d, data: %v", f.interval, f.offset, i, datas)
	return datas
}

// De 固定间隔数据恢复
func (f *Fixed) De(datas []byte) []byte {
	logx.Debugf("[beg] --shuffle:de-- data: %v", datas)
	size := len(datas)
	i := size - (size % (int(f.interval) + 1))
	if i == size {
		i = size - (int(f.interval) + 1)
	}
	for ; 0 <= i; i = i - int(f.interval) - 1 {
		j := (i + int(f.interval)) % size
		datas[i], datas[j] = datas[j], datas[i]
	}
	offset := int(f.offset) % size
	for i = 1; i <= offset; i++ {
		datas[offset-i], datas[size-i] = datas[size-i], datas[offset-i]
	}
	logx.Debugf("[end] --shuffle:de-- interval: %d, offset: %d, i: %d, data: %v", f.interval, f.offset, i, datas)
	return datas
}

// Accum 累加间隔替换
type Accum struct {
	offset   byte
	interval byte
}

// Init 初始
func (a *Accum) Init() {
	a.offset = byte(rand.Int())
	a.interval = byte(rand.Int()) & 0x07
	if a.interval == 0 {
		a.interval = 1
	}
}

// GetHpn 头部编码
func (a *Accum) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Shuffle<<enum.Bit8 + enum.Fixed,
		Data: []byte{a.offset, a.interval},
	}
}

// SetHpn 头部编码
func (a *Accum) SetHpn(h *utils.Hpn) {
	if len(h.Data) != 2 {
		return
	}
	a.offset = h.Data[0]
	a.interval = h.Data[1]
}

// En 固定间隔数据替换
func (a *Accum) En(datas []byte) []byte {
	size := len(datas)
	offset := int(a.offset) % size
	i := 0
	for i = 0; i < offset; i++ {
		datas[i], datas[size-offset+i] = datas[size-offset+i], datas[i]
	}
	for i = 0; i < size; i = i + int(a.interval) {
		j := i + int(a.interval)
		if j >= size {
			break
		}
		datas[i], datas[j] = datas[j], datas[i]
		a.interval = datas[i]&0x07&a.interval + 1
		i = j // 注意：必须跳过当前交互区间，否则可能在区间内发生交换导致无法恢复
	}
	return datas
}

// De 固定间隔数据恢复
func (a *Accum) De(datas []byte) []byte {
	size := len(datas)
	i := 0
	for ; i < size; i = i + int(a.interval) {
		j := i + int(a.interval)
		if j >= size {
			break
		}
		interval := datas[i]&0x07&a.interval + 1
		datas[i], datas[j] = datas[j], datas[i]
		a.interval = interval
		i = j // 注意：必须跳过当前交互区间，否则可能在区间内发生交换导致无法恢复
	}
	offset := int(a.offset) % size
	for i = 0; i < offset; i++ {
		datas[i], datas[size-offset+i] = datas[size-offset+i], datas[i]
	}
	return datas
}
