package chat

import (
	"context"
	"time"
)

// Chatインターフェースは、中身がSlackであるかどうかを知らない

type Chat interface {
	PostMessage(ctx context.Context, input *PostMessageRequest) (*PostMessageResponse, error)
}

type PostMessageRequest struct {
	Message   string
	ChannelID string
}
type PostMessageResponse struct {
	Message   string
	ChannelID string
	PostAt    time.Time
}
