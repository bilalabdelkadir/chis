package delivery

import (
	"context"
	"log"

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
	log.Printf("[Delivery] Received QueueMessage request for URL: %s", req.Url)
	// org_id comes from req.OrgId
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
	log.Printf("[Delivery] Saved message %s to database", message.ID)
	s.queue.Push(ctx, message.ID.String())
	log.Printf("[Delivery] Pushed message %s to queue", message.ID)

	res := pb.QueueMessageResponse{
		MessageId: message.ID.String(),
		Status:    message.Status,
	}

	return &res, nil

}
