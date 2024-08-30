package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

func ToHashHex(filePath string) string {
	h := md5.New()
	io.WriteString(h, filePath)
	return fmt.Sprintf("%x", h.Sum(nil))
}
