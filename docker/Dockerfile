# syntax=docker/dockerfile:1
# Dockerfileのバージョン指定（https://docs.docker.com/engine/reference/builder/#syntax）

# Dockerfileの書き方については、公式ドキュメントを参照（https://docs.docker.com/engine/reference/builder/）

# ベースイメージの指定（https://hub.docker.com/_/golangのイメージを使用）
FROM golang:1.21-alpine

# ログに出力する時間をJSTにするため、タイムゾーンを設定
ENV TZ /usr/share/zoneinfo/Asia/Tokyo
# ModuleモードをON
ENV GO111MODULE=on

# ワーキングディレクトリ（コンテナ内の作業ディレクトリ）を指定
WORKDIR /go/src/app
# ホストPCのDockerfileが存在するディレクトリ（ビルドコンテキスト）の内容を、コンテナの作業ディレクトリにコピー
COPY . .

# apkでシステムレベルのパッケージをインストール（gitとmake）
RUN apk add --no-cache --update git make

# go installでgoのライブラリをインストール（airとgqlgen、gomigrate）
RUN go install github.com/cosmtrek/air@v1.49.0
RUN go install github.com/99designs/gqlgen@latest
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# delve（dockeコンテナ内でデバッグを行うためのツール）をインストール
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# go.modを参照し、go.sumファイルの更新を行う
RUN go mod tidy

# ポート番号の指定（番号は任意で、他のコンテナと被らないようにする）
EXPOSE 4000

# 起動時のコマンド（ホットリロードするため、airを使用）
CMD ["air", "-c", ".air.toml"]