package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresApiKeyRepository struct {
	pool *pgxpool.Pool
}

func NewApiKeyRepository(pool *pgxpool.Pool) ApiKeyRepository {
	return &PostgresApiKeyRepository{
		pool: pool,
	}
}

func (r *PostgresApiKeyRepository) Create(ctx context.Context, apiKey *model.APIKey) error {
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

func (r *PostgresApiKeyRepository) FindByHashedKey(ctx context.Context, hashedKey string) (*model.APIKey, error) {
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

func (r *PostgresApiKeyRepository) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]*model.APIKey, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, org_id, name, prefix, expires_at, last_used_at, created_at, updated_at
		FROM api_keys
		WHERE org_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*model.APIKey
	for rows.Next() {
		k := &model.APIKey{}
		if err := rows.Scan(
			&k.ID, &k.OrgID, &k.Name, &k.Prefix,
			&k.ExpiresAt, &k.LastUsedAt, &k.CreatedAt, &k.UpdatedAt,
		); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	return keys, rows.Err()
}

func (r *PostgresApiKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM api_keys WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
