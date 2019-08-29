package base62

import (
	"math"
	"testing"
)

func Test_FormatUint_ParseUint(t *testing.T) {
	x := uint64(math.MaxUint64)
	dst := FormatUint(x)

	got, err := ParseUint(dst)
	t.Log(string(got))
	if err != nil {
		t.Fatalf("failed parse uint, err = %v", err)
	}
	if got != x {
		t.Fatalf("failed parse uint, got = %v, want = %v", got, x)
	}
}
