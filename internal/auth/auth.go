package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateSalt() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Return as base64 for storage
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// HashPassword hashes the password using Argon2id with the provided salt.
// Returns the base64 encoded hash.
func HashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	// Argon2id parameters
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), saltBytes, time, memory, threads, keyLen)
	return base64.RawStdEncoding.EncodeToString(hash), nil
}

func CheckPasswordHash(password, salt, hash string) bool {
	newHash, err := HashPassword(password, salt)
	if err != nil {
		return false
	}
	return newHash == hash
}

func GenerateToken(userID, role, secret string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
