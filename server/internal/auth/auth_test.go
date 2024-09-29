package auth

import (
	"context"
	"os"
	"testing"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserFromContext(t *testing.T) {
	t.Run("User found in context", func(t *testing.T) {
		expectedUser := &models.User{ID: uuid.New(), Name: "testuser"}
		ctx := context.WithValue(context.Background(), "user", expectedUser)

		user, err := GetUserFromContext(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("User not found in context", func(t *testing.T) {
		ctx := context.Background()

		user, err := GetUserFromContext(ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})
}

func TestRandomBytesInHex(t *testing.T) {
	t.Run("Generate random bytes", func(t *testing.T) {
		count := 16
		hex1, err1 := randomBytesInHex(count)
		hex2, err2 := randomBytesInHex(count)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Len(t, hex1, count*2) // Each byte is represented by 2 hex characters
		assert.Len(t, hex2, count*2)
		assert.NotEqual(t, hex1, hex2) // Ensure randomness
	})

	t.Run("Invalid count", func(t *testing.T) {
		_, err := randomBytesInHex(-1)

		assert.Error(t, err)
	})
}

func TestEncryptDecryptToken(t *testing.T) {
	t.Run("Encrypt and decrypt token", func(t *testing.T) {
		os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef") // 32-byte key
		originalToken := "test-token"

		encryptedToken, err := EncryptToken(originalToken)
		assert.NoError(t, err)
		assert.NotEqual(t, originalToken, encryptedToken)

		decryptedToken, err := decryptToken(encryptedToken)
		assert.NoError(t, err)
		assert.Equal(t, originalToken, decryptedToken)
	})

	t.Run("Decrypt invalid token", func(t *testing.T) {
		os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
		invalidToken := "invalid-token"

		_, err := decryptToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("Encrypt with invalid key", func(t *testing.T) {
		os.Setenv("ENCRYPTION_KEY", "invalid-key")
		token := "test-token"

		_, err := EncryptToken(token)
		assert.Error(t, err)
	})
}
