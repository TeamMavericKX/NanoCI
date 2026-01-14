package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
)

type projectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) domain.ProjectRepository {
	return &projectRepository{pool: pool}
}

func (r *projectRepository) Create(ctx context.Context, p *domain.Project) error {
	query := `
		INSERT INTO projects (user_id, name, repo_url, github_repo_id, default_branch, webhook_secret)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, p.UserID, p.Name, p.RepoURL, p.GithubRepoID, p.DefaultBranch, p.WebhookSecret).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	query := `SELECT id, user_id, name, repo_url, github_repo_id, default_branch, webhook_secret, created_at, updated_at 
			  FROM projects WHERE id = $1`
	var p domain.Project
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.UserID, &p.Name, &p.RepoURL, &p.GithubRepoID, &p.DefaultBranch, &p.WebhookSecret, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *projectRepository) GetByGithubRepoID(ctx context.Context, githubRepoID string) (*domain.Project, error) {
	query := `SELECT id, user_id, name, repo_url, github_repo_id, default_branch, webhook_secret, created_at, updated_at 
			  FROM projects WHERE github_repo_id = $1`
	var p domain.Project
	err := r.pool.QueryRow(ctx, query, githubRepoID).
		Scan(&p.ID, &p.UserID, &p.Name, &p.RepoURL, &p.GithubRepoID, &p.DefaultBranch, &p.WebhookSecret, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *projectRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	query := `SELECT id, user_id, name, repo_url, github_repo_id, default_branch, webhook_secret, created_at, updated_at 
			  FROM projects WHERE user_id = $1`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.RepoURL, &p.GithubRepoID, &p.DefaultBranch, &p.WebhookSecret, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}
