package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fuku01/test-v2-api/pkg/gateway/slack"
	graph "github.com/fuku01/test-v2-api/pkg/graph/generated"

	"github.com/rs/cors"

	"github.com/fuku01/test-v2-api/config"

	h "github.com/fuku01/test-v2-api/pkg/handler"
	todo_handler "github.com/fuku01/test-v2-api/pkg/handler"
	todo_repository "github.com/fuku01/test-v2-api/pkg/infrastructure"
	todo_usecase "github.com/fuku01/test-v2-api/pkg/usecase"
)

func main() {
	port := 4000

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	slackClient, err := config.NewSlack()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("slackClient: ", slackClient) //!後で消す

	// 依存性の注入
	tr := todo_repository.NewTodoRepository(db)
	tu := todo_usecase.NewTodoUsecase(tr)
	th := todo_handler.NewTodoHandler(tu)

	h := h.Handler{
		TodoHandler: th,
	}

	// Slack Webhookのエンドポイントの設定
	http.HandleFunc("/slack/events/verification", func(w http.ResponseWriter, r *http.Request) {
		slack.SlackURLVerification(w, r)
	})
	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		log.Println("slack/eventsが呼ばれました")
		th.ListTodos()
	})

	// GraphQL ルーティングするハンドラーの設定
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			Handler: h,
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

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)

	// サーバーの起動
	err = http.ListenAndServe(":"+strconv.Itoa(port), httpHandler)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
