package worker

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/google/uuid"
)

const (
	MaxAttempts = 5
	BaseBackoff = 1 * time.Second
)

type Worker struct {
	messageRepo repository.MessageRepository
	attemptRepo repository.DeliveryAttemptRepository
	queue       *queue.Queue
	httpClient  *http.Client
}

func NewWorker(messageRepo repository.MessageRepository, attemptRepo repository.DeliveryAttemptRepository,
	queue *queue.Queue,
) *Worker {
	return &Worker{
		messageRepo: messageRepo,
		attemptRepo: attemptRepo,
		queue:       queue,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageID, err := w.queue.Pop(ctx)
			if err != nil {
				continue
			}
			id, err := uuid.Parse(messageID)
			if err != nil {
				continue
			}

			message, err := w.messageRepo.FindById(ctx, id)
			if err != nil {
				continue
			}

			w.deliver(ctx, message)

		}
	}
}

func (w *Worker) deliver(ctx context.Context, msg *model.Message) {
	req, err := http.NewRequestWithContext(
		ctx,
		msg.Method,
		msg.URL,
		bytes.NewReader(msg.Payload),
	)
	if err != nil {
		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "failed")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	slog.Info("webhook_delivering", "message_id", msg.ID, "org_id", msg.OrgID, "url", msg.URL)

	start := time.Now()
	resp, err := w.httpClient.Do(req)
	duration := time.Since(start)

	var (
		statusCode   *int
		errorMessage *string
		responseBody *string
		durationMS   *int
		success      bool
	)

	ms := int(duration.Milliseconds())
	durationMS = &ms

	if err != nil {
		errMsg := err.Error()
		errorMessage = &errMsg
		success = false
	} else {
		defer resp.Body.Close()

		code := resp.StatusCode
		statusCode = &code

		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			bodyStr := string(body)
			responseBody = &bodyStr
		}

		success = code >= 200 && code < 300
	}

	attempt := &model.DeliveryAttempt{
		MessageID:     msg.ID,
		AttemptNumber: msg.AttemptCount + 1,
		StatusCode:    statusCode,
		ErrorMessage:  errorMessage,
		DurationMS:    durationMS,
		ResponseBody:  responseBody,
	}

	_ = w.attemptRepo.Create(ctx, attempt)

	if success {
		slog.Info("webhook_delivered", "message_id", msg.ID, "org_id", msg.OrgID, "status_code", *statusCode, "duration_ms", *durationMS)
		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "success")
	} else {
		if errorMessage != nil {
			slog.Warn("webhook_failed", "message_id", msg.ID, "org_id", msg.OrgID, "error", *errorMessage)
		} else {
			slog.Warn("webhook_failed", "message_id", msg.ID, "org_id", msg.OrgID, "status_code", *statusCode)
		}
		if msg.AttemptCount < MaxAttempts {
			backoff := BaseBackoff * time.Duration(1<<msg.AttemptCount) // 1<<n is 2^n
			nextRetry := time.Now().Add(backoff)
			updatedData := &model.Message{
				ID:           msg.ID,
				AttemptCount: msg.AttemptCount + 1,
				NextRetryAt:  &nextRetry,
				Status:       "retry",
			}
			err = w.messageRepo.Update(ctx, updatedData)

		} else {
			msg.Status = "failed"
			err = w.messageRepo.Update(ctx, msg)
		}
	}
}
