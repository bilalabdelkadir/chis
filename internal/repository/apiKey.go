package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiKeyRepository struct {
	pool *pgxpool.Pool
}

func NewApiKeyRepository(pool *pgxpool.Pool) *ApiKeyRepository {
	return &ApiKeyRepository{
		pool: pool,
	}
}

func (r *ApiKeyRepository) Create(ctx context.Context, apiKey *model.APIKey) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO api_keys (org_id, name,prefix, hashed_key, expires_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at, updated_at

	`,
		apiKey.OrgID,
		apiKey.Name,
		apiKey.Prefix,
		apiKey.HashedKey,
		apiKey.ExpiresAt,
	).Scan(&apiKey.ID, &apiKey.CreatedAt, &apiKey.UpdatedAt)

	return err
}

func (r *ApiKeyRepository) FindByHashedKey(ctx context.Context, hashedKey string) (*model.APIKey, error) {
	apiKey := &model.APIKey{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, org_id, name, prefix, hashed_key, expires_at, created_at, updated_at
		FROM api_keys
		WHERE hashed_key = $1
		LIMIT 1
	`, hashedKey).Scan(
		&apiKey.ID,
		&apiKey.OrgID,
		&apiKey.Name,
		&apiKey.Prefix,
		&apiKey.HashedKey,
		&apiKey.ExpiresAt,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return apiKey, nil
}
