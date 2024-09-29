package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleAuthClient_AuthorizationURL(t *testing.T) {
	client := NewGoogleAuthClient()

	url, err := client.AuthorizationURL()

	assert.NoError(t, err)
	assert.Contains(t, url, "accounts.google.com/o/oauth2/auth")
	assert.Contains(t, url, "code_challenge_method=S256")
	assert.Contains(t, url, "code_challenge=")
	assert.Contains(t, url, "state=")

	assert.Len(t, client.codeVerifiers, 1)
}
