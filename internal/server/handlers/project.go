package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"github.com/princetheprogrammerbtw/nanoci/pkg/response"
)

type ProjectHandler struct {
	repo domain.ProjectRepository
}

func NewProjectHandler(repo domain.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo: repo}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	// For now, assume a hardcoded user ID or get from context if middleware is ready
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		cookie, err := r.Cookie("user_id")
		if err == nil {
			userIDStr = cookie.Value
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	projects, err := h.repo.ListByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if project == nil {
		response.Error(w, http.StatusNotFound, "project not found")
		return
	}

	response.JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		cookie, err := r.Cookie("user_id")
		if err == nil {
			userIDStr = cookie.Value
		}
	}
	userID, _ := uuid.Parse(userIDStr)
	p.UserID = userID

	if err := h.repo.Create(r.Context(), &p); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, p)
}
