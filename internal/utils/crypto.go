package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

var encryptionKey []byte

// InitCrypto initializes the encryption key from environment variables.
// It must be called once at the start of the application.
func InitCrypto() error {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if len(keyStr) != 64 { // A chave deve ter 32 bytes, que são 64 caracteres hex
		return errors.New("ENCRYPTION_KEY environment variable must be a 64-character hex string (32 bytes)")
	}
	var err error
	encryptionKey, err = hex.DecodeString(keyStr)
	if err != nil {
		return fmt.Errorf("failed to decode ENCRYPTION_KEY: %w", err)
	}
	return nil
}

// Encrypt criptografa dados usando AES-GCM.
func Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt descriptografa dados usando AES-GCM.
func Decrypt(ciphertextHex string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
