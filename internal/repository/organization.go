package repository

import (
	"context"
	"fmt"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) OrganizationRepository {
	return &PostgresOrganizationRepository{
		pool: pool,
	}
}

func (r *PostgresOrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO organizations (name, slug, signing_secret)
		VALUES ($1,$2,$3)
		RETURNING id, created_at, updated_at
	`,
		organization.Name,
		organization.Slug,
		organization.SigningSecret,
	).Scan(&organization.ID, &organization.CreatedAt, &organization.UpdatedAt)

	return err
}

func (r *PostgresOrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	org := &model.Organization{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, signing_secret, created_at, updated_at
		FROM organizations WHERE id = $1
	`, id).Scan(&org.ID, &org.Name, &org.Slug, &org.SigningSecret, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	return err
}

func (r *PostgresOrganizationRepository) GetSigningSecret(ctx context.Context, orgID uuid.UUID) (string, error) {
	var secret string
	err := r.pool.QueryRow(ctx, `
		SELECT signing_secret FROM organizations WHERE id = $1
	`, orgID).Scan(&secret)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func (r *PostgresOrganizationRepository) RotateSigningSecret(ctx context.Context, orgID uuid.UUID, newSecret string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE organizations SET signing_secret = $1, updated_at = NOW() WHERE id = $2
	`, newSecret, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}
