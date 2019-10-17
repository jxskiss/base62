package base62

func (enc *Encoding) FormatInt(num int64) []byte {
	dst := make([]byte, 0)
	return enc.AppendUint(dst, uint64(num))
}

func (enc *Encoding) FormatUint(num uint64) []byte {
	dst := make([]byte, 0)
	return enc.AppendUint(dst, num)
}

func (enc *Encoding) AppendInt(dst []byte, num int64) []byte {
	return enc.AppendUint(dst, uint64(num))
}

func (enc *Encoding) AppendUint(dst []byte, num uint64) []byte {
	if num == 0 {
		dst = append(dst, enc.encode[0])
		return dst
	}

	var buf [11]byte
	var i = 11
	for num > 0 {
		r := num % base
		num /= base
		i--
		buf[i] = enc.encode[r]
	}
	dst = append(dst, buf[i:]...)
	return dst
}

func (enc *Encoding) ParseInt(src []byte) (int64, error) {
	num, err := enc.ParseUint(src)
	if err != nil {
		return 0, err
	}
	return int64(num), nil
}

func (enc *Encoding) ParseUint(src []byte) (uint64, error) {
	var num uint64
	for i, c := range src {
		x := enc.decodeMap[c]
		if x == 0xFF {
			return 0, CorruptInputError(i)
		}
		num = num*base + uint64(x)
	}
	return num, nil
}

func FormatInt(num int64) []byte {
	return stdEncoding.FormatInt(num)
}

func FormatUint(num uint64) []byte {
	return stdEncoding.FormatUint(num)
}

func AppendInt(dst []byte, num int64) []byte {
	return stdEncoding.AppendInt(dst, num)
}

func AppendUint(dst []byte, num uint64) []byte {
	return stdEncoding.AppendUint(dst, num)
}

func ParseInt(src []byte) (int64, error) {
	return stdEncoding.ParseInt(src)
}

func ParseUint(src []byte) (uint64, error) {
	return stdEncoding.ParseUint(src)
}
