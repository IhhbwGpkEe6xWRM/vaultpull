// Package envencrypt provides symmetric encryption and decryption
// for secret values before writing them to .env files.
package envencrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the key length is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("envencrypt: key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when decryption fails.
var ErrInvalidCiphertext = errors.New("envencrypt: invalid ciphertext")

// Encrypter encrypts and decrypts string values using AES-GCM.
type Encrypter struct {
	key []byte
}

// New creates a new Encrypter with the given key.
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func New(key []byte) (*Encrypter, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKey
	}
	return &Encrypter{key: key}, nil
}

// Encrypt encrypts a plaintext string and returns a base64-encoded ciphertext.
func (e *Encrypter) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts a base64-encoded ciphertext and returns the plaintext.
func (e *Encrypter) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", ErrInvalidCiphertext
	}
	nonce, data := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plain), nil
}

// EncryptMap encrypts all values in a map, returning a new map.
func (e *Encrypter) EncryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := e.Encrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in a map, returning a new map.
func (e *Encrypter) DecryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := e.Decrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = dec
	}
	return out, nil
}
