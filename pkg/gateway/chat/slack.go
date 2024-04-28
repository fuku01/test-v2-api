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

func (s *Slack) PostMessage(req *PostMessageRequest) (*PostMessageResponse, error) {

	channel, postAt, err := s.client.PostMessage(req.ChannelID, slack.MsgOptionText(req.Message, false))
	if err != nil {
		return nil, err
	}

	timePostAt, err := time.Parse(time.RFC3339, postAt)
	if err != nil {
		return nil, err
	}

	return &PostMessageResponse{
		ChannelID: channel,
		Message:   req.Message,
		PostAt:    timePostAt,
	}, nil
}
