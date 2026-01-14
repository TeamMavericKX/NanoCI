package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
)

type secretRepository struct {
	pool *pgxpool.Pool
}

func NewSecretRepository(pool *pgxpool.Pool) domain.SecretRepository {
	return &secretRepository{pool: pool}
}

func (r *secretRepository) Create(ctx context.Context, s *domain.Secret) error {
	query := `
		INSERT INTO secrets (project_id, key, encrypted_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, key) DO UPDATE SET encrypted_value = EXCLUDED.encrypted_value
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, s.ProjectID, s.Key, s.EncryptedValue).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *secretRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*domain.Secret, error) {
	query := `SELECT id, project_id, key, encrypted_value, created_at FROM secrets WHERE project_id = $1`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*domain.Secret
	for rows.Next() {
		var s domain.Secret
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Key, &s.EncryptedValue, &s.CreatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, &s)
	}
	return secrets, nil
}

func (r *secretRepository) Delete(ctx context.Context, projectID uuid.UUID, key string) error {
	query := `DELETE FROM secrets WHERE project_id = $1 AND key = $2`
	_, err := r.pool.Exec(ctx, query, projectID, key)
	return err
}
