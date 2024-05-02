package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/rs/cors"

	"github.com/fuku01/test-v2-api/db/config"

	graph "github.com/fuku01/test-v2-api/pkg/graph/generated"
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
	// 1. ポート番号を設定
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("環境変数 PORT が設定されていません")
	}

	// 2. DB接続
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Slackクライアントの作成
	chatClient, err := chat.NewSlack()
	if err != nil {
		log.Fatal(err)
	}

	// @ 4. 依存性の注入
	// repository
	tr := repository.NewMessageRepository(db)
	tu := usecase.NewMessageUsecase(tr, chatClient)
	mh := graphql_handler.NewMessageHandler(tu)
	// webhook
	wh := webhook_handler.NewSlackHandler(tu)

	graphQLHandler := graphql_handler.GraphQLHandler{
		MessageHandler: mh,
	}

	// 5. GraphQL サーバーの作成
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			Handler: graphQLHandler,
		}}),
	)

	// 6. 各エンドポイントの設定
	// GraphQL
	http.Handle("/query", (srv))
	// GraphQL Playground
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// Webhook（Slack Event API）
	http.HandleFunc("/slack/events/verification", httpHandlerFuncMiddleware(wh.SlackURLVerification))
	http.HandleFunc("/slack/events/create", httpHandlerFuncMiddleware(wh.CreateTodo))

	// 7. CORSの設定（異なるオリジン(異なるドメインやプロトコル、ポート番号)からのリクエストを許可する）
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                            // すべてのオリジンを許可
		AllowCredentials: true,                                     // クレデンシャル情報（Cookieなど）を許可
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // 許可するHTTPメソッド
	})
	httpHandler := c.Handler(http.DefaultServeMux)

	// 8. サーバーの設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      httpHandler,
		ReadTimeout:  requestTimeout,
		WriteTimeout: requestTimeout,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	// 9. サーバーの起動
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
