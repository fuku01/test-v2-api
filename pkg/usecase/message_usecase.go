package usecase

import (
	"context"

	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
)

type MessageUsecase interface {
	ListMessages(ctx context.Context) ([]*domain_model.Message, error)
	CreateMessage(ctx context.Context, req *domain_model.CreateMessageRequest) (*domain_model.Message, error)
}

type messageUsecase struct {
	tr repository.MessageRepository
}

func NewMessageUsecase(tr repository.MessageRepository) MessageUsecase {
	return &messageUsecase{
		tr: tr,
	}
}

func (u *messageUsecase) ListMessages(ctx context.Context) ([]*domain_model.Message, error) {
	msgs, err := u.tr.ListMessages(ctx)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (u *messageUsecase) CreateMessage(ctx context.Context, req *domain_model.CreateMessageRequest) (*domain_model.Message, error) {
	msg, err := u.tr.CreateMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
