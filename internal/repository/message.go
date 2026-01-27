package repository

import (
	"context"
	"fmt"
	"math"
	"time"

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
		SELECT id, org_id, method, url, payload, status, created_at, updated_at, attempt_count, next_retry_at
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
		&msg.AttemptCount,
		&msg.NextRetryAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &msg, nil
}

func (r *PostgresMessageRepository) Update(ctx context.Context, msg *model.Message) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE messages 
        SET status = $1, attempt_count = $2, next_retry_at = $3
        WHERE id = $4
    `, msg.Status, msg.AttemptCount, msg.NextRetryAt, msg.ID)
	return err
}

func (r *PostgresMessageRepository) FindRetryReady(ctx context.Context, limit int) ([]*model.Message, error) {

	rows, err := r.pool.Query(ctx, `
        SELECT id, org_id, method, url, payload, status,
               created_at, updated_at, attempt_count, next_retry_at
        FROM messages
        WHERE status = 'retry'
          AND next_retry_at IS NOT NULL
          AND next_retry_at <= NOW()
        ORDER BY next_retry_at ASC
        LIMIT $1
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message

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
			&msg.AttemptCount,
			&msg.NextRetryAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *PostgresMessageRepository) GetStatsByOrgID(ctx context.Context, orgID uuid.UUID) (*MessageStats, error) {
	var total, sent, failed, queued int

	err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'success') AS sent,
			COUNT(*) FILTER (WHERE status = 'failed') AS failed,
			COUNT(*) FILTER (WHERE status = 'pending' OR status = 'retry') AS queued
		FROM messages
		WHERE org_id = $1
	`, orgID).Scan(&total, &sent, &failed, &queued)
	if err != nil {
		return nil, err
	}

	var successRate float64
	if total > 0 {
		successRate = math.Round(float64(sent)/float64(total)*10000) / 100
	}

	return &MessageStats{
		TotalSent:   sent,
		TotalFailed: failed,
		TotalQueued: queued,
		SuccessRate: successRate,
	}, nil
}

func (r *PostgresMessageRepository) FindWebhookLogs(ctx context.Context, orgID uuid.UUID, status string, search string, page int, limit int) (*WebhookLogsResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build WHERE clause
	where := "WHERE m.org_id = $1"
	args := []any{orgID}
	argIdx := 2

	if status != "" {
		where += fmt.Sprintf(" AND m.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if search != "" {
		where += fmt.Sprintf(" AND m.url ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM messages m %s", where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	// Fetch page
	dataQuery := fmt.Sprintf(`
		SELECT m.id, m.url, m.status, m.method, m.created_at,
			COALESCE(da.status_code, 0),
			COALESCE(da.duration_ms, 0),
			COALESCE(da.attempted_at, m.created_at)
		FROM messages m
		LEFT JOIN LATERAL (
			SELECT status_code, duration_ms, attempted_at
			FROM delivery_attempts
			WHERE message_id = m.id
			ORDER BY attempt_number DESC
			LIMIT 1
		) da ON true
		%s
		ORDER BY m.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []WebhookLogEntry
	for rows.Next() {
		var (
			id          uuid.UUID
			url, msgStatus, method string
			createdAt, attemptedAt time.Time
			statusCode, durationMs int
		)
		if err := rows.Scan(&id, &url, &msgStatus, &method, &createdAt, &statusCode, &durationMs, &attemptedAt); err != nil {
			return nil, err
		}
		entries = append(entries, WebhookLogEntry{
			ID:             id,
			Endpoint:       url,
			Status:         msgStatus,
			StatusCode:     statusCode,
			EventType:      method,
			AttemptedAt:    attemptedAt.Format(time.RFC3339),
			ResponseTimeMs: durationMs,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if entries == nil {
		entries = []WebhookLogEntry{}
	}

	return &WebhookLogsResult{
		Data:       entries,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
