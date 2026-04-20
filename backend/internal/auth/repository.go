package auth

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, email, password, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
	)

	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
