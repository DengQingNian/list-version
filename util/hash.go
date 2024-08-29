package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func ToHashHex(filePath string) string {
	s1 := sha1.New()
	buf := s1.Sum([]byte(filePath))
	return hex.EncodeToString(buf)
}
