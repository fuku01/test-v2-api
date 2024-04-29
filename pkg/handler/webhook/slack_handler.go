package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/fuku01/test-v2-api/context/logger"
	"github.com/fuku01/test-v2-api/pkg/domain/entity"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/slack-go/slack/slackevents"
)

type SlackHandler interface {
	SlackURLVerification(w http.ResponseWriter, r *http.Request)
	CreateTodo(w http.ResponseWriter, r *http.Request)
}

type slackHandler struct {
	tu usecase.MessageUsecase
}

func NewSlackHandler(tu usecase.MessageUsecase) SlackHandler {
	return &slackHandler{
		tu: tu,
	}
}

func (h *slackHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------スラックイベントAPIのCreateTodo()が呼ばれました---------------------")

	// HTTPリクエストからコンテキストを取得
	ctx := r.Context()

	// SlackEventAPIのリクエストをチェック
	body, err := h.checkSlackEventsAPIRequest(r)
	if err != nil {
		logger.Error("checkSlackEventsAPIRequest", err)
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

		req := &entity.CreateMessageRequest{
			Content: message,
		}
		_, err = h.tu.CreateMessage(ctx, req)
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
			logger.Error("goroutine in CreateTodo", err)
		}
	}()
}

// SlackEventAPIのリクエストをチェックして、再送リクエストの場合は無視、それ以外の場合はBodyを返す処理
func (h *slackHandler) checkSlackEventsAPIRequest(r *http.Request) ([]byte, error) {
	// リクエストをチェックして、再送リクエスト（ヘッダーにX-Slack-Retry-Numが存在する）の場合は無視する。※ SlackEventAPIの仕様により、再送リクエストは3回まで行われる（3秒,1分,5分後）（https://dev.classmethod.jp/articles/slack-resend-matome/）
	retryNum := r.Header.Get("X-Slack-Retry-Num")
	if retryNum != "" {
		return nil, nil
	}

	// HTTPリクエストからBodyを読み込む
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SlackEventAPIのHTTPリクエストBodyから、イベントデータをパースする処理
func (h *slackHandler) parseSlackEventsAPIRequestBody(body []byte) (*slackevents.EventsAPIEvent, error) {
	// BodyをパースしてSlackEventsAPIのイベントデータを取得
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return nil, err
	}

	return &eventsAPIEvent, nil
}

// SlackEventAPIのイベントが「コールバックイベントかつメンションイベント」であるかをチェックし、イベントデータを取得する処理
func (h *slackHandler) checkSlackEventsAPIEvent(events *slackevents.EventsAPIEvent) (*slackevents.AppMentionEvent, error) {
	// イベントがコールバックイベントであるかをチェック
	if events.Type != slackevents.CallbackEvent {
		return nil, errors.New("event type is not callback event")
	}

	// イベントがメンションイベントであるかをチェック
	event, ok := events.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		return nil, errors.New("inner event is not message event")
	}

	return event, nil

}

// Slackメンションイベントのテキストから不要な文字列を削除する処理
func (h *slackHandler) replaceSlackMentionEventText(event *slackevents.AppMentionEvent) string {
	re := regexp.MustCompile(`<@[^>]+>`)
	message := re.ReplaceAllString(event.Text, "")
	return message
}

/*
!Slack Event APIを初めて利用する際にURLVerification(検証)を行うための処理
?https://api.slack.com/events/url_verification

1.初めて登録するエンドポイントに対して、Slackからのリクエストが送れる
2.リクエストに含まれるChallengeをそのままレスポンスで返すことで、そのエンドポイントの検証が完了し登録可能となる
*/
func (h *slackHandler) SlackURLVerification(w http.ResponseWriter, r *http.Request) {
	// Slackからのリクエストを格納する構造体
	var reqBody struct {
		Token     string `json:"token"`
		Challenge string `json:"challenge"`
		Type      string `json:"type"`
	}

	// HTTPリクエストからボディを読み込む
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	// リクエストボディを構造体に変換
	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "could not decode request body", http.StatusBadRequest)
		return
	}

	// リクエストのタイプがURL検証（url_verification）の場合は、チャレンジ（Challenge）を返す
	if reqBody.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(reqBody.Challenge))
		return
	}
}
