package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (github_id, username, email, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, user.GithubID, user.Username, user.Email, user.AvatarURL).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByGithubID(ctx context.Context, githubID string) (*domain.User, error) {
	query := `SELECT id, github_id, username, email, avatar_url, created_at, updated_at FROM users WHERE github_id = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, githubID).
		Scan(&user.ID, &user.GithubID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, github_id, username, email, avatar_url, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.GithubID, &user.Username, &user.Email, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
