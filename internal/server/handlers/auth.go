package handlers

import (
	"net/http"

	"github.com/princetheprogrammerbtw/nanoci/internal/auth"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// In a real app, use a secure random state and store it in a cookie/session
	state := "random-state"
	url := h.authService.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != "random-state" {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	user, err := h.authService.HandleCallback(r.Context(), code)
	if err != nil {
		zap.L().Error("auth callback failed", zap.Error(err))
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// In a real app, create a JWT or session here
	zap.L().Info("user logged in", zap.String("username", user.Username))
	w.Write([]byte("Welcome " + user.Username))
}
