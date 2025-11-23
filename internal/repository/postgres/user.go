package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Snake1-1eyes/auth-service-it/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	row := r.db.QueryRowContext(ctx, query, login)

	var user domain.User
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}
