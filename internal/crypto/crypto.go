package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func MakeHash(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst)
}
