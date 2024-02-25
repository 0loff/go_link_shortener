package base62

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBase62Encode(t *testing.T) {
	testCases := []struct {
		name string
		seed int
		want string
	}{
		{
			name: "test POST request, text body",
			seed: 123456,
			want: "UI8H2R2QLnp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedString := NewBase62Encoder().EncodeString(uint64(tc.seed))
			require.Equal(t, tc.want, encodedString, "Base62 кодирование строки не соответствует ожиданию!")
		})
	}
}

func BenchmarkBase62Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBase62Encoder().EncodeString(uint64(123456))
	}
}
