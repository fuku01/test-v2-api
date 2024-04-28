package mock

import (
	"context"

	"github.com/fuku01/test-v2-api/pkg/gateway/chat"
)

type MockChat interface {
	PostMessage(ctx context.Context, input *chat.PostMessageRequest) (*chat.PostMessageResponse, error)
}

type mockChat struct{}

func NewMockChat() MockChat {
	return &mockChat{}
}

func (mock *mockChat) PostMessage(ctx context.Context, input *chat.PostMessageRequest) (*chat.PostMessageResponse, error) {
	return &chat.PostMessageResponse{}, nil
}
