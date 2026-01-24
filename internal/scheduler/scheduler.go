package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
)

type Scheduler struct {
	messageRepo repository.MessageRepository
	queue       *queue.Queue
}

func NewScheduler(messageRepo repository.MessageRepository,
	queue *queue.Queue,
) *Scheduler {
	return &Scheduler{
		messageRepo: messageRepo,
		queue:       queue,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := s.processRetries(ctx)
			if err != nil {
				log.Printf("[Scheduler] Error: %v", err)
			}
			time.Sleep(5 * time.Second) // Check every 5 seconds
		}
	}
}

func (s *Scheduler) processRetries(ctx context.Context) error {
	messages, err := s.messageRepo.FindRetryReady(ctx, 10)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		log.Printf("[Scheduler] Requeueing message %s (attempt %d)", msg.ID, msg.AttemptCount)
		s.queue.Push(ctx, msg.ID.String())

		msg.Status = "pending"
		s.messageRepo.Update(ctx, msg)
	}
	return nil
}
