package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const jwtPrivateKeyPathEnv = "AUTH_JWT_PRIVATE_KEY_PATH"

// LoadRSAPrivateKeyFromFile loads an RSA private key from a file path specified in environment variable.
func LoadRSAPrivateKeyFromFile() (*rsa.PrivateKey, error) {
	filePath := os.Getenv(jwtPrivateKeyPathEnv)
	if filePath == "" {
		return nil, fmt.Errorf("%s environment variable is not set", jwtPrivateKeyPathEnv)
	}

	pemBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read private key file %q: %w", filePath, err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM for RSA private key")
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q, want %q", block.Type, "PRIVATE KEY")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#8 private key: %w", err)
	}

	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key")
	}

	return key, nil
}

// GenerateToken generates a JWT token for the given user ID and email.
// Returns the token string, expiration timestamp, and any error that occurred.
func GenerateToken(userID uuid.UUID, email string, expirationDuration time.Duration) (string, int64, error) {
	exp := time.Now().Add(expirationDuration).Unix()
	privateKey, err := LoadRSAPrivateKeyFromFile()
	if err != nil {
		return "", 0, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"exp":   exp,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return tokenString, exp, nil
}
