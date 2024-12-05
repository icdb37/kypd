package utils

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"unsafe"
)

// CharPassword 密码字符串
const CharPassword = "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~`!@#$%^&*()_+-=[]{}|;':\",./<>?\\"

// Copy 切片拷贝
func Copy[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// Integer 整数类型
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// GetLittleEndian 数据小端字节序
func GetLittleEndian[T Integer](t T) []byte {
	size := unsafe.Sizeof(t)
	var data []byte
	switch size {
	case 1, 2:
		data = make([]byte, 2)
		binary.LittleEndian.PutUint16(data, uint16(t))
	case 4:
		data = make([]byte, 4)
		binary.LittleEndian.PutUint32(data, uint32(t))
	case 8:
		data = make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(t))
	}
	return data
}

// RandomPassword 随机密码
func RandomPassword(size int) string {
	p := make([]byte, size)
	for i := 0; i < size; i++ {
		p[i] = CharPassword[rand.Intn(len(CharPassword))]
	}
	return string(p)
}

// CalcCbit 计算校验和
func CalcCbit(data []byte, cbit uint16) uint16 {
	for i := 1; i < len(data); i = i + 2 {
		cbit = cbit ^ (uint16(data[i-1])<<8 + uint16(data[i]))
	}
	return cbit
}

// Hpn 头部处理节点
type Hpn struct {
	Code uint16
	Size uint16
	Data []byte
}

func (h *Hpn) En() []byte {
	h.Size = uint16(len(h.Data))
	data := make([]byte, 4+len(h.Data))
	binary.LittleEndian.PutUint32(data[:4], uint32(h.Code)<<16+uint32(h.Size))
	copy(data[4:], h.Data)
	return data
}

func (h *Hpn) De(data []byte) uint16 {
	h.Code, h.Size = 0, 0
	if len(data) < 4 {
		return 0
	}
	cs := binary.LittleEndian.Uint32(data[:4])
	h.Code = uint16(cs >> 16)
	h.Size = uint16(cs)
	h.Data = data[4 : 4+h.Size]
	return h.Size + 4
}

// Hpns 头部处理节点列表
type Hpns struct {
	Items []*Hpn // 数据处理流程
}

// NewHpns 创建头部处理节点列表
func NewHpns(items ...*Hpn) *Hpns {
	return &Hpns{
		Items: items,
	}
}

// En 头部节点加密
func (h *Hpns) En() []byte {
	buf := bytes.NewBuffer(nil)
	for _, i := range h.Items {
		buf.Write(i.En())
	}
	return buf.Bytes()
}

// De 头部节点解密
func (h *Hpns) De(data []byte) {
	h.Items = []*Hpn{}
	for len(data) > 4 {
		i := &Hpn{}
		step := i.De(data)
		if step == 0 {
			break
		}
		data = data[step:]
		h.Items = append(h.Items, i)
	}
}
