package credentials

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// MasterKeyLengthBytes is the length of the master key in bytes.
	// AES-128-GCM requires a 16-byte key.
	MasterKeyLengthBytes = 16
)

// RandomMasterKey generates a random master key.
func RandomMasterKey() (string, error) {
	key := make([]byte, MasterKeyLengthBytes)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("random geneartor error: %w", err)
	}
	return hex.EncodeToString(key), nil
}

func SanitizeMasterKey(in string) string {
	return strings.Trim(string(in), "\r\n")
}
