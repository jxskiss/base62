# base62

base62 is a compact and fast implementation of `Base62` encoding/decoding algorithm,
which is inspired by the [java implementation by glowfall](https://github.com/glowfall/base62).

As this package was written, I did not find good base62 implementation on github.
The following packages can only encode/decode integers:

    - https://github.com/gravityblast/go-base62/
    - https://github.com/kare/base62
    - https://github.com/catinello/base62

While this package encodes/decodes bytes of arbitrary length, as well as integers.

This `Base62` implementation is much faster than `big.Int` based implementation,
and is not much slower than typical `Base64` implementations. See the following
benchmark results.

## Benchmark

```text
Benchmark_Encode-12                      7054132               146.4 ns/op
Benchmark_EncodeToString-12              8101567               146.2 ns/op
Benchmark_Decode-12                     15481666                73.60 ns/op
Benchmark_DecodeString-12               16301325                74.36 ns/op

Benchmark_EncodeToBuf-12                 9724098               126.8 ns/op
Benchmark_DecodeToBuf-12                97695962                12.21 ns/op

Benchmark_Encode_Integer-12             29119437                41.30 ns/op
Benchmark_Decode_Integer-12             120328183                9.917 ns/op

Benchmark_Encode_BigInt-12               1000000              1048 ns/op

Benchmark_Base64_EncodeToString-12      19974897                57.41 ns/op
Benchmark_Base64_DecodeString-12        19884616                55.09 ns/op

Benchmark_Base64_Encode-12              68163142                17.93 ns/op
Benchmark_Base64_Decode-12              41990004                28.25 ns/op
```
