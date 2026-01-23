package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{
		pool: pool,
	}
}

func (r *OrganizationRepository) Create(ctx context.Context, organization *model.Organization) error {
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
