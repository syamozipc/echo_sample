package models

// リクエストパラメータを埋め込む構造体
type User struct {
	// query：GETのクエリストリングのkey、json：POSTリクエストボディのkey、pamam：パスパラメータ名、validate：バリデーションの内容、ja：フィールドの日本語名
	// is-messiはカスタムバリデーション
	// uuidは不正な値がきたときにecho.Bindが失敗してエラーを返すので、バリデーションまでこない。string型にしてvalidateにuuidを追加し、詰め替え時にuuid型にする
	// Null〇〇もBind失敗するので、ポインタ型でバリデーションし、domainなどへの詰め替え時にNull〇〇にする
	Id     string  `query:"id" param:"id" json:"id" validate:"required,uuid" ja:"ID"`
	Name   string  `query:"name" json:"name" validate:"required,is-messi" ja:"ユーザー名"`
	Age    *int    `query:"age" json:"age" validate:"omitempty,number" ja:"年齢"`
	Email  string  `query:"email" json:"email" validate:"required,email" ja:"メールアドレス"`
	Gender *string `query:"gender" json:"gender" validate:"omitempty,oneof=男性 女性 その他" ja:"性別"`
}
