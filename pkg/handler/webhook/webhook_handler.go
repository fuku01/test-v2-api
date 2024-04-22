package webhook

import (
	"encoding/json"
	"io"
	"net/http"

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

	eventsAPIEvent, err := h.parseWebhookSlackEvents(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	var input *model.CreateTodoInput
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			input = &model.CreateTodoInput{
				Content: ev.Text,
			}
		}
	}

	if input == nil {
		return nil
	}

	todo, err := h.tu.CreateTodo(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
	return nil
}

func (h *webhookHandler) parseWebhookSlackEvents(w http.ResponseWriter, r *http.Request) (*slackevents.EventsAPIEvent, error) {
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
