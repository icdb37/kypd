// Package encrypt 数据加密
package encrypt

import (
	"fmt"
	"math/rand"

	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/utils"
	"github.com/icdb37/kypd/utils/logx"
)

// NewXorml 异或加密
func NewXorml(prevMask, password string) *Xorml {
	x := &Xorml{
		PrevMask: prevMask,
		Password: password,
	}
	x.Init()
	return x
}

// Xorml 多层异或
type Xorml struct {
	Nseq     uint8  // 干扰因子，数据层
	PrevMask string // 前一掩码
	CurrMask string // 当前掩码
	Password string // 用户密码
}

func (x *Xorml) Init() {
	x.CurrMask = utils.RandomPassword(32)
	x.Nseq = uint8(rand.Int())
}

// Code 处理流程编号
func (x *Xorml) Code() byte {
	return enum.Encrypt
}

// SetHpn 头部编码
func (x *Xorml) SetHpn(h *utils.Hpn) {
	x.Nseq = h.Data[0]
	x.CurrMask = string(h.Data[1:])
}

// GetHpn 头部编码
func (x *Xorml) GetHpn() *utils.Hpn {
	return &utils.Hpn{
		Code: enum.Encrypt<<enum.Bit8 + enum.Fixed,
		Data: append([]byte{x.Nseq}, []byte(x.CurrMask)...),
	}
}

// En 数据加密
func (x *Xorml) En(datas []byte) []byte {
	logx.Debugf("--encrypt:en-- nseq: %d, prev_mask: %s, curr_mask: %s, password: %s", x.Nseq, x.PrevMask, x.CurrMask, x.Password)
	utils.Xor3([]string{fmt.Sprint(len(datas)), x.PrevMask, x.CurrMask, x.Password}, datas, int(x.Nseq))
	return datas
}

// De 数据解密
func (x *Xorml) De(datas []byte) []byte {
	logx.Debugf("--encrypt:de-- nseq: %d, prev_mask: %s, curr_mask: %s, password: %s", x.Nseq, x.PrevMask, x.CurrMask, x.Password)
	utils.Xor3([]string{fmt.Sprint(len(datas)), x.PrevMask, x.CurrMask, x.Password}, datas, int(x.Nseq))
	return datas
}
