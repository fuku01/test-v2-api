package repository

import (
	"context"

	gorm_model "github.com/fuku01/test-v2-api/db/model"
	"github.com/fuku01/test-v2-api/pkg/domain/entity"
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

func (r *messageRepository) ListMessages(ctx context.Context) ([]*entity.Message, error) {

	msgs := []*gorm_model.Message{}
	err := r.db.Find(&msgs).Error // 論理削除されたデータは取得しない
	if err != nil {
		return nil, err
	}

	convMsgs := lo.Map(msgs, func(msg *gorm_model.Message, _ int) *entity.Message {
		return r.convMessage(msg)
	})

	return convMsgs, nil
}

func (r *messageRepository) CreateMessage(ctx context.Context, req *entity.CreateMessageRequest) (*entity.Message, error) {
	msg := &entity.Message{
		Content: req.Content,
	}

	err := r.db.Create(msg).Error
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// gormの型をドメインモデルの型に変換
func (r *messageRepository) convMessage(msg *gorm_model.Message) *entity.Message {
	if msg == nil {
		return nil
	}

	return &entity.Message{
		ID:        msg.ID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
}
