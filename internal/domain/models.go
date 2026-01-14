package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	GithubID  string    `json:"github_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByGithubID(ctx context.Context, githubID string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type Project struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Name          string    `json:"name"`
	RepoURL       string    `json:"repo_url"`
	GithubRepoID  string    `json:"github_repo_id"`
	DefaultBranch string    `json:"default_branch"`
	WebhookSecret string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	GetByGithubRepoID(ctx context.Context, githubRepoID string) (*Project, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Project, error)
}
