package worker

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/google/uuid"
)

type Worker struct {
	messageRepo repository.MessageRepository
	attemptRepo repository.DeliveryAttemptRepository
	queue       *queue.Queue
}

func NewWorker(messageRepo repository.MessageRepository, attemptRepo repository.DeliveryAttemptRepository,
	queue *queue.Queue,
) *Worker {
	return &Worker{
		messageRepo: messageRepo,
		attemptRepo: attemptRepo,
		queue:       queue,
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

	client := &http.Client{Timeout: 10 * time.Second}
	log.Printf("[Worker] Delivering message %s to %s", msg.ID, msg.URL)

	start := time.Now()
	resp, err := client.Do(req)
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
		msg := err.Error()
		errorMessage = &msg
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
		AttemptNumber: 1,
		StatusCode:    statusCode,
		ErrorMessage:  errorMessage,
		DurationMS:    durationMS,
		ResponseBody:  responseBody,
	}

	_ = w.attemptRepo.Create(ctx, attempt)

	if success {
		log.Printf("[Worker] Success: %s (status=%d, duration=%dms)", msg.ID, *statusCode, *durationMS)

		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "success")
	} else {
		log.Printf("[Worker] Failed: %s (error=%s)", msg.ID, *errorMessage)

		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "failed")
	}
}
