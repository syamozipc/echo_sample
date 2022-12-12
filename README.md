# Go × Echoで色々試すリポジトリ

## バリデーション

### 1. 以下のいずれかのバリデーションサンプルを立ち上げる

go-playground validator（日本語対応）
```shell
go run ja_playground.go
```

go-playground validator
```shell
go run playground.go
```

Ozzo Validation
```shell
go run ozzo.go
```

### 2. リクエストを送ってバリデーションが機能するか確認する 
nameやemailを送らなかったり、無効なメール形式にしたり...etc
```shell
http post http://localhost:1323/users name=test email=test@test.com
```
