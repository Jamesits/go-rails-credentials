package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// implements ActiveSupport::MessageEncryptor
// https://github.com/rails/rails/blob/04df9bc3d120b51447bde54caa56e9237cb8da0e/activesupport/lib/active_support/message_encryptor.rb

const (
	MasterKeyLengthBytes = 16

	// https://github.com/rails/rails/blob/04df9bc3d120b51447bde54caa56e9237cb8da0e/activesupport/lib/active_support/message_encryptor.rb#L119
	Separator = "--"

	GcmStandardNonceSize = 12

	// https://github.com/rails/rails/blob/04df9bc3d120b51447bde54caa56e9237cb8da0e/activesupport/lib/active_support/message_encryptor.rb#L118
	GcmTagSize = 16
)

var Base64Encoding = base64.StdEncoding

func RandomMasterKey() (string, error) {
	key := make([]byte, MasterKeyLengthBytes)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("random geneartor error: %w", err)
	}
	return hex.EncodeToString(key), nil
}

func Decrypt(MasterKey string, EncryptedFileContent string) (DecryptedFileContent []byte, err error) {
	key, err := hex.DecodeString(MasterKey)
	if err != nil {
		return nil, fmt.Errorf("decode master key failed: %w", err)
	}

	content := strings.SplitN(EncryptedFileContent, Separator, 3)
	if len(content) != 3 {
		return nil, fmt.Errorf("parse encrypted file failed")
	}
	cipherText, err := Base64Encoding.DecodeString(content[0])
	if err != nil {
		return nil, fmt.Errorf("parse content failed: %w", err)
	}
	iv, err := Base64Encoding.DecodeString(content[1])
	if err != nil {
		return nil, fmt.Errorf("parse IV failed: %w", err)
	}
	tag, err := Base64Encoding.DecodeString(content[2])
	if err != nil {
		return nil, fmt.Errorf("parse tag failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("initialize AES parser failed: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("initialize GCM parser failed: %w", err)
	}

	decryptedFileContent, err := gcm.Open(nil, iv, []byte(string(cipherText)+string(tag)), nil)

	if err != nil {
		return decryptedFileContent, fmt.Errorf("decrypt failed: %w", err)
	}
	return decryptedFileContent, nil
}

func Encrypt(MasterKey string, RawFileContent []byte) (EncryptedFileContent string, err error) {
	key, err := hex.DecodeString(MasterKey)
	if err != nil {
		return "", fmt.Errorf("decode master key failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("initialize AES parser failed: %w", err)
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return "", fmt.Errorf("initialize GCM parser failed: %w", err)
	}

	encryptedStream := gcm.Seal(nil, nil, RawFileContent, nil)
	if len(encryptedStream) < (GcmStandardNonceSize + len(RawFileContent) + GcmTagSize) {
		return "", fmt.Errorf("parser internal error")
	}
	iv := encryptedStream[0:GcmStandardNonceSize]
	content := encryptedStream[GcmStandardNonceSize : len(encryptedStream)-GcmTagSize]
	tag := encryptedStream[len(encryptedStream)-GcmTagSize:]

	sb := strings.Builder{}
	sb.WriteString(Base64Encoding.EncodeToString(content))
	sb.WriteString(Separator)
	sb.WriteString(Base64Encoding.EncodeToString(iv))
	sb.WriteString(Separator)
	sb.WriteString(Base64Encoding.EncodeToString(tag))
	return sb.String(), nil
}
