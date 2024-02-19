package base62

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Базовая структура средства кодирования в формат base63
type Base62Encoder struct{}

// Конструктор средства кодирования в base62 формат
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

// Метод кодирования строки в base62 формат
func (B62E *Base62Encoder) EncodeString() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return encode(rnd.Uint64())
}
