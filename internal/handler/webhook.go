package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	pb "github.com/bilalabdelkadir/chis/proto/delivery"
	"github.com/google/uuid"
)

type WebhookHandler struct {
	grpcClient pb.DeliveryServiceClient
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
	grpcClient pb.DeliveryServiceClient,
) *WebhookHandler {
	return &WebhookHandler{
		grpcClient: grpcClient,
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

	v, ok := pb.HttpMethod_value[method]
	if !ok {
		return fmt.Errorf("invalid http method: %s", method)
	}

	methodEnum := pb.HttpMethod(v)

	log.Printf("[API] Received webhook request for URL: %s", req.URL)

	grpcReq := &pb.QueueMessageRequest{
		Url:     req.URL,
		Method:  methodEnum,
		Payload: payload,
		OrgId:   orgId.String(),
	}

	grpcRes, err := h.grpcClient.QueueMessage(r.Context(), grpcReq)
	if err != nil {
		return err
	}
	log.Printf("[API] Delivery service returned message_id: %s", grpcRes.MessageId)

	msgId, err := uuid.Parse(grpcRes.MessageId)
	if err != nil {
		return err
	}

	res := SendWebhookResponse{
		MessageID: msgId,
		Status:    grpcRes.Status,
	}

	response.WriteJSON(w, http.StatusCreated, res)

	return nil

}
