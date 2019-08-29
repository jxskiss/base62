package base62

import (
	"bytes"
	"crypto/rand"
	mathrand "math/rand"
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
		dst := Encode(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecode_0xFF(t *testing.T) {
	for i := 0; i < 1000; i++ {
		src := make([]byte, i)
		for i := range src {
			src[i] = 0xff
		}
		dst := Encode(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecode_RandomBytes(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		src := make([]byte, 32+mathrand.Intn(32))
		_, _ = rand.Read(src)
		dst := Encode(src)
		got, err := Decode(dst)
		if err != nil {
			t.Fatalf("failed decode, err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
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
