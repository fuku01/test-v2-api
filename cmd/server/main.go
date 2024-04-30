package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/fuku01/test-v2-api/pkg/graph/generated"

	"github.com/rs/cors"

	"github.com/fuku01/test-v2-api/db/config"

	graphql_handler "github.com/fuku01/test-v2-api/pkg/handler/graphql"
	webhook_handler "github.com/fuku01/test-v2-api/pkg/handler/webhook"
	repository "github.com/fuku01/test-v2-api/pkg/infrastructure/mysql"
	chat "github.com/fuku01/test-v2-api/pkg/infrastructure/slack"
	"github.com/fuku01/test-v2-api/pkg/usecase"
)

func httpHandlerFuncMiddleware(handlerFunc func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	httpHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
	return httpHandlerFunc
}

const requestTimeout = 5 * time.Minute

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("環境変数 PORT が設定されていません")
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	chatClient, err := chat.NewSlack()
	if err != nil {
		log.Fatal(err)
	}

	// 依存性の注入
	tr := repository.NewMessageRepository(db)
	tu := usecase.NewMessageUsecase(tr, chatClient)
	th := graphql_handler.NewMessageHandler(tu)

	wh := webhook_handler.NewSlackHandler(tu)

	gh := graphql_handler.GraphQLHandler{
		MessageHandler: th,
	}

	// SlackEventAPI(Webhook) エンドポイントの設定
	http.HandleFunc("/slack/events/verification", httpHandlerFuncMiddleware(wh.SlackURLVerification))
	http.HandleFunc("/slack/events/create", httpHandlerFuncMiddleware(wh.CreateTodo))

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

	// サーバーの設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      httpHandler,
		ReadTimeout:  requestTimeout,
		WriteTimeout: requestTimeout,
	}

	// サーバーの起動
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
