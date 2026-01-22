package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"firstName"`
	LastName      *string   `json:"lastName"` // nullable
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Account struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"userId"`
	Provider          string     `json:"provider"`
	ProviderAccountID string     `json:"providerAccountId"`
	PasswordHash      *string    `json:"-"`
	AccessToken       *string    `json:"-"`
	RefreshToken      *string    `json:"-"`
	TokenExpiresAt    *time.Time `json:"tokenExpiresAt,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
