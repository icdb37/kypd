package cmd

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/utils"
)

const (
	minPasswordSize = 8
	maxPasswordSize = 32
	modeEn          = "EN"
	modeDe          = "DE"

	codeFlow = 1
	codeMask = 2
	codeBuff = 3
)

var p KypdPacakge

// KypdPacakge 加解密包
type KypdPacakge struct {
	password string
	initmask string
	r        io.ReadCloser
	w        io.WriteCloser
}

// Init 数据初始化
func (p *KypdPacakge) Init() (err error) {
	password = strings.TrimSpace(password)
	if len(password) < minPasswordSize {
		return fmt.Errorf("password too short")
	}
	if len(password) > maxPasswordSize {
		return fmt.Errorf("password too long")
	}
	p.r, err = os.Open(input)
	if err != nil {
		return fmt.Errorf("open input file '%s' err: %v", input, err)
	}
	p.w, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open output file '%s' err: %v", output, err)
	}
	p.password = password
	p.initmask = utils.RandomPassword(maxPasswordSize)
	if mode == modeDe {
		// 解密
		p.r = utils.NewXorReader(p.r, []byte(mask))
	} else {
		// 加密
		p.w = utils.NewXorWriter(p.w, []byte(mask))
	}
	return nil
}

// Close 关闭文件
func (p *KypdPacakge) Close() {
	p.r.Close() //nolint
	p.w.Close() //nolint
}

// En 加密
func (p *KypdPacakge) En() error {
	h := &KypdHead{}
	h.preEn(p.password, p.initmask)
	if err := h.En(p.w); err != nil {
		return err
	}
	b := NewKypdBody(h)
	return b.En(p.r, p.w)
}

// De 解密
func (p *KypdPacakge) De() error {
	h := &KypdHead{
		password: p.password,
	}
	if err := h.De(p.r); err != nil {
		return err
	}
	if err := h.pstDe(); err != nil {
		return err
	}
	b := NewKypdBody(h)
	return b.De(p.r, p.w)
}

// KypdHead 加密描述
type KypdHead struct {
	version  uint32      // 版本
	hpns     *utils.Hpns // 加密策略：1.初始掩码；2.加密顺序；3.缓存范围；
	password string      // 密码
	initmask string      // 初始掩码
	minBuff  uint32
	maxBuff  uint32
	steps    []byte // 数据加密处理步骤
}

func (h *KypdHead) innerKey() []byte {
	ik := []byte{}
	k := byte('S') ^ byte('M')
	for _, b := range []byte(h.password) {
		k ^= b
	}
	ik = append(ik, k)
	switch h.version {
	case 1:
		ik = append(ik, ' ', '!', '0', '9', 'a', 'z', 'A', 'Z')
	}
	return ik
}

// preEn 加密之前设置
func (h *KypdHead) preEn(password, initmask string) {
	h.password = password
	h.initmask = initmask
	h.minBuff = enum.MB * 2
	h.maxBuff = enum.MB * 4
	h.steps = []byte{enum.Compress, enum.Encrypt, enum.Replace, enum.Shuffle}
	h.version = 1
	h.hpns = utils.NewHpns(
		&utils.Hpn{
			Code: codeFlow,
			Data: h.steps, // 压缩->加密->替换->乱序
		},
		&utils.Hpn{
			Code: codeMask,
			Data: []byte(h.initmask), // 初始掩码
		},
		&utils.Hpn{
			Code: codeBuff,
			Data: utils.GetLittleEndian(uint64(h.minBuff)<<enum.Bit32 + uint64(h.maxBuff)),
		})
}

// pstDe 加密之后设置
func (h *KypdHead) pstDe() error {
	for _, i := range h.hpns.Items {
		switch i.Code {
		case codeMask:
			h.initmask = string(i.Data)
		case codeFlow:
			h.steps = i.Data
		case codeBuff:
			buffSize := binary.LittleEndian.Uint64(i.Data)
			h.minBuff = uint32(buffSize >> enum.Bit32)
			h.maxBuff = uint32(buffSize)
		}
	}
	if len(h.initmask) == 0 {
		return fmt.Errorf("解析描述数据失败，非法的掩码")
	}
	if len(h.steps) == 0 {
		return fmt.Errorf("解析描述数据失败，非法加密流程")
	}
	if h.minBuff == 0 ||
		h.maxBuff == 0 ||
		h.minBuff > h.maxBuff {
		return fmt.Errorf("解析描述数据失败，无效的缓存值")
	}
	return nil
}

// En 加密
func (h *KypdHead) En(w io.Writer) error {
	data := h.hpns.En()
	utils.Xor1(h.innerKey(), data)
	cbit := utils.CalcCbit(data, uint16(h.version))
	vsc := utils.GetLittleEndian(uint64(h.version)<<enum.Bit32 + uint64(len(data))<<enum.Bit16 + uint64(cbit))
	if _, err := w.Write(vsc); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

// De 解密
func (h *KypdHead) De(r io.Reader) error {
	data := utils.GetLittleEndian(uint64(0))
	n, err := r.Read(data)
	if err != nil {
		return fmt.Errorf("读取描述数据失败，err: %v", err)
	}
	if n != len(data) {
		return fmt.Errorf("描述数据缺失")
	}
	vsc := binary.LittleEndian.Uint64(data)
	h.version = uint32(vsc >> enum.Bit32)
	size := uint16(vsc >> enum.Bit16)
	cbit1 := uint16(vsc)
	data = make([]byte, size)
	n, err = r.Read(data)
	if err != nil {
		return fmt.Errorf("读取加密描述数据失败，err: %v", err)
	}
	if n != int(size) {
		return fmt.Errorf("描述信息缺失")
	}
	cbit2 := utils.CalcCbit(data, uint16(h.version))
	if cbit1 != cbit2 {
		return fmt.Errorf("描述信息校验失败，数据被破坏")
	}
	utils.Xor1(h.innerKey(), data)
	h.hpns = &utils.Hpns{}
	h.hpns.De(data)
	return nil
}
