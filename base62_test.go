package base62

import (
	"bytes"
	"crypto/rand"
	mathrand "math/rand"
	"strings"
	"testing"
)

func Test_EncodeDecode(t *testing.T) {
	src := []byte("Hello, 世界！")
	dst := Encode(src)
	got, err := Decode(dst)
	if err != nil {
		t.Fatalf("failed decode, err = %v", err)
	}
	if !bytes.Equal(src, got) {
		t.Fatalf("failed decode, got = %v, want = %v", got, src)
	}

	dstStr := EncodeToString(src)
	got, _ = DecodeString(dstStr)
	if !bytes.Equal(src, got) {
		t.Fatalf("failed decode string, got = %v, want = %v", got, src)
	}
}

func Test_EncodeDecode_Zeros(t *testing.T) {
	for i := 0; i < 1000; i++ {
		src := make([]byte, i)
		dst := StdEncoding._encodeV1(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}

		// Make sure the new implementation is compatible with the old.
		v2Dst := StdEncoding._encodeV2(src)
		if !bytes.Equal(dst, v2Dst) {
			t.Logf("src= %v\n  v1= %v\n  v2= %v", src, dst, v2Dst)
			t.Fatalf("encode new implementation not equal to v1")
		}
	}
}

func Test_EncodeDecode_0xFF(t *testing.T) {
	for i := 0; i < 1000; i++ {
		src := make([]byte, i)
		for i := range src {
			src[i] = 0xff
		}
		dst := StdEncoding._encodeV1(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}

		// Make sure the new implementation is compatible with the old.
		v2Dst := StdEncoding._encodeV2(src)
		if !bytes.Equal(dst, v2Dst) {
			t.Logf("src= %v\n  v1= %v\n  v2= %v", src, dst, v2Dst)
			t.Fatalf("encode new implementation not equal to v1")
		}
	}
}

func Test_EncodeDecode_RandomBytes(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		src := make([]byte, 32+mathrand.Intn(32))
		_, _ = rand.Read(src)
		dst := StdEncoding._encodeV1(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode, err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}

		// Make sure the new implementation is compatible with the old.
		v2Dst := StdEncoding._encodeV2(src)
		if !bytes.Equal(dst, v2Dst) {
			t.Logf("src= %v\n  v1= %v\n  v2= %v", src, dst, v2Dst)
			t.Fatalf("encode new implementation not equal to v1")
		}
	}
}

func Test_EncodeToBuf(t *testing.T) {
	buf := make([]byte, 0, 1000)
	for i := 0; i < 10000; i++ {
		src := make([]byte, 32+mathrand.Intn(100))
		_, _ = rand.Read(src)
		want := Encode(src)

		got1 := EncodeToBuf(make([]byte, 0, 2), src)
		if !bytes.Equal(want, got1) {
			t.Fatal("incorrect result from EncodeToBuf")
		}

		got2 := EncodeToBuf(buf, src)
		if !bytes.Equal(want, got2) {
			t.Fatal("incorrect result from EncodeToBuf")
		}
	}
}

func TestDecodeToBuf(t *testing.T) {
	buf := make([]byte, 0, 1000)
	for i := 0; i < 10000; i++ {
		src := make([]byte, 32+mathrand.Intn(100))
		_, _ = rand.Read(src)
		encoded := Encode(src)

		got1, err := DecodeToBuf(make([]byte, 0, 2), encoded)
		if err != nil {
			t.Fatalf("failed DecodeToBuf, err = %v", err)
		}
		if !bytes.Equal(src, got1) {
			t.Fatalf("incorrect result from DecodeToBuf, encoded = %v", encoded)
		}

		got2, err := DecodeToBuf(buf, encoded)
		if err != nil {
			t.Fatalf("failed DecodeToBuf, err = %v", err)
		}
		if !bytes.Equal(src, got2) {
			t.Fatalf("incorrect result from DecodeToBuf, encoded = %v", encoded)
		}
	}
}

// ----------

func Test_encoder_next(t *testing.T) {
	src := []byte{123, 234, 255}
	enc := newEncoder(src)

	//for _, w := range src {
	//	fmt.Printf("%08b", w)
	//}
	//fmt.Println()

	x, hasMore := enc.next()
	for {
		_ = x
		if !hasMore {
			break
		}
		x, hasMore = enc.next()
	}
}

// ----------

func Test_NewEncoding_panic(t *testing.T) {
	func() {
		encoder := "abcdef"
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewEncoding did not panic with encoder %q", encoder)
			}
		}()
		_ = NewEncoding(encoder)
	}()

	func() {
		encoder := []byte(encodeStd)
		encoder[1] = '\n'
		defer func() {
			if r := recover(); r == nil {
				t.Error("NewEncoding did not panic with encoder contains \\n")
			}
		}()
		_ = NewEncoding(string(encoder))
	}()

	func() {
		encoder := []byte(encodeStd)
		encoder[1] = '\r'
		defer func() {
			if r := recover(); r == nil {
				t.Error("NewEncoding did not panic with encoder contains \\r")
			}
		}()
		_ = NewEncoding(string(encoder))
	}()
}

func Test_Decode_CorruptInputError(t *testing.T) {
	src := make([]byte, 256)
	for i := range src {
		src[i] = byte(i)
	}
	_, err := StdEncoding.Decode(src)
	if err == nil || !strings.Contains(err.Error(), "illegal base62 data at input byte") {
		t.Fatal("decoding invalid data did not return CorruptInputError")
	}
}
