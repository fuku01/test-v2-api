package webhook

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/slack-go/slack/slackevents"
)

type WebhookHandler interface {
	CreateTodo(w http.ResponseWriter, r *http.Request)
}

type webhookHandler struct {
	tu usecase.MessageUsecase
}

func NewWebhookHandler(tu usecase.MessageUsecase) WebhookHandler {
	return &webhookHandler{
		tu: tu,
	}
}

func (h *webhookHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// SlackEventAPIのリクエストをチェック
	body, err := h.checkSlackEventsAPIRequest(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // 再送リクエストされないように200 OKを即座に返す

	// エラーを受け取るためのチャネルを作成
	errChan := make(chan error)
	// 非同期処理でSlackのイベントを処理
	go func() {
		defer close(errChan)

		eventsAPIEvent, err := h.parseSlackEventsAPIRequestBody(body)
		if err != nil {
			errChan <- err
			return
		}
		if eventsAPIEvent == nil {
			errChan <- nil
			return
		}
		event, err := h.checkSlackEventsAPIEvent(eventsAPIEvent)
		if err != nil {
			errChan <- err
			return
		}
		if event == nil {
			errChan <- nil
			return
		}
		message := h.replaceSlackMentionEventText(event)

		input := &domain_model.CreateMessageInput{
			Content: message,
		}
		_, err = h.tu.CreateMessage(input)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	// 非同期処理で発生したエラーをログに記録
	go func() {
		err := <-errChan
		if err != nil {
			slog.Error("failed to goroutine: in CreateTodo: ", err)
		}
	}()
}

// SlackEventAPIのリクエストをチェックして、再送リクエストの場合は無視、それ以外の場合はBodyを返す処理
func (h *webhookHandler) checkSlackEventsAPIRequest(r *http.Request) ([]byte, error) {
	// リクエストをチェックして、再送リクエスト（ヘッダーにX-Slack-Retry-Numが存在する）の場合は無視する。※ SlackEventAPIの仕様により、再送リクエストは3回まで行われる（3秒,1分,5分後）（https://dev.classmethod.jp/articles/slack-resend-matome/）
	retryNum := r.Header.Get("X-Slack-Retry-Num")
	if retryNum != "" {
		return nil, nil
	}

	// HTTPリクエストからBodyを読み込む
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("io.ReadAll: failed to read http request body: ", err)
		return nil, err
	}

	return body, nil
}

// SlackEventAPIのHTTPリクエストBodyから、イベントデータをパースする処理
func (h *webhookHandler) parseSlackEventsAPIRequestBody(body []byte) (*slackevents.EventsAPIEvent, error) {
	// BodyをパースしてSlackEventsAPIのイベントデータを取得
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		slog.Error("slackevents.ParseEvent: failed to parse request body: ", err)
		return nil, err
	}

	return &eventsAPIEvent, nil
}

// SlackEventAPIのイベントが「コールバックイベントかつメンションイベント」であるかをチェックし、イベントデータを取得する処理
func (h *webhookHandler) checkSlackEventsAPIEvent(events *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	// イベントがコールバックイベントであるかをチェック
	if events.Type != slackevents.CallbackEvent {
		slog.Error("checkSlackEventsAPIEvent: event type is not callback event")
		return nil, errors.New("event type is not callback event")
	}

	// イベントがメンションイベントであるかをチェック
	event, ok := events.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		slog.Error("checkSlackEventsAPIEvent: inner event is not message event")
		return nil, errors.New("inner event is not message event")
	}

	return event, nil

}

// Slackメンションイベントのテキストから不要な文字列を削除する処理
func (h *webhookHandler) replaceSlackMentionEventText(event *slackevents.AppMentionEvent) string {
	re := regexp.MustCompile(`<@[^>]+>`)
	message := re.ReplaceAllString(event.Text, "")
	return message
}
