package chat

import "time"

// Chatインターフェースは、中身がSlackであるかどうかを知らない

type Chat interface {
	PostMessage(input *PostMessageRequest) (*PostMessageResponse, error)
}

type PostMessageRequest struct {
	ChannelID string
	Message   string
}
type PostMessageResponse struct {
	ChannelID string
	Message   string
	PostAt    time.Time
}
