package handler

import (
	"net/http"
	"strconv"

	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	messageRepo         repository.MessageRepository
	deliveryAttemptRepo repository.DeliveryAttemptRepository
}

func NewDashboardHandler(
	messageRepo repository.MessageRepository,
	deliveryAttemptRepo repository.DeliveryAttemptRepository,
) *DashboardHandler {
	return &DashboardHandler{
		messageRepo:         messageRepo,
		deliveryAttemptRepo: deliveryAttemptRepo,
	}
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) error {
	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	stats, err := h.messageRepo.GetStatsByOrgID(r.Context(), orgID)
	if err != nil {
		return apperror.Internal("failed to fetch dashboard stats")
	}

	response.WriteJSON(w, http.StatusOK, stats)
	return nil
}

func (h *DashboardHandler) WebhookLogs(w http.ResponseWriter, r *http.Request) error {
	orgID, err := extractOrgID(r)
	if err != nil {
		return err
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

	result, err := h.messageRepo.FindWebhookLogs(r.Context(), orgID, status, search, page, limit)
	if err != nil {
		return apperror.Internal("failed to fetch webhook logs")
	}

	response.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *DashboardHandler) WebhookLogDetail(w http.ResponseWriter, r *http.Request) error {
	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	idParam := chi.URLParam(r, "id")
	msgID, err := uuid.Parse(idParam)
	if err != nil {
		return apperror.BadRequest("invalid log ID")
	}

	msg, err := h.messageRepo.FindById(r.Context(), msgID)
	if err != nil {
		return apperror.NotFound("webhook log not found")
	}

	if msg.OrgID != orgID {
		return apperror.NotFound("webhook log not found")
	}

	attempts, err := h.deliveryAttemptRepo.FindByMessageID(r.Context(), msgID)
	if err != nil {
		return apperror.Internal("failed to fetch delivery attempts")
	}

	attemptDetails := make([]repository.DeliveryAttemptDetail, len(attempts))
	for i, a := range attempts {
		attemptDetails[i] = repository.DeliveryAttemptDetail{
			ID:            a.ID,
			AttemptNumber: a.AttemptNumber,
			StatusCode:    a.StatusCode,
			ResponseBody:  a.ResponseBody,
			ErrorMessage:  a.ErrorMessage,
			DurationMS:    a.DurationMS,
			AttemptedAt:   a.AttemptedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	var nextRetry *string
	if msg.NextRetryAt != nil {
		t := msg.NextRetryAt.Format("2006-01-02T15:04:05Z")
		nextRetry = &t
	}

	detail := repository.WebhookLogDetail{
		ID:               msg.ID,
		Method:           msg.Method,
		URL:              msg.URL,
		Status:           msg.Status,
		Payload:          msg.Payload,
		AttemptCount:     msg.AttemptCount,
		CreatedAt:        msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        msg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		NextRetryAt:      nextRetry,
		DeliveryAttempts: attemptDetails,
	}

	response.WriteJSON(w, http.StatusOK, detail)
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

func extractOrgID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(middleware.OrgIDKey)
	if val == nil {
		return uuid.Nil, apperror.BadRequest("organization not found in context")
	}
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperror.BadRequest("invalid organization ID type")
	}
	return orgID, nil
}
