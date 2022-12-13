# Go × Echoで色々試すリポジトリ

## バリデーション

### 1. アプリケーション立ち上げ
```shell
go run main.go
```

### 2. リクエストを送ってバリデーションが機能するか確認する 
nameやemailを送らなかったり、無効なメール形式にしたり...etc
```shell
# HTTPie
http http://localhost:1323/users\?name\=test\&email\=test
http POST http://localhost:1323/users name=test email=test

# cURL
curl http://localhost:1323/users\?name\=test\&email\=test 
curl -X POST -d '{"name":"test", "email":"test"}' http://localhost:1323/users
```
