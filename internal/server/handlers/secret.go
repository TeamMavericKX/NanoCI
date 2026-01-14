package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"github.com/princetheprogrammerbtw/nanoci/pkg/crypto"
	"github.com/princetheprogrammerbtw/nanoci/pkg/response"
)

type SecretHandler struct {
	repo          domain.SecretRepository
	encryptionKey []byte
}

func NewSecretHandler(repo domain.SecretRepository, key string) *SecretHandler {
	return &SecretHandler{
		repo:          repo,
		encryptionKey: []byte(key),
	}
}

func (h *SecretHandler) List(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	secrets, err := h.repo.ListByProjectID(r.Context(), projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, secrets)
}

func (h *SecretHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	encrypted, err := crypto.Encrypt(req.Value, h.encryptionKey)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "encryption failed")
		return
	}

	secret := &domain.Secret{
		ProjectID:      projectID,
		Key:            req.Key,
		EncryptedValue: encrypted,
	}

	if err := h.repo.Create(r.Context(), secret); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, secret)
}
