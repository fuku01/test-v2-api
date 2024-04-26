package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/slack-go/slack/slackevents"
)

type WebhookHandler interface {
	CreateTodo(w http.ResponseWriter, r *http.Request) error
}

type webhookHandler struct {
	tu usecase.TodoUsecase
}

func NewWebhookHandler(tu usecase.TodoUsecase) WebhookHandler {
	return &webhookHandler{
		tu: tu,
	}
}

func (h *webhookHandler) CreateTodo(w http.ResponseWriter, r *http.Request) error {

	eventsAPIEvent, err := h.parseSlackEventsAPIRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	event, err := h.checkSlackEventsAPIEvent(eventsAPIEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	// メンションされたユーザー名を削除
	re := regexp.MustCompile(`<@[^>]+>`)
	message := re.ReplaceAllString(event.Text, "")

	input := &model.CreateTodoInput{
		Content: message,
	}
	todo, err := h.tu.CreateTodo(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
	return nil
}

// SlackEventAPIからのHTTPリクエスト(Body)から、イベントデータをパースする処理
func (h *webhookHandler) parseSlackEventsAPIRequest(w http.ResponseWriter, r *http.Request) (*slackevents.EventsAPIEvent, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return nil, err
	}

	return &eventsAPIEvent, nil
}

// SlackEventAPIのEventがコールバックイベントかつ、メンションイベントであるかをチェックする処理
func (h *webhookHandler) checkSlackEventsAPIEvent(events *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	// イベントがコールバックイベントであるかをチェック
	if events.Type != slackevents.CallbackEvent {
		return nil, fmt.Errorf("event type is not callback event")
	}

	// イベントがメンションイベントであるかをチェック
	event, ok := events.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		return nil, fmt.Errorf("inner event is not message event")
	}

	return event, nil

}
