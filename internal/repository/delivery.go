package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDeliveryAttemptsRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryAttemptsRepository(pool *pgxpool.Pool) DeliveryAttemptRepository {
	return &PostgresDeliveryAttemptsRepository{
		pool: pool,
	}
}

func (r *PostgresDeliveryAttemptsRepository) FindByMessageID(ctx context.Context, messageID uuid.UUID) ([]*model.DeliveryAttempt, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, message_id, attempt_number, status_code, response_body, error_message, duration_ms, attempted_at
		FROM delivery_attempts
		WHERE message_id = $1
		ORDER BY attempt_number ASC
	`, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*model.DeliveryAttempt
	for rows.Next() {
		a := &model.DeliveryAttempt{}
		err := rows.Scan(&a.ID, &a.MessageID, &a.AttemptNumber, &a.StatusCode, &a.ResponseBody, &a.ErrorMessage, &a.DurationMS, &a.AttemptedAt)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}

	return attempts, rows.Err()
}

func (r *PostgresDeliveryAttemptsRepository) Create(ctx context.Context, deliveryAttempt *model.DeliveryAttempt) error {
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
