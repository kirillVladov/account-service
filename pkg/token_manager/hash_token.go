package token_manager

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
