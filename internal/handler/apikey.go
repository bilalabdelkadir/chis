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
	"github.com/go-chi/chi/v5"
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
	ID        string `json:"id"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt"`
}

type ListApiKeyItem struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Prefix     string  `json:"prefix"`
	CreatedAt  string  `json:"createdAt"`
	LastUsedAt *string `json:"lastUsedAt"`
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
		ID:        apiKey.ID.String(),
		Name:      apiKey.Name,
		Prefix:    apiKey.Prefix,
		Key:       fullApiKey,
		CreatedAt: apiKey.CreatedAt.Format(time.RFC3339),
	}

	response.WriteJSON(w, http.StatusCreated, res)
	return nil
}

func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	membership, err := h.membershipRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		return apperror.NotFound("membership not found")
	}

	keys, err := h.apiKeyRepo.FindByOrgID(r.Context(), membership.OrgID)
	if err != nil {
		return apperror.Internal("failed to fetch api keys")
	}

	items := make([]ListApiKeyItem, 0, len(keys))
	for _, k := range keys {
		item := ListApiKeyItem{
			ID:        k.ID.String(),
			Name:      k.Name,
			Prefix:    k.Prefix,
			CreatedAt: k.CreatedAt.Format(time.RFC3339),
		}
		if k.LastUsedAt != nil {
			formatted := k.LastUsedAt.Format(time.RFC3339)
			item.LastUsedAt = &formatted
		}
		items = append(items, item)
	}

	response.WriteJSON(w, http.StatusOK, items)
	return nil
}

func (h *ApiKeyHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")
	keyID, err := uuid.Parse(idParam)
	if err != nil {
		return apperror.BadRequest("invalid api key id")
	}

	if err := h.apiKeyRepo.Delete(r.Context(), keyID); err != nil {
		if err == repository.ErrNotFound {
			return apperror.NotFound("api key not found")
		}
		return apperror.Internal("failed to delete api key")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
