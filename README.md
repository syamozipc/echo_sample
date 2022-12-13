# Go × Echoで色々試すリポジトリ

## バリデーション

### 1. バリデーションサンプルを立ち上げる

```shell
# go-playground validator
go run playground.go

# Ozzo Validation
go run ozzo.go
```

### 2. リクエストを送ってバリデーションが機能するか確認する 
nameやemailを送らなかったり、無効なメール形式にしたり...etc
```shell
# httpie
http http http://localhost:1323/users\?name\=test
http http://localhost:1323/users\?name\=test\&email\=test

# curl
curl http://localhost:1323/users\?name\=test\&email\=test 
curl -X POST -d '{"name":"test", "email":"test"}' http://localhost:1323/users
```
