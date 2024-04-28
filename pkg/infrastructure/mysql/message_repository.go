package infrastructure

import (
	"context"

	gorm_model "github.com/fuku01/test-v2-api/db/model"
	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) ListMessages(ctx context.Context) ([]*domain_model.Message, error) {

	msgs := []*gorm_model.Message{}
	err := r.db.Find(&msgs).Error // 論理削除されたデータは取得しない
	if err != nil {
		return nil, err
	}

	convMsgs := lo.Map(msgs, func(msg *gorm_model.Message, _ int) *domain_model.Message {
		return r.convMessage(msg)
	})

	return convMsgs, nil
}

func (r *messageRepository) CreateMessage(ctx context.Context, req *domain_model.CreateMessageRequest) (*domain_model.Message, error) {
	msg := &domain_model.Message{
		Content: req.Content,
	}

	err := r.db.Create(msg).Error
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// gormの型をドメインモデルの型に変換
func (r *messageRepository) convMessage(msg *gorm_model.Message) *domain_model.Message {
	if msg == nil {
		return nil
	}

	return &domain_model.Message{
		ID:        msg.ID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
}
