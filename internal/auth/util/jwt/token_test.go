package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestPrivateKeyFile creates a temporary file with a valid RSA private key
func setupTestPrivateKeyFile(t *testing.T) (string, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	tmpFile, err := os.CreateTemp("", "test-private-key-*.pem")
	require.NoError(t, err)

	_, err = tmpFile.Write(privateKeyPEM)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name(), privateKey
}

func TestLoadRSAPrivateKeyFromFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		filePath, expectedKey := setupTestPrivateKeyFile(t)

		key, err := LoadRSAPrivateKeyFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, key)
		assert.Equal(t, expectedKey.N, key.N)
		assert.Equal(t, expectedKey.E, key.E)
	})

	t.Run("EmptyPath", func(t *testing.T) {
		key, err := LoadRSAPrivateKeyFromFile("")

		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "read private key file")
	})

	t.Run("FileNotFound", func(t *testing.T) {
		filePath := "/nonexistent/path/to/key.pem"

		key, err := LoadRSAPrivateKeyFromFile(filePath)

		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "read private key file")
	})

	t.Run("InvalidPEM", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-invalid-*.pem")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("not a valid PEM content")
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		key, err := LoadRSAPrivateKeyFromFile(tmpFile.Name())

		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "failed to parse PEM")
	})

	t.Run("WrongPEMType", func(t *testing.T) {
		// Create a PEM file with wrong type (e.g., CERTIFICATE instead of PRIVATE KEY)
		tmpFile, err := os.CreateTemp("", "test-wrong-type-*.pem")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte("fake certificate data"),
		}
		pemData := pem.EncodeToMemory(block)

		_, err = tmpFile.Write(pemData)
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		key, err := LoadRSAPrivateKeyFromFile(tmpFile.Name())

		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "unexpected PEM type")
	})

	t.Run("InvalidKeyFormat", func(t *testing.T) {
		// Create a PEM file with invalid key data
		tmpFile, err := os.CreateTemp("", "test-invalid-key-*.pem")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: []byte("invalid key data"),
		}
		pemData := pem.EncodeToMemory(block)

		_, err = tmpFile.Write(pemData)
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		key, err := LoadRSAPrivateKeyFromFile(tmpFile.Name())

		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Contains(t, err.Error(), "parse PKCS#8 private key")
	})
}

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	expirationDuration := 15 * time.Minute

	t.Run("Success", func(t *testing.T) {
		filePath, _ := setupTestPrivateKeyFile(t)
		tokenClaim := TokenClaim{TokenPath: filePath}

		tokenString, exp, err := tokenClaim.GenerateToken(userID, email, expirationDuration)

		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		assert.Greater(t, exp, int64(0))

		// Verify expiration is approximately correct (within 1 second)
		expectedExp := time.Now().Add(expirationDuration).Unix()
		assert.InDelta(t, expectedExp, exp, 1)
	})

	t.Run("LoadKeyError", func(t *testing.T) {
		filePath := "/nonexistent/path/to/key.pem"
		tokenClaim := TokenClaim{TokenPath: filePath}

		tokenString, exp, err := tokenClaim.GenerateToken(userID, email, expirationDuration)

		assert.Error(t, err)
		assert.Empty(t, tokenString)
		assert.Equal(t, int64(0), exp)
	})

	t.Run("DifferentUsers", func(t *testing.T) {
		filePath, _ := setupTestPrivateKeyFile(t)
		tokenClaim := TokenClaim{TokenPath: filePath}

		userID1 := uuid.New()
		userID2 := uuid.New()

		token1, _, err1 := tokenClaim.GenerateToken(userID1, "user1@example.com", expirationDuration)
		token2, _, err2 := tokenClaim.GenerateToken(userID2, "user2@example.com", expirationDuration)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		assert.NotEqual(t, token1, token2, "Tokens for different users should be different")
	})

	t.Run("DifferentExpirationDurations", func(t *testing.T) {
		filePath, _ := setupTestPrivateKeyFile(t)
		tokenClaim := TokenClaim{TokenPath: filePath}

		token1, exp1, err1 := tokenClaim.GenerateToken(userID, email, 5*time.Minute)
		token2, exp2, err2 := tokenClaim.GenerateToken(userID, email, 30*time.Minute)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		assert.NotEqual(t, exp1, exp2, "Expiration times should be different")
		assert.Greater(t, exp2, exp1, "Longer duration should have later expiration")
	})
}
