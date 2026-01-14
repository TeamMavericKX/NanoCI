package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"github.com/princetheprogrammerbtw/nanoci/pkg/response"
)

type BuildHandler struct {
	repo domain.BuildRepository
}

func NewBuildHandler(repo domain.BuildRepository) *BuildHandler {
	return &BuildHandler{repo: repo}
}

func (h *BuildHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	builds, err := h.repo.ListByProjectID(r.Context(), projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, builds)
}

func (h *BuildHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid build id")
		return
	}

	build, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if build == nil {
		response.Error(w, http.StatusNotFound, "build not found")
		return
	}

	response.JSON(w, http.StatusOK, build)
}
