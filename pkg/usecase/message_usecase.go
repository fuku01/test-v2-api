package usecase

import (
	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
)

type MessageUsecase interface {
	ListMessages() ([]*domain_model.Message, error)
	CreateMessage(input *domain_model.CreateMessageInput) (*domain_model.Message, error)
}

type messageUsecase struct {
	tr repository.MessageRepository
}

func NewMessageUsecase(tr repository.MessageRepository) MessageUsecase {
	return &messageUsecase{
		tr: tr,
	}
}

func (u *messageUsecase) ListMessages() ([]*domain_model.Message, error) {
	msgs, err := u.tr.ListMessages()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (u *messageUsecase) CreateMessage(input *domain_model.CreateMessageInput) (*domain_model.Message, error) {
	msg, err := u.tr.CreateMessage(input)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
