package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/fuku01/test-v2-api/pkg/graph/generated"
	chat "github.com/fuku01/test-v2-api/pkg/infrastructure/slack"

	"github.com/rs/cors"

	"github.com/fuku01/test-v2-api/db/config"

	graph_handler "github.com/fuku01/test-v2-api/pkg/handler/graph"
	message_handler "github.com/fuku01/test-v2-api/pkg/handler/graph"
	webhook_handler "github.com/fuku01/test-v2-api/pkg/handler/webhook"
	message_repository "github.com/fuku01/test-v2-api/pkg/infrastructure/mysql"
	message_usecase "github.com/fuku01/test-v2-api/pkg/usecase"
)

func httpHandlerFuncMiddleware(handlerFunc func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	httpHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
	return httpHandlerFunc
}

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
	tr := message_repository.NewMessageRepository(db)
	tu := message_usecase.NewMessageUsecase(tr, chatClient)
	th := message_handler.NewMessageHandler(tu)

	wh := webhook_handler.NewSlackHandler(tu)

	gh := graph_handler.GraphQLHandler{
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

	// サーバーの起動
	err = http.ListenAndServe(":"+port, httpHandler)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
