package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashToken returns the hex-encoded SHA-256 hash of the given raw string.
// Used internally to store refresh token hashes in the database.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
