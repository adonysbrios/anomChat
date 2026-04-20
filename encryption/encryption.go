package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

// deriveKey convierte una clave arbitraria en una clave AES-256 (32 bytes) usando SHA-256.
func deriveKey(key []byte) []byte {
	sum := sha256.Sum256(key)
	out := make([]byte, 32)
	copy(out, sum[:])
	return out
}

// EncryptMessage devuelve: [NONCE (12 bytes)] + [CIPHERTEXT || TAG]
func EncryptMessage(plaintext []byte, key []byte) ([]byte, error) {
	key = deriveKey(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	out := make([]byte, 0, len(nonce)+len(ciphertext))
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out, nil
}

// DecryptMessage espera: [NONCE (12 bytes)] + [CIPHERTEXT || TAG]
func DecryptMessage(ciphertext []byte, key []byte) ([]byte, error) {
	key = deriveKey(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize+16 { // 16 bytes = tag GCM
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return plaintext, nil
}
