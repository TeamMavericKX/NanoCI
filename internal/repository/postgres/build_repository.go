package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
)

type buildRepository struct {
	pool *pgxpool.Pool
}

func NewBuildRepository(pool *pgxpool.Pool) domain.BuildRepository {
	return &buildRepository{pool: pool}
}

func (r *buildRepository) Create(ctx context.Context, b *domain.Build) error {
	query := `
		INSERT INTO builds (project_id, commit_hash, commit_message, branch, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, b.ProjectID, b.CommitHash, b.CommitMessage, b.Branch, b.Status).
		Scan(&b.ID, &b.CreatedAt)
}

func (r *buildRepository) Update(ctx context.Context, b *domain.Build) error {
	query := `
		UPDATE builds 
		SET status = $1, started_at = $2, finished_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, b.Status, b.StartedAt, b.FinishedAt, b.ID)
	return err
}

func (r *buildRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Build, error) {
	query := `SELECT id, project_id, commit_hash, commit_message, branch, status, started_at, finished_at, created_at 
			  FROM builds WHERE id = $1`
	var b domain.Build
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&b.ID, &b.ProjectID, &b.CommitHash, &b.CommitMessage, &b.Branch, &b.Status, &b.StartedAt, &b.FinishedAt, &b.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &b, err
}

func (r *buildRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domain.Build, error) {
	query := `SELECT id, project_id, commit_hash, commit_message, branch, status, started_at, finished_at, created_at 
			  FROM builds WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []*domain.Build
	for rows.Next() {
		var b domain.Build
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.CommitHash, &b.CommitMessage, &b.Branch, &b.Status, &b.StartedAt, &b.FinishedAt, &b.CreatedAt); err != nil {
			return nil, err
		}
		builds = append(builds, &b)
	}
	return builds, nil
}
