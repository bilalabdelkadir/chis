package repository

import (
	"context"
	"time"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInvitationRepository struct {
	pool *pgxpool.Pool
}

func NewInvitationRepository(pool *pgxpool.Pool) InvitationRepository {
	return &PostgresInvitationRepository{
		pool: pool,
	}
}

func (r *PostgresInvitationRepository) Create(ctx context.Context, invitation *model.Invitation) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO invitations (org_id, email, role, token, status, invited_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`,
		invitation.OrgID,
		invitation.Email,
		invitation.Role,
		invitation.Token,
		invitation.Status,
		invitation.InvitedBy,
		invitation.ExpiresAt,
	).Scan(&invitation.ID, &invitation.CreatedAt, &invitation.UpdatedAt)

	return err
}

func (r *PostgresInvitationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Invitation, error) {
	inv := &model.Invitation{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, org_id, email, role, token, status, invited_by, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE id = $1
	`, id).Scan(
		&inv.ID,
		&inv.OrgID,
		&inv.Email,
		&inv.Role,
		&inv.Token,
		&inv.Status,
		&inv.InvitedBy,
		&inv.ExpiresAt,
		&inv.AcceptedAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return inv, nil
}

func (r *PostgresInvitationRepository) FindPendingByEmail(ctx context.Context, email string) ([]PendingInvitation, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT i.id, o.name, o.slug, i.role, i.expires_at
		FROM invitations i
		JOIN organizations o ON o.id = i.org_id
		WHERE i.email = $1 AND i.status = 'pending' AND i.expires_at > NOW()
		ORDER BY i.created_at DESC
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PendingInvitation
	for rows.Next() {
		var p PendingInvitation
		if err := rows.Scan(&p.ID, &p.OrgName, &p.OrgSlug, &p.Role, &p.ExpiresAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, rows.Err()
}

func (r *PostgresInvitationRepository) FindByToken(ctx context.Context, token string) (*model.Invitation, error) {
	inv := &model.Invitation{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, org_id, email, role, token, status, invited_by, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE token = $1
	`, token).Scan(
		&inv.ID,
		&inv.OrgID,
		&inv.Email,
		&inv.Role,
		&inv.Token,
		&inv.Status,
		&inv.InvitedBy,
		&inv.ExpiresAt,
		&inv.AcceptedAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return inv, nil
}

func (r *PostgresInvitationRepository) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]*model.Invitation, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, org_id, email, role, token, status, invited_by, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE org_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Invitation
	for rows.Next() {
		inv := &model.Invitation{}
		if err := rows.Scan(
			&inv.ID,
			&inv.OrgID,
			&inv.Email,
			&inv.Role,
			&inv.Token,
			&inv.Status,
			&inv.InvitedBy,
			&inv.ExpiresAt,
			&inv.AcceptedAt,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, inv)
	}

	return result, rows.Err()
}

func (r *PostgresInvitationRepository) FindByOrgAndEmail(ctx context.Context, orgID uuid.UUID, email string) (*model.Invitation, error) {
	inv := &model.Invitation{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, org_id, email, role, token, status, invited_by, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE org_id = $1 AND email = $2
	`, orgID, email).Scan(
		&inv.ID,
		&inv.OrgID,
		&inv.Email,
		&inv.Role,
		&inv.Token,
		&inv.Status,
		&inv.InvitedBy,
		&inv.ExpiresAt,
		&inv.AcceptedAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return inv, nil
}

func (r *PostgresInvitationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE invitations
		SET status = $1
		WHERE id = $2
	`, status, id)

	return err
}

func (r *PostgresInvitationRepository) AcceptInvitation(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `
		UPDATE invitations
		SET status = 'accepted', accepted_at = $1
		WHERE id = $2
	`, now, id)

	return err
}

func (r *PostgresInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM invitations
		WHERE id = $1
	`, id)

	return err
}
