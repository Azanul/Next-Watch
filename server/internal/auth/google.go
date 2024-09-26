package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleAuthClient struct {
	*oauth2.Config
	codeVerifiers map[string]string
	mu            sync.Mutex
}

func NewGoogleAuthClient() *GoogleAuthClient {
	var googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuthClient{
		googleOauthConfig,
		make(map[string]string),
		sync.Mutex{},
	}
}

func (g *GoogleAuthClient) AuthorizationURL() (string, error) {
	codeVerifier, verifierErr := randomBytesInHex(32)
	if verifierErr != nil {
		return "", fmt.Errorf("could not create a code verifier: %v", verifierErr)
	}

	sha2 := sha256.New()
	io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	state, err := randomBytesInHex(24)
	if err != nil {
		return "", fmt.Errorf("could not generate random state: %v", err)
	}

	g.mu.Lock()
	g.codeVerifiers[state] = codeVerifier
	g.mu.Unlock()

	return g.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	), nil
}

func (g *GoogleAuthClient) Callback(code string, state string) (*oauth2.Token, error) {
	g.mu.Lock()
	codeVerifier, exists := g.codeVerifiers[state]
	if !exists {
		g.mu.Unlock()
		return nil, errors.New("no matching code verifier found for state")
	}
	delete(g.codeVerifiers, state)
	g.mu.Unlock()

	token, err := g.Exchange(
		context.Background(),
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return nil, fmt.Errorf("error while exchanging token: %v", err)
	}

	return token, nil
}

func (g *GoogleAuthClient) ValidateToken(token string) (*GoogleClaims, error) {
	tokenString, err := decryptToken(token)
	if err != nil {
		return nil, errors.New("failed to decrypt token")
	}

	return g.GetUserInfo(tokenString)
}

func (g *GoogleAuthClient) GetUserInfo(accessToken string) (*GoogleClaims, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	oauth2Service, err := oauth2v2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.New("failed to create OAuth2 service")
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, errors.New("failed to get user info")
	}

	claims := &GoogleClaims{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}

	return claims, nil
}
