// Package crypto provides encryption and decryption functionality
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// NewEncryptedWriter creates an encrypted writer that uses AES-256-GCM
func NewEncryptedWriter(w io.Writer, password string) (io.Writer, error) {
	// Derive key from password
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Write salt and nonce first
	if _, err := w.Write(salt); err != nil {
		return nil, err
	}
	if _, err := w.Write(nonce); err != nil {
		return nil, err
	}

	return &EncryptedWriter{
		writer: w,
		gcm:    gcm,
		nonce:  nonce,
	}, nil
}

// EncryptedWriter wraps an io.Writer to encrypt data
type EncryptedWriter struct {
	writer io.Writer
	gcm    cipher.AEAD
	nonce  []byte
}

// Write encrypts data and writes it to the underlying writer
func (ew *EncryptedWriter) Write(p []byte) (n int, err error) {
	encrypted := ew.gcm.Seal(nil, ew.nonce, p, nil)
	return ew.writer.Write(encrypted)
}

// NewEncryptedReader creates an encrypted reader that uses AES-256-GCM
func NewEncryptedReader(r io.Reader, password string) (io.Reader, error) {
	// Read salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(r, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Read nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	return &EncryptedReader{
		reader: r,
		gcm:    gcm,
		nonce:  nonce,
	}, nil
}

// EncryptedReader wraps an io.Reader to decrypt data
type EncryptedReader struct {
	reader io.Reader
	gcm    cipher.AEAD
	nonce  []byte
}

// Read decrypts data from the underlying reader
func (er *EncryptedReader) Read(p []byte) (n int, err error) {
	encrypted := make([]byte, len(p)+er.gcm.Overhead())
	n, err = er.reader.Read(encrypted)
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("read error: %w", err)
	}

	decrypted, err := er.gcm.Open(nil, er.nonce, encrypted[:n], nil)
	if err != nil {
		return 0, fmt.Errorf("decryption failed: %w", err)
	}

	copy(p, decrypted)
	return len(decrypted), nil
}
