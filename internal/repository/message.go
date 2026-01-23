package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		pool: pool,
	}
}

func (r *MessageRepository) Create(ctx context.Context, message *model.Message) error {
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
