package utils

import (
	"crypto/sha256"
	"fmt"
)

// HashString returns the SHA256 hash of the input string as a hexadecimal string.
func HashString(s string) string {
	return HashBytes([]byte(s))
}

// HashBytes returns the SHA256 hash of the input byte slice as a hexadecimal string.
func HashBytes(input []byte) string {
	hasher := sha256.New()
	hasher.Write(input)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// HashBytesWithSalt returns the SHA256 hash of the input string with a salt as a hexadecimal string.
func HashBytesWithSalt(input []byte, salt string) string {
	hasher := sha256.New()
	hasher.Write(input)
	hasher.Write([]byte(salt))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
