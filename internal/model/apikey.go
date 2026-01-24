package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID         uuid.UUID  `json:"id"`
	OrgID      uuid.UUID  `json:"orgId"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	HashedKey  string     `json:"-"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type Message struct {
	ID           uuid.UUID       `json:"id"`
	OrgID        uuid.UUID       `json:"orgId"`
	Method       string          `json:"method"` // e.g., "POST"
	URL          string          `json:"url"`
	Payload      json.RawMessage `json:"payload"` // JSONB stored as []byte
	Status       string          `json:"status"`  // 'pending', 'success', 'failed'
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	AttemptCount int             `json:"attemptCount"`
	NextRetryAt  *time.Time      `json:"nextRetryAt"`
}

type DeliveryAttempt struct {
	ID            uuid.UUID `json:"id"`
	MessageID     uuid.UUID `json:"messageId"`
	AttemptNumber int       `json:"attemptNumber"`
	StatusCode    *int      `json:"statusCode"`   // nullable but always present
	ResponseBody  *string   `json:"responseBody"` // nullable but always present
	ErrorMessage  *string   `json:"errorMessage"` // nullable but always present
	DurationMS    *int      `json:"durationMs"`   // nullable but always present
	AttemptedAt   time.Time `json:"attemptedAt"`
}
