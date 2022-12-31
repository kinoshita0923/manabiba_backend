package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
    sha256Binary := getSHA256Binary(password)
    hashedPassword := hex.EncodeToString(sha256Binary)
	return hashedPassword
}

func getSHA256Binary(s string) []byte {
    sha256Binary := sha256.Sum256([]byte(s))
    return sha256Binary[:]
}
