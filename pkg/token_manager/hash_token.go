package token_manager

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(plainToken string) string {
	sum := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(sum[:])
}
