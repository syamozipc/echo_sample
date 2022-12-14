package models

import uuid "github.com/satori/go.uuid"

// リクエストパラメータを埋め込む構造体
type User struct {
	// validate：バリデーションの内容、ja：フィールドの日本語名、query：GETのクエリストリングのkey、json：POSTリクエストボディのkey
	// is-messiはカスタムバリデーション
	// uuidはstring型でないと正しくチェックできなそうなので、カスタマイズが必要そう
	Id    uuid.UUID `query:"id" json:"id" validate:"required" ja:"ID"`
	Name  string    `query:"name" json:"name" validate:"required,is-messi" ja:"ユーザー名"`
	Age   int       `query:"age" json:"age" validate:"required,number" ja:"年齢"`
	Email string    `query:"email" json:"email" validate:"required,email" ja:"メールアドレス"`
}
