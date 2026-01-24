package handler

import (
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/helper"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	"github.com/google/uuid"
)

type ApiKeyHandler struct {
	membershipRepo repository.MembershipRepository
	apiKeyRepo     repository.ApiKeyRepository
}

func NewApiKeyHandler(
	membershipRepo repository.MembershipRepository,
	apiKeyRepo repository.ApiKeyRepository,

) *ApiKeyHandler {
	return &ApiKeyHandler{
		membershipRepo: membershipRepo,
		apiKeyRepo:     apiKeyRepo,
	}
}

type CreateApiKeyRequest struct {
	Name      string    `json:"name" validate:"required,min=3"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type CreateApiKeyResponse struct {
	ApiKey string `json:"apiKey"`
}

func (h *ApiKeyHandler) Create(w http.ResponseWriter, r *http.Request) error {

	var req CreateApiKeyRequest

	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	userIdValue := r.Context().Value(middleware.UserIDKey)
	if userIdValue == nil {
		return apperror.Unauthorized("user not found in context")
	}
	userId, ok := userIdValue.(uuid.UUID)
	if !ok {
		return apperror.Unauthorized("invalid user ID type")
	}

	membership, err := h.membershipRepo.FindByUserID(r.Context(), userId)
	if err != nil {
		return apperror.NotFound("membership not found")
	}

	fullApiKey, err := helper.GenerateRandomApiKey("chis_sk_")
	if err != nil {
		return apperror.Internal("couldn't generate api key")
	}

	hashKey := helper.HashSHA256(fullApiKey)
	displayPrefix := fullApiKey[:12]

	var expiresAt *time.Time
	if !req.ExpiresAt.IsZero() {
		expiresAt = &req.ExpiresAt
	}

	apiKey := &model.APIKey{
		OrgID:     membership.OrgID,
		Name:      req.Name,
		Prefix:    displayPrefix,
		HashedKey: hashKey,
		ExpiresAt: expiresAt,
	}

	err = h.apiKeyRepo.Create(r.Context(), apiKey)
	if err != nil {
		return err
	}

	res := CreateApiKeyResponse{
		ApiKey: fullApiKey,
	}

	response.WriteJSON(w, http.StatusCreated, res)
	return nil

}
