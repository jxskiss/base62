package base62

import (
	"math/bits"
	"reflect"
	"strconv"
	"unsafe"
)

const (
	base        = 62
	compactMask = 0x1E // 00011110
	mask5bits   = 0x1F // 00011111
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type Encoding struct {
	encode    [base]byte
	decodeMap [256]byte
}

func NewEncoding(encoder string) *Encoding {
	if len(encoder) != base {
		panic("encoding alphabet is not 62-bytes long")
	}
	for i := 0; i < len(encoder); i++ {
		if encoder[i] == '\n' || encoder[i] == '\r' {
			panic("encoding alphabet contains newline character")
		}
	}

	e := new(Encoding)
	copy(e.encode[:], encoder)
	for i := 0; i < len(e.decodeMap); i++ {
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}
	return e
}

var stdEncoding = NewEncoding(encodeStd)

func (enc *Encoding) Encode(src []byte) []byte {
	if len(src) == 0 {
		return []byte{}
	}
	encoder := newEncoder(src)
	return encoder.encode(enc.encode[:])
}

func (enc *Encoding) EncodeToString(src []byte) string {
	ret := enc.Encode(src)
	return b2s(ret)
}

type encoder struct {
	b   []byte
	pos int
}

func newEncoder(src []byte) *encoder {
	return &encoder{
		b:   src,
		pos: len(src) * 8,
	}
}

func (enc *encoder) next() (byte, bool) {
	var i, pos int
	var j, blen byte
	pos = enc.pos - 6
	if pos <= 0 {
		pos = 0
		blen = byte(enc.pos)
	} else {
		i = pos / 8
		j = byte(pos % 8)
		blen = byte((i+1)*8 - pos)
		if blen > 6 {
			blen = 6
		}
	}
	shift := 8 - j - blen
	b := enc.b[i] >> shift & (1<<blen - 1)

	if blen < 6 && pos > 0 {
		blen1 := 6 - blen
		b = b<<blen1 | enc.b[i+1]>>(8-blen1)
	}
	if b&compactMask == compactMask {
		if pos > 0 || b > mask5bits {
			pos++
		}
		b &= mask5bits
	}
	enc.pos = pos

	return b, pos > 0
}

func (enc *encoder) encode(encTable []byte) []byte {
	ret := make([]byte, 0, len(enc.b)*8/5+1)
	x, hasMore := enc.next()
	for {
		ret = append(ret, encTable[x])
		if !hasMore {
			break
		}
		x, hasMore = enc.next()
	}
	return ret
}

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base62 data at input byte " + strconv.FormatInt(int64(e), 10)
}

func (enc *Encoding) Decode(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	dec := decoder(src)
	return dec.decode(enc.decodeMap[:])
}

func (enc *Encoding) DecodeString(src string) ([]byte, error) {
	b := s2b(src)
	return enc.Decode(b)
}

type decoder []byte

func (dec decoder) decode(decTable []byte) ([]byte, error) {
	ret := make([]byte, len(dec)*6/8+1)
	idx := len(ret)
	pos := byte(0)
	b := 0
	for i, c := range dec {
		x := decTable[c]
		if x == 0xFF {
			return nil, CorruptInputError(i)
		}
		if i == len(dec)-1 {
			b |= int(x) << pos
			pos += byte(bits.Len8(x))
		} else if x&compactMask == compactMask {
			b |= int(x) << pos
			pos += 5
		} else {
			b |= int(x) << pos
			pos += 6
		}
		if pos >= 8 {
			idx--
			ret[idx] = byte(b)
			pos %= 8
			b >>= 8
		}
	}
	if pos > 0 {
		idx--
		ret[idx] = byte(b)
	}

	return ret[idx:], nil
}

func Encode(src []byte) []byte {
	return stdEncoding.Encode(src)
}

func EncodeToString(src []byte) string {
	return stdEncoding.EncodeToString(src)
}

func Decode(src []byte) ([]byte, error) {
	return stdEncoding.Decode(src)
}

func DecodeString(src string) ([]byte, error) {
	return stdEncoding.DecodeString(src)
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := &reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bh))
}
