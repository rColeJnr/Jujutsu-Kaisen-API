package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rcolejnr/jjk-api/internal/jwt"
	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/service"
)

// AuthHandler handles auth-related HTTP requests.
type AuthHandler struct {
	authSvc service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register creates a user and returns the created user and an access token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: "+err.Error()))
		return
	}

	resp, err := h.authSvc.Register(&req)
	if err != nil {
		if isValidationErr(err) {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJson(w, http.StatusCreated, resp)
}

// Login authenticates the user and returns a user and an access token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.authSvc.Login(&req)
	if err != nil {
		http.Error(w, "failed to authenticate:"+err.Error(), http.StatusUnauthorized)
		return
	}

	writeJson(w, http.StatusOK, resp)
}

// Update internal route  to update the admin password
// TODO: Invalidate current access_token after update.
func (h *AuthHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	claims := jwt.ClaimsFromContext(r.Context())

	err := h.authSvc.Update(claims.Email, &req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	writeJson(w, http.StatusNoContent, nil)
}
