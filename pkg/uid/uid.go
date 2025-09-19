package uid

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

// New 生成带前缀的随机字符串，用于各业务主键。
func New(prefix string) string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return prefix + hex.EncodeToString(b[:])
	}

	return prefix + strconv.FormatInt(time.Now().UnixNano(), 10)
}
