package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var key []byte

func init() {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET is not set")
	}
	sum := sha256.Sum256([]byte(secret)) // 32-Byte key
	key = sum[:]
}

// Sign generates an encrypted token for <id>, valid until <exp>.
func Sign(id string, exp time.Time) (string, error) {
	plaintext := fmt.Sprintf("%s|%d", id, exp.Unix())

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	tokenBytes := append(nonce, ciphertext...) // nonce|ciphertext|tag
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// Verify ----- Verify token -----
func Verify(tok string) (string, bool) {
	raw, err := base64.URLEncoding.DecodeString(tok)
	if err != nil {
		return "", false
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}
	gcm, _ := cipher.NewGCM(block)

	if len(raw) < gcm.NonceSize() {
		return "", false
	}
	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", false //Auth-Tag or decryption failed.
	}
	parts := strings.Split(string(plain), "|")
	if len(parts) != 2 {
		return "", false // format error
	}
	id := parts[0]
	expUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", false
	}
	if time.Now().Unix() > expUnix {
		return "", false // expired
	}
	return id, true
}
