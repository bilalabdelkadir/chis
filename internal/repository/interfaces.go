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
}

type DeliveryAttemptRepository interface {
	Create(ctx context.Context, deliveryAttempt *model.DeliveryAttempt) error
}
type MembershipRepository interface {
	Create(ctx context.Context, membership *model.Membership) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Membership, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	FindPending(ctx context.Context, limit int) ([]*model.Message, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Message, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.Message, error)
	Update(ctx context.Context, msg *model.Message) error
	FindRetryReady(ctx context.Context, limit int) ([]*model.Message, error)
}

type OrganizationRepository interface {
	Create(ctx context.Context, organization *model.Organization) error
}
