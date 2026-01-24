package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &PostgresMessageRepository{
		pool: pool,
	}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, message *model.Message) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO messages (org_id, method, url,payload)
		VALUES ($1,$2,$3,$4)
		RETURNING id, status,created_at, updated_at

	`,
		message.OrgID,
		message.Method,
		message.URL,
		message.Payload,
	).Scan(&message.ID, &message.Status, &message.CreatedAt, &message.UpdatedAt)

	return err
}

func (r *PostgresMessageRepository) FindPending(ctx context.Context, limit int) ([]*model.Message, error) {
	messages := []*model.Message{}

	rows, err := r.pool.Query(ctx, `
		SELECT id, org_id, method, url, payload, status, created_at, updated_at
		FROM messages
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through results
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.OrgID,
			&msg.Method,
			&msg.URL,
			&msg.Payload,
			&msg.Status,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *PostgresMessageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Message, error) {
	msg := &model.Message{}

	err := r.pool.QueryRow(ctx, `
		UPDATE messages
		SET status = $2
		WHERE id = $1
		RETURNING id, org_id, method, url, payload, status, created_at, updated_at
	`, id, status).Scan(
		&msg.ID,
		&msg.OrgID,
		&msg.Method,
		&msg.URL,
		&msg.Payload,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *PostgresMessageRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	var msg model.Message

	err := r.pool.QueryRow(ctx, `
		SELECT id, org_id, method, url, payload, status, created_at, updated_at
		FROM messages
		WHERE id = $1
	`, id).Scan(
		&msg.ID,
		&msg.OrgID,
		&msg.Method,
		&msg.URL,
		&msg.Payload,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &msg, nil
}
