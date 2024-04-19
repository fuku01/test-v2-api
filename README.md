# test-v2-api

## migration
golang-migrate を使用
https://github.com/golang-migrate/migrate

### migration SQLファイル作成(local)
```
migrate create -ext sql -dir db/migrations -seq create_samples_table
```

### migrate up
dockerコンテナ内からコマンド実行するため、host.docker.internalを使用
mysqlの場合はtcp()を指定する必要がある
末尾の番号はmigrationファイルの番号
```
migrate -database "mysql://root:pass@tcp(host.docker.internal:3306)/v2_development?" -path db/migrations up 1
```
### migrate down
```
migrate -database "mysql://root:pass@tcp(host.docker.internal:3306)/v2_development?" -path db/migrations down 1
```
