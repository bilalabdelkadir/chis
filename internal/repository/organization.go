package repository

import (
	"context"

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
		INSERT INTO organizations (name, slug)
		VALUES ($1,$2)
		RETURNING id, created_at, updated_at

	`,
		organization.Name,
		organization.Slug,
	).Scan(&organization.ID, &organization.CreatedAt, &organization.UpdatedAt)

	return err
}

func (r *PostgresOrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	org := &model.Organization{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, created_at, updated_at
		FROM organizations WHERE id = $1
	`, id).Scan(&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *PostgresOrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	return err
}
