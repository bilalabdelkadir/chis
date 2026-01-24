package delivery

import (
	"context"
	"log/slog"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bilalabdelkadir/chis/proto/delivery"
)

type ServiceRepo struct {
	messageRepo repository.MessageRepository
	queue       *queue.Queue
	pb.UnimplementedDeliveryServiceServer
}

func NewDeliveryService(messageRepo repository.MessageRepository, queue *queue.Queue) *ServiceRepo {
	return &ServiceRepo{
		messageRepo: messageRepo,
		queue:       queue,
	}
}

func (s *ServiceRepo) QueueMessage(ctx context.Context, req *pb.QueueMessageRequest) (*pb.QueueMessageResponse, error) {
	slog.Info("message_received", "url", req.Url, "method", req.Method.String())

	orgId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid org_id")
	}
	method := req.Method.String()

	message := &model.Message{
		OrgID:   orgId,
		Method:  method,
		URL:     req.Url,
		Payload: req.Payload,
	}

	err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save message")
	}
	slog.Info("message_saved", "message_id", message.ID, "org_id", message.OrgID)

	s.queue.Push(ctx, message.ID.String())
	slog.Info("message_queued", "message_id", message.ID, "org_id", message.OrgID)

	res := pb.QueueMessageResponse{
		MessageId: message.ID.String(),
		Status:    message.Status,
	}

	return &res, nil

}
