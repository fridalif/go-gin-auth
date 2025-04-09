package security

import (
	"crypto/sha256"
	"fmt"
)

func HashPassword(password string, salt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(salt+password)))
}
