package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Azanul/Next-Watch/internal/auth"
	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/services"
)

type Handler struct {
	userService      *services.UserService
	googleAuthClient *auth.GoogleAuthClient
}

func NewHandler(userService *services.UserService, googleAuthClient *auth.GoogleAuthClient) *Handler {
	return &Handler{
		userService:      userService,
		googleAuthClient: googleAuthClient,
	}
}

// AuthMiddleware checks for a user in the request and adds it to the context
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "No access token found", http.StatusUnauthorized)
			return
		}
		claims, err := h.googleAuthClient.ValidateToken(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Error(w, "Invalid Google token", http.StatusUnauthorized)
			return
		}

		user, err := h.userService.GetUserByEmail(ctx, claims.Email)
		if err != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Google sign in initiation handler
func (h *Handler) GoogleSignin(w http.ResponseWriter, r *http.Request) {
	authorizationURL, err := h.googleAuthClient.AuthorizationURL()
	if err != nil {
		http.Error(w, "Error signing in with google", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, authorizationURL, http.StatusFound)
}

// Google OAuth callback handler
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := h.googleAuthClient.Callback(r.FormValue("code"), r.FormValue("state"))
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	claims, err := h.googleAuthClient.GetUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		user = &models.User{
			Email: claims.Email,
			Name:  claims.Name,
		}
		err := h.userService.CreateUser(ctx, user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	}

	encryptedToken, err := auth.EncryptToken(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to encrypt token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    encryptedToken,
		Expires:  time.Now().Add(3600 * 24 * 7 * time.Second), // 1 week
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // Set to true in production
		HttpOnly: false, // Set to true in production
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
