package base62

import "testing"

func BenchmarkBase62Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBase62Encoder().EncodeString()
	}
}
