package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fuku01/test-v2-api/pkg/gateway/chat"
	graph "github.com/fuku01/test-v2-api/pkg/graph/generated"

	"github.com/rs/cors"

	"github.com/fuku01/test-v2-api/config"

	graph_handler "github.com/fuku01/test-v2-api/pkg/handler/graph"
	message_handler "github.com/fuku01/test-v2-api/pkg/handler/graph"
	webhook_handler "github.com/fuku01/test-v2-api/pkg/handler/webhook"
	message_repository "github.com/fuku01/test-v2-api/pkg/infrastructure/mysql"
	message_usecase "github.com/fuku01/test-v2-api/pkg/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("環境変数が不足しています。PORT: %s", port)
	}
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackToken == "" {
		log.Fatalf("環境変数が不足しています。SLACK_BOT_TOKEN: %s", slackToken)
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	slackClient, err := chat.NewSlack(slackToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("slackClient: ", slackClient) //!後で消す

	// 依存性の注入
	tr := message_repository.NewMessageRepository(db)
	tu := message_usecase.NewMessageUsecase(tr)
	th := message_handler.NewMessageHandler(tu)

	wh := webhook_handler.NewSlackHandler(tu)

	gh := graph_handler.GraphQLHandler{
		MessageHandler: th,
	}

	// SlackEventAPI(Webhook) エンドポイントの設定
	http.HandleFunc("/slack/events/verification", func(w http.ResponseWriter, r *http.Request) {
		wh.SlackURLVerification(w, r)
	})
	http.HandleFunc("/slack/events/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("========================SlackEventAPI(Webhook) が呼ばれました==============================")
		wh.CreateTodo(w, r)
	})

	// GraphQL ルーティングするハンドラーの設定
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			Handler: gh,
		}}),
	)
	// GraphQL ルーティングの設定
	http.Handle("/query", (srv))
	// GraphQL Playgroundの設定
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// CORSの設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // すべてのオリジンを許可
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	httpHandler := c.Handler(http.DefaultServeMux)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	// サーバーの起動
	err = http.ListenAndServe(":"+port, httpHandler)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
