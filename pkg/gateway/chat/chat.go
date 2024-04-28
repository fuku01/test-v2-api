package chat

import "time"

// Chatインターフェースは、中身がSlackであるかどうかを知らない

type Chat interface {
	PostMessage(input *ChatMessageInput) (*ChatMessageResponse, error)
}

type ChatMessageInput struct {
	ChannelID string
	Message   string
}
type ChatMessageResponse struct {
	ChannelID string
	Message   string
	PostAt    time.Time
}
