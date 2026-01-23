package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, first_name, last_name, email_verified)
		VALUES ($1, $2,$3,$4)
		RETURNING id, created_at, updated_at

	`,
		user.Email,
		user.FirstName,
		user.LastName,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, first_name, last_name, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.EmailVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// No user found
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}
