package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryAttemptsRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryAttemptsRepository(pool *pgxpool.Pool) *DeliveryAttemptsRepository {
	return &DeliveryAttemptsRepository{
		pool: pool,
	}
}

func (r *DeliveryAttemptsRepository) Create(ctx context.Context, deliveryAttempt *model.DeliveryAttempt) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO delivery_attempts (message_id, attempt_number,status_code,response_body, error_message, duration_ms)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, attempted_at

	`,
		deliveryAttempt.MessageID,
		deliveryAttempt.AttemptNumber,
		deliveryAttempt.StatusCode,
		deliveryAttempt.ResponseBody,
		deliveryAttempt.ErrorMessage,
		deliveryAttempt.DurationMS,
	).Scan(&deliveryAttempt.ID, &deliveryAttempt.AttemptedAt)

	return err
}
