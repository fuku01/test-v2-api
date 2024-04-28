package repository

import (
	"context"

	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
)

type MessageRepository interface {
	ListMessages(ctx context.Context) ([]*domain_model.Message, error)
	CreateMessage(ctx context.Context, req *domain_model.CreateMessageRequest) (*domain_model.Message, error)
}
