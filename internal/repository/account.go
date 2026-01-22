package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {

	return &AccountRepository{
		pool: pool,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *model.Account) error {
	err := r.pool.QueryRow(ctx, `
	INSERT INTO accounts (user_id, provider, provider_account_id, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
`,
		account.UserID,
		account.Provider,
		account.ProviderAccountID,
		account.PasswordHash,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	return err
}

func (r *AccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	var a model.Account

	err := r.pool.QueryRow(ctx, `
SELECT id, user_id, provider, provider_account_id, password_hash, 
		access_token, refresh_token, token_expires_at, created_at, updated_at
FROM accounts
WHERE user_id = $1
`, userID).Scan(
		&a.ID,
		&a.UserID,
		&a.Provider,
		&a.ProviderAccountID,
		&a.PasswordHash,
		&a.AccessToken,
		&a.RefreshToken,
		&a.TokenExpiresAt,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No user found
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &a, nil
}
