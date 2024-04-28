package chat

import (
	"time"

	"github.com/slack-go/slack"
)

type Slack struct {
	client *slack.Client
}

func NewSlack(slackToken string) (*Slack, error) {
	client := slack.New(slackToken)
	return &Slack{
		client: client,
	}, nil
}

func (s *Slack) PostMessage(input *ChatMessageInput) (*ChatMessageResponse, error) {

	channel, postAt, err := s.client.PostMessage(input.ChannelID, slack.MsgOptionText(input.Message, false))
	if err != nil {
		return nil, err
	}

	timePostAt, err := time.Parse(time.RFC3339, postAt)
	if err != nil {
		return nil, err
	}

	return &ChatMessageResponse{
		ChannelID: channel,
		Message:   input.Message,
		PostAt:    timePostAt,
	}, nil
}
