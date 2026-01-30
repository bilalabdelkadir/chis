package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMembershipRepository struct {
	pool *pgxpool.Pool
}

func NewMembershipRepository(pool *pgxpool.Pool) MembershipRepository {
	return &PostgresMembershipRepository{
		pool: pool,
	}
}

func (r *PostgresMembershipRepository) Create(ctx context.Context, membership *model.Membership) error {
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

func (r *PostgresMembershipRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Membership, error) {
	m := &model.Membership{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, org_id, role, created_at, updated_at
		FROM memberships
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(
		&m.ID,
		&m.UserID,
		&m.OrgID,
		&m.Role,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *PostgresMembershipRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*MembershipWithOrg, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.org_id, o.name, o.slug, m.role
		FROM memberships m
		JOIN organizations o ON o.id = m.org_id
		WHERE m.user_id = $1
		ORDER BY m.created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*MembershipWithOrg
	for rows.Next() {
		m := &MembershipWithOrg{}
		if err := rows.Scan(&m.OrgID, &m.OrgName, &m.OrgSlug, &m.Role); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *PostgresMembershipRepository) FindByUserAndOrgID(ctx context.Context, userID, orgID uuid.UUID) (*model.Membership, error) {
	m := &model.Membership{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, org_id, role, created_at, updated_at
		FROM memberships
		WHERE user_id = $1 AND org_id = $2
	`, userID, orgID).Scan(
		&m.ID,
		&m.UserID,
		&m.OrgID,
		&m.Role,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}
