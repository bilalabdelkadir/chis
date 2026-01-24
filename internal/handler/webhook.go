package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	"github.com/google/uuid"
)

type WebhookHandler struct {
	messageRepo repository.MessageRepository
	queue       *queue.Queue
}

type SendWebhookRequest struct {
	URL     string      `json:"url" validate:"required,url"`
	Method  string      `json:"method"` // optional, default POST
	Payload interface{} `json:"payload" validate:"required"`
}

type SendWebhookResponse struct {
	MessageID uuid.UUID `json:"messageId"`
	Status    string    `json:"status"`
}

func NewWebhookHandler(
	messageRepo repository.MessageRepository,
	queue *queue.Queue,
) *WebhookHandler {
	return &WebhookHandler{
		messageRepo: messageRepo,
		queue:       queue,
	}
}

func (h *WebhookHandler) Send(w http.ResponseWriter, r *http.Request) error {
	var req SendWebhookRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	orgIdValue := r.Context().Value(middleware.OrgIDKey)
	if orgIdValue == nil {
		return apperror.Unauthorized("org not found in context")
	}
	orgId, ok := orgIdValue.(uuid.UUID)
	if !ok {
		return apperror.Unauthorized("invalid Org ID type")
	}

	method := req.Method
	if method == "" {
		method = "POST"
	}

	payload, err := json.Marshal(req.Payload)

	if err != nil {
		return err
	}

	message := &model.Message{
		OrgID:   orgId,
		Method:  method,
		URL:     req.URL,
		Payload: payload,
	}

	err = h.messageRepo.Create(r.Context(), message)
	if err != nil {
		return err
	}

	h.queue.Push(r.Context(), message.ID.String())

	res := SendWebhookResponse{
		MessageID: message.ID,
		Status:    message.Status,
	}

	response.WriteJSON(w, http.StatusCreated, res)

	return nil

}
