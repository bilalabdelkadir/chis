package repository

import (
	"context"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error)
}

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey *model.APIKey) error
	FindByHashedKey(ctx context.Context, hashedKey string) (*model.APIKey, error)
	FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]*model.APIKey, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeliveryAttemptRepository interface {
	Create(ctx context.Context, deliveryAttempt *model.DeliveryAttempt) error
}
type MembershipRepository interface {
	Create(ctx context.Context, membership *model.Membership) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
}

type MessageStats struct {
	TotalSent   int     `json:"totalWebhooksSent"`
	TotalFailed int     `json:"totalWebhooksFailed"`
	TotalQueued int     `json:"totalWebhooksQueued"`
	SuccessRate float64 `json:"successRate"`
}

type WebhookLogEntry struct {
	ID             uuid.UUID `json:"id"`
	Endpoint       string    `json:"endpoint"`
	Status         string    `json:"status"`
	StatusCode     int       `json:"statusCode"`
	EventType      string    `json:"eventType"`
	AttemptedAt    string    `json:"attemptedAt"`
	ResponseTimeMs int       `json:"responseTimeMs"`
}

type WebhookLogsResult struct {
	Data       []WebhookLogEntry `json:"data"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	FindPending(ctx context.Context, limit int) ([]*model.Message, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Message, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.Message, error)
	Update(ctx context.Context, msg *model.Message) error
	FindRetryReady(ctx context.Context, limit int) ([]*model.Message, error)
	GetStatsByOrgID(ctx context.Context, orgID uuid.UUID) (*MessageStats, error)
	FindWebhookLogs(ctx context.Context, orgID uuid.UUID, status string, search string, page int, limit int) (*WebhookLogsResult, error)
}

type OrganizationRepository interface {
	Create(ctx context.Context, organization *model.Organization) error
}
