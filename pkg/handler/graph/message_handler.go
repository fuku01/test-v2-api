package graph

import (
	"context"
	"fmt"
	"strconv"

	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/graph/generated/model"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/samber/lo"
)

type MessageHandler interface {
	ListMessages(ctx context.Context) ([]*model.Message, error)
}

type messageHandler struct {
	tu usecase.MessageUsecase
}

func NewMessageHandler(tu usecase.MessageUsecase) MessageHandler {
	return &messageHandler{
		tu: tu,
	}
}

func (h *messageHandler) ListMessages(ctx context.Context) ([]*model.Message, error) {
	fmt.Println("========================ListMessages()が呼ばれました==============================")

	msgs, err := h.tu.ListMessages(ctx)
	if err != nil {
		return nil, err
	}

	convMegs := lo.Map(msgs, func(m *domain_model.Message, _ int) *model.Message {
		return convMessage(m)
	})

	return convMegs, nil
}

// ドメインモデルの型をGraphQLの型に変換
func convMessage(msg *domain_model.Message) *model.Message {
	if msg == nil {
		return nil
	}

	return &model.Message{
		ID:        strconv.FormatUint(uint64(msg.ID), 10),
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
}
