package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fuku01/test-v2-api/pkg/context/logger"
	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/gateway/chat"
	"github.com/fuku01/test-v2-api/pkg/graph/generated/model"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/samber/lo"
)

type MessageHandler interface {
	ListMessages(ctx context.Context) ([]*model.Message, error)
	PostMessage(ctx context.Context, req model.PostMessageInput) (*model.PostMessagePayload, error)
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
		logger.Error("ListMessages", err)
		return nil, InternalServerError
	}

	convMegs := lo.Map(msgs, func(m *domain_model.Message, _ int) *model.Message {
		return convMessage(m)
	})

	return convMegs, nil
}

// Slackにメッセージを投稿する処理
func (h *messageHandler) PostMessage(ctx context.Context, req model.PostMessageInput) (*model.PostMessagePayload, error) {
	if req.Message == "" || req.ChannelID == "" {
		return nil, InvalidRequest
	}

	input := &chat.PostMessageRequest{
		Message:   req.Message,
		ChannelID: req.ChannelID,
	}

	res, err := h.tu.PostMessage(ctx, input)
	if err != nil {
		logger.Error("PostMessage", err)
		return nil, InternalServerError
	}

	convRes := convPostMessageResponse(res)

	return convRes, nil
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

func convPostMessageResponse(res *chat.PostMessageResponse) *model.PostMessagePayload {
	if res == nil {
		return nil
	}

	return &model.PostMessagePayload{
		Message:   res.Message,
		ChannelID: res.ChannelID,
		PostAt:    res.PostAt,
	}
}
