package graphql

import "fmt"

// 複数のHandlerを1つにまとめる
type GraphQLHandler struct {
	MessageHandler MessageHandler
}

// エラー文を定義
var (
	InternalServerError = fmt.Errorf("Internal Server Error: 内部エラーが発生しました")
	InvalidRequest      = fmt.Errorf("Invalid Request: 無効なリクエストです")
)
