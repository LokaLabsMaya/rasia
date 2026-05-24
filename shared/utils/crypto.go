package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/redhajuanda/komon/fail"
	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16
)

// DeriveKey derives a 32-byte AES key from password using Argon2id.
func DeriveKey(password string) []byte {
	salt := []byte("rasia-secrets-00") // fixed salt — key is already a secret
	return argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
}

// Encrypt encrypts plaintext with AES-256-GCM and returns base64-encoded ciphertext.
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fail.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fail.Wrap(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fail.Wrap(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded AES-256-GCM ciphertext.
func Decrypt(key []byte, encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fail.Wrap(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fail.Wrap(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fail.Wrap(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fail.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fail.Wrap(err)
	}

	return string(plaintext), nil
}
