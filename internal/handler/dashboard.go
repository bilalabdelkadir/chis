package handler

import (
	"net/http"
	"strconv"

	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	membershipRepo repository.MembershipRepository
	messageRepo    repository.MessageRepository
}

func NewDashboardHandler(
	membershipRepo repository.MembershipRepository,
	messageRepo repository.MessageRepository,
) *DashboardHandler {
	return &DashboardHandler{
		membershipRepo: membershipRepo,
		messageRepo:    messageRepo,
	}
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	membership, err := h.membershipRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		return apperror.NotFound("membership not found")
	}

	stats, err := h.messageRepo.GetStatsByOrgID(r.Context(), membership.OrgID)
	if err != nil {
		return apperror.Internal("failed to fetch dashboard stats")
	}

	response.WriteJSON(w, http.StatusOK, stats)
	return nil
}

func (h *DashboardHandler) WebhookLogs(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	membership, err := h.membershipRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		return apperror.NotFound("membership not found")
	}

	query := r.URL.Query()
	status := query.Get("status")
	search := query.Get("search")

	page := 1
	if v := query.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if v := query.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	result, err := h.messageRepo.FindWebhookLogs(r.Context(), membership.OrgID, status, search, page, limit)
	if err != nil {
		return apperror.Internal("failed to fetch webhook logs")
	}

	response.WriteJSON(w, http.StatusOK, result)
	return nil
}

func extractUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(middleware.UserIDKey)
	if val == nil {
		return uuid.Nil, apperror.Unauthorized("user not found in context")
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperror.Unauthorized("invalid user ID type")
	}
	return userID, nil
}
