package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func Sha256(content string) string {
	h := sha256.New()

	h.Write([]byte(content))

	return hex.EncodeToString(h.Sum(nil))
}

func Sha256Bin(content string) string {
	h := sha256.New()

	h.Write([]byte(content))
	hex := hex.EncodeToString(h.Sum(nil))

	var sb strings.Builder

	for i := 0; i < len(hex); i += 8 {
		n, _ := strconv.ParseUint(hex[i:i+8], 16, 32)
		sb.WriteString(fmt.Sprintf("%032b", n))
	}

	return sb.String()
}
