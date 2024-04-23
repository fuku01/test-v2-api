package slack

import (
	"encoding/json"
	"io"
	"net/http"
)

/*
Slack Event APIを初めて利用する際にURLVerification(検証)を行うための処理
https://api.slack.com/events/url_verification

1.初めて登録するエンドポイントに対して、Slackからのリクエストが送れる
2.リクエストに含まれるChallengeをそのままレスポンスで返すことで、そのエンドポイントの検証が完了し登録可能となる
*/

func SlackURLVerification(w http.ResponseWriter, r *http.Request) {
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
