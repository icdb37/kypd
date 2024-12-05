package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"

	"github.com/icdb37/kypd/compress"
	"github.com/icdb37/kypd/encrypt"
	"github.com/icdb37/kypd/enum"
	"github.com/icdb37/kypd/replace"
	"github.com/icdb37/kypd/shuffle"
	"github.com/icdb37/kypd/utils"
	"github.com/icdb37/kypd/utils/logx"
)

// processor 数据处理器
type processor interface {
	SetHpn(*utils.Hpn)
	GetHpn() *utils.Hpn
	En(datas []byte) []byte
	De(datas []byte) []byte
}

// KypdBody 分段包处理
type KypdBody struct {
	Password  string
	MinSize   uint32
	MaxSize   uint32
	WholeCbit uint16     // 整体校验码
	Header    utils.Hpns // 段数据头部
	Data      []byte     // 段数据内容
	process   []processor
	pEncrypt  *encrypt.Xorml
	iPwdpos   uint8  // 整体加密密码读取
	iDatpos   uint16 //
}

// NewKypdBody 创建加密数据处理对象
func NewKypdBody(h *KypdHead) *KypdBody {
	pCompres := compress.NewZlib()
	pEncrypt := encrypt.NewXorml(h.initmask, h.password)
	pReplace := replace.NewFixed()
	pShuffle := shuffle.NewFixed()
	procFloew := map[byte]processor{
		pCompres.Code(): pCompres,
		pEncrypt.Code(): pEncrypt,
		pReplace.Code(): pReplace,
		pShuffle.Code(): pShuffle,
	}
	b := &KypdBody{
		Password: h.password,
	}
	for _, i := range h.hpns.Items {
		switch i.Code {
		case codeMask:
			h.initmask = string(i.Data)
		case codeFlow:
			h.steps = i.Data
			for _, c := range i.Data {
				fp, ok := procFloew[c]
				if ok {
					b.process = append(b.process, fp)
				}
			}
		case codeBuff:
			buffSize := binary.LittleEndian.Uint64(i.Data)
			h.minBuff = uint32(buffSize >> enum.Bit32)
			h.maxBuff = uint32(buffSize)
			b.MaxSize = h.maxBuff
			b.MinSize = h.minBuff
			if h.maxBuff < enum.MB {
				b.Data = make([]byte, h.maxBuff*6/5+enum.KB)
			} else {
				b.Data = make([]byte, h.maxBuff*6/5)
			}
		}
	}
	b.pEncrypt = pEncrypt
	b.iPwdpos = uint8(len(b.Password)) - 3 // 【优化】根据版本或者整体干扰因子确定
	b.iDatpos = 0
	return b
}

// enSection 分段加密
func (s *KypdBody) enSection(data []byte) []byte {
	defer func() {
		s.pEncrypt.PrevMask = s.pEncrypt.CurrMask
		s.pEncrypt.Init()
	}()
	buf := bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	for _, p := range s.process {
		data = p.En(data)
		s.Header.Items = append(s.Header.Items, p.GetHpn())
	}
	headerData := s.Header.En()
	buf.Write(headerData)
	buf.Write(data)
	headerSize := uint16(len(headerData))
	bodySize := uint16(len(data))
	data = buf.Bytes()
	utils.Xor3([]string{s.Password, string([]byte{byte(headerSize), byte(bodySize)})}, data[enum.Byte8:], 0)
	s.WholeCbit = 0
	s.WholeCbit = utils.CalcCbit(data[enum.Byte8:], s.WholeCbit) // 计算校验和
	binary.LittleEndian.PutUint64(data[:enum.Byte8], uint64(s.WholeCbit)<<enum.Bit48+uint64(headerSize)<<enum.Bit32+uint64(bodySize))
	return data
}

// deSection 分段解密
func (s *KypdBody) deSection(data []byte, w io.Writer) error {
	for i := 0; ; i++ {
		if len(data) < enum.Byte8 {
			if i == 0 {
				return fmt.Errorf("section data size < 8byte") // 数据格式错误
			}
			break
		}
		secHeader := binary.LittleEndian.Uint64(data[:enum.Byte8])
		wholeCbit := uint16(secHeader >> enum.Bit48)
		headerSize := uint16(secHeader >> enum.Bit32)
		bodySize := uint16(secHeader)
		if len(data) < int(headerSize+bodySize+enum.Byte8) {
			if i == 0 {
				return fmt.Errorf("section data broken")
			}
			break
		}
		s.WholeCbit = 0
		s.WholeCbit = utils.CalcCbit(data[enum.Byte8:enum.Byte8+headerSize+bodySize], s.WholeCbit) // 计算校验和
		if s.WholeCbit != wholeCbit {
			return fmt.Errorf("check cbit failed, maybe data chaged") // 校验失败
		}
		utils.Xor3([]string{s.Password, string([]byte{byte(headerSize), byte(bodySize)})}, data[enum.Byte8:], 0)
		headerData := data[enum.Byte8 : enum.Byte8+headerSize]
		bodyData := data[enum.Byte8+headerSize : enum.Byte8+headerSize+bodySize]
		s.Header.De(headerData)
		for i := len(s.process) - 1; i >= 0; i-- {
			p := s.process[i]
			p.SetHpn(s.Header.Items[i])
			bodyData = p.De(bodyData)
		}
		if _, err := w.Write(bodyData); err != nil {
			return fmt.Errorf("encrypted section data write err: %v", err)
		}
		s.pEncrypt.PrevMask = s.pEncrypt.CurrMask
		logx.Debugf("section info cbit: %d, head_size: %d, body_size: %d", wholeCbit, headerSize, bodySize)
		data = data[enum.Byte8+headerSize+bodySize:]
	}
	copy(s.Data, data)
	s.iDatpos = uint16(len(data))
	return nil
}

// En 加密
func (s *KypdBody) En(r io.Reader, w io.Writer) error {
	for {
		size := s.MinSize + uint32(rand.Int())%(s.MaxSize-s.MinSize)
		n, err := r.Read(s.Data[:size])
		if err != nil && err != io.EOF {
			return fmt.Errorf("read original data err: %v", err)
		}
		if n > 0 {
			data := s.enSection(s.Data[:n])
			s.iPwdpos = utils.Xor0([]byte(s.Password), data, s.iPwdpos)
			if _, err := w.Write(data); err != nil {
				return err
			}
		}
		if err == io.EOF || n < int(size) {
			break
		}
	}
	return nil
}

// De 解密
func (s *KypdBody) De(r io.Reader, w io.Writer) error {
	for {
		n, err := r.Read(s.Data[s.iDatpos:])
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			s.iPwdpos = utils.Xor0([]byte(s.Password), s.Data[s.iDatpos:s.iDatpos+uint16(n)], s.iPwdpos)
			if err := s.deSection(s.Data[:s.iDatpos+uint16(n)], w); err != nil {
				return err
			}
		}
		if err == io.EOF || n < len(s.Data) {
			break
		}
	}
	return nil
}
