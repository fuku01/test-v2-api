package slack

import (
	"encoding/json"
	"io"
	"net/http"
)

// Slack Event APIを初めて利用する際にURLVerification(検証)を行うための処理
// https://api.slack.com/events/url_verification

func SlackURLVerification(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Token     string `json:"token"`
		Challenge string `json:"challenge"`
		Type      string `json:"type"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "could not decode request body", http.StatusBadRequest)
		return
	}

	if reqBody.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(reqBody.Challenge))
		return
	}
}
