package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"github.com/princetheprogrammerbtw/nanoci/internal/queue"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	projectRepo domain.ProjectRepository
	buildRepo   domain.BuildRepository
	queue       *queue.RedisQueue
}

func NewWebhookHandler(p domain.ProjectRepository, b domain.BuildRepository, q *queue.RedisQueue) *WebhookHandler {
	return &WebhookHandler{
		projectRepo: p,
		buildRepo:   b,
		queue:       q,
	}
}

type githubPushPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Repository struct {
		ID   int64  `json:"id"`
		Name string `json:"full_name"`
	} `json:"repository"`
	HeadCommit struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	} `json:"head_commit"`
}

func (h *WebhookHandler) HandleGithub(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// In production, verify signature here using hmac with project.WebhookSecret
	// For now, we'll skip or log it.

	var event githubPushPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		zap.L().Error("failed to unmarshal github payload", zap.Error(err))
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Only handle push to branches (not tags)
	if !strings.HasPrefix(event.Ref, "refs/heads/") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	branch := strings.TrimPrefix(event.Ref, "refs/heads/")

	repoID := fmt.Sprintf("%d", event.Repository.ID)
	project, err := h.projectRepo.GetByGithubRepoID(r.Context(), repoID)
	if err != nil {
		zap.L().Error("failed to find project", zap.Error(err))
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	if project == nil {
		zap.L().Warn("received webhook for unknown project", zap.Int64("repo_id", event.Repository.ID))
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	build := &domain.Build{
		ProjectID:     project.ID,
		CommitHash:    event.HeadCommit.ID,
		CommitMessage: event.HeadCommit.Message,
		Branch:        branch,
		Status:        domain.BuildStatusPending,
	}

	if err := h.buildRepo.Create(r.Context(), build); err != nil {
		zap.L().Error("failed to create build", zap.Error(err))
		http.Error(w, "failed to trigger build", http.StatusInternalServerError)
		return
	}

	if err := h.queue.Enqueue(r.Context(), &queue.Job{BuildID: build.ID.String()}); err != nil {
		zap.L().Error("failed to enqueue job", zap.Error(err))
		// We might want to mark build as failed here
	}

	zap.L().Info("build triggered", zap.String("project", project.Name), zap.String("commit", build.CommitHash))
	w.WriteHeader(http.StatusAccepted)
}

func verifySignature(secret, signature string, payload []byte) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	signature = strings.TrimPrefix(signature, "sha256=")
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
