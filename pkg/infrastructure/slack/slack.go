package chat

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fuku01/test-v2-api/pkg/domain/chat"
	"github.com/slack-go/slack"
)

type Slack struct {
	client *slack.Client
}

func NewSlack() (*Slack, error) {
	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("環境変数 SLACK_BOT_TOKEN が設定されていません")
	}
	client := slack.New(token)

	return &Slack{
		client: client,
	}, nil
}

func (s *Slack) PostMessage(ctx context.Context, req *chat.PostMessageRequest) (*chat.PostMessageResponse, error) {
	channel, postAt, err := s.client.PostMessage(req.ChannelID, slack.MsgOptionText(req.Message, false))
	if err != nil {
		return nil, err
	}

	// SlackのPostAtはUnixタイムスタンプで返ってくるので、time.Timeに変換する
	timestamp, err := strconv.ParseFloat(postAt, 64)
	if err != nil {
		return nil, err
	}
	timePostAt := time.Unix(int64(timestamp), 0)

	return &chat.PostMessageResponse{
		ChannelID: channel,
		Message:   req.Message,
		PostAt:    timePostAt,
	}, nil
}
