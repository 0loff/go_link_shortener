package base62

import (
	"math/rand"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Base62Encoder struct{}

func NewBase62Encoder() *Base62Encoder {
	return &Base62Encoder{}
}

func encode(id uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder

	encodedBuilder.Grow(10)

	for ; id > 0; id = id / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(id % uint64(length))])
	}

	return encodedBuilder.String()
}

func (B62E *Base62Encoder) EncodeString() string {
	return encode(rand.Uint64())
}
