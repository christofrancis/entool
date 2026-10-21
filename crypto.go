package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/scrypt"
)

const (
	formatVersion  = "ET1"
	maxCharacters  = 1_000_000
	saltSize       = 16
	nonceSize      = 12
	authTagSize    = 16
	lengthSize     = 4
	fixedOverhead  = saltSize + nonceSize + authTagSize + lengthSize
	fillerAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
)

var errInvalidEncryptedFile = errors.New("wrong key, unsupported format, or damaged file")

// minimumCharacters returns the minimum ET1 length for a UTF-8 secret.
// ET1 stays compatible with older Python versions of EnTool.
func minimumCharacters(secret string) int {
	return len(formatVersion) + 4*((fixedOverhead+len([]byte(secret))+2)/3)
}

func deriveKey(passphrase string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
}

func encryptSecret(secret, passphrase string, characters int) (string, error) {
	if len(secret) == 0 || !utf8.ValidString(secret) {
		return "", errors.New("enter a valid, nonempty UTF-8 secret")
	}
	if len(passphrase) == 0 {
		return "", errors.New("the encryption key cannot be empty")
	}
	if characters < minimumCharacters(secret) {
		return "", fmt.Errorf("secret is too long for this file size; choose at least %d characters", minimumCharacters(secret))
	}
	if characters > maxCharacters {
		return "", fmt.Errorf("maximum file length is %d characters", maxCharacters)
	}

	// A fixed ET1 overhead and floor division preserve older EnTool file lengths.
	encodedBytes := 3 * ((characters - len(formatVersion)) / 4)
	paddingSize := encodedBytes - fixedOverhead - len([]byte(secret))
	if paddingSize < 0 {
		return "", errors.New("invalid ET1 file length")
	}
	plain := make([]byte, lengthSize+len([]byte(secret))+paddingSize)
	binary.BigEndian.PutUint32(plain[:lengthSize], uint32(len([]byte(secret))))
	copy(plain[lengthSize:], secret)
	if _, err := io.ReadFull(rand.Reader, plain[lengthSize+len([]byte(secret)):]); err != nil {
		return "", fmt.Errorf("generate random padding: %w", err)
	}
	salt := make([]byte, saltSize)
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return "", fmt.Errorf("derive encryption key: %w", err)
	}
	defer clear(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := aead.Seal(nil, nonce, plain, []byte(formatVersion))
	payload := make([]byte, 0, saltSize+nonceSize+len(ciphertext))
	payload = append(payload, salt...)
	payload = append(payload, nonce...)
	payload = append(payload, ciphertext...)
	result := formatVersion + base64.URLEncoding.EncodeToString(payload)

	// Exactly 0-3 random printable filler characters for the requested length.
	if extra := characters - len(result); extra > 0 {
		filler := make([]byte, extra)
		if _, err := io.ReadFull(rand.Reader, filler); err != nil {
			return "", fmt.Errorf("generate filler: %w", err)
		}
		for i, randomByte := range filler {
			filler[i] = fillerAlphabet[int(randomByte)&63]
		}
		result += string(filler)
	}
	if len(result) != characters {
		return "", errors.New("internal error: incorrect ET1 length")
	}
	return result, nil
}

func decryptSecret(contents, passphrase string) (string, error) {
	if passphrase == "" || len(contents) < minimumCharacters("") || len(contents) > maxCharacters || !strings.HasPrefix(contents, formatVersion) {
		return "", errInvalidEncryptedFile
	}
	encodedLength := 4 * ((len(contents) - len(formatVersion)) / 4)
	raw, err := base64.URLEncoding.Strict().DecodeString(contents[len(formatVersion) : len(formatVersion)+encodedLength])
	if err != nil || len(raw) < saltSize+nonceSize+authTagSize+lengthSize {
		return "", errInvalidEncryptedFile
	}
	salt := raw[:saltSize]
	nonce := raw[saltSize : saltSize+nonceSize]
	ciphertext := raw[saltSize+nonceSize:]
	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return "", errInvalidEncryptedFile
	}
	defer clear(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errInvalidEncryptedFile
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", errInvalidEncryptedFile
	}
	plain, err := aead.Open(nil, nonce, ciphertext, []byte(formatVersion))
	if err != nil || len(plain) < lengthSize {
		return "", errInvalidEncryptedFile
	}
	secretLength := binary.BigEndian.Uint32(plain[:lengthSize])
	if uint64(secretLength) > uint64(len(plain)-lengthSize) {
		return "", errInvalidEncryptedFile
	}
	secret := plain[lengthSize : lengthSize+int(secretLength)]
	if !utf8.Valid(secret) {
		return "", errInvalidEncryptedFile
	}
	return string(secret), nil
}
