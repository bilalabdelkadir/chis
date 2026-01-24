package worker

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
)

type Worker struct {
	messageRepo repository.MessageRepository
	attemptRepo repository.DeliveryAttemptRepository
}

func NewWorker(messageRepo repository.MessageRepository, attemptRepo repository.DeliveryAttemptRepository) *Worker {
	return &Worker{
		messageRepo: messageRepo,
		attemptRepo: attemptRepo,
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := w.messageRepo.FindPending(ctx, 10)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, message := range messages {
				w.deliver(ctx, message)
			}

			time.Sleep(5 * time.Second)
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
		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "success")
	} else {
		_, err = w.messageRepo.UpdateStatus(ctx, msg.ID, "failed")
	}
}
