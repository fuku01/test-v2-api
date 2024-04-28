package repository

import domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"

type MessageRepository interface {
	ListMessages() ([]*domain_model.Message, error)
	CreateMessage(input *domain_model.CreateMessageInput) (*domain_model.Message, error)
}
