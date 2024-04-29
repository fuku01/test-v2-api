package usecase

import (
	"context"

	"github.com/fuku01/test-v2-api/pkg/domain/chat"
	"github.com/fuku01/test-v2-api/pkg/domain/entity"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
)

type MessageUsecase interface {
	ListMessages(ctx context.Context) ([]*entity.Message, error)
	CreateMessage(ctx context.Context, req *entity.CreateMessageRequest) (*entity.Message, error)

	PostMessage(ctx context.Context, req *chat.PostMessageRequest) (*chat.PostMessageResponse, error)
}

type messageUsecase struct {
	tr   repository.MessageRepository
	chat chat.Chat
}

func NewMessageUsecase(tr repository.MessageRepository, chat chat.Chat) MessageUsecase {
	return &messageUsecase{
		tr:   tr,
		chat: chat,
	}
}

func (u *messageUsecase) ListMessages(ctx context.Context) ([]*entity.Message, error) {
	msgs, err := u.tr.ListMessages(ctx)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (u *messageUsecase) CreateMessage(ctx context.Context, req *entity.CreateMessageRequest) (*entity.Message, error) {
	msg, err := u.tr.CreateMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (u *messageUsecase) PostMessage(ctx context.Context, req *chat.PostMessageRequest) (*chat.PostMessageResponse, error) {
	res, err := u.chat.PostMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
