package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipRepository struct {
	pool *pgxpool.Pool
}

func NewMembershipRepository(pool *pgxpool.Pool) *MembershipRepository {
	return &MembershipRepository{
		pool: pool,
	}
}

func (r *MembershipRepository) Create(ctx context.Context, membership *model.Membership) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO memberships (user_id, org_id, role)
		VALUES ($1,$2,$3)
		RETURNING id, created_at, updated_at

	`,
		membership.UserID,
		membership.OrgID,
		membership.Role,
	).Scan(&membership.ID, &membership.CreatedAt, &membership.UpdatedAt)

	return err
}
