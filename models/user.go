package models

// リクエストパラメータを埋め込む構造体
type User struct {
	// query：GETのクエリストリングのkey、json：POSTリクエストボディのkey、pamam：パスパラメータ名、validate：バリデーションの内容、ja：フィールドの日本語名
	// is-messiはカスタムバリデーション
	// uuidは不正な値がきたときにecho.Bindが失敗してエラーを返すので、バリデーションまでこない。string型にしてvalidateにuuidを追加し、詰め替え時にuuid型にする
	// Null〇〇もBind失敗するので、ポインタ型でバリデーションし、domainなどへの詰め替え時にNull〇〇にする
	Id       string  `query:"id" param:"id" json:"id" validate:"required,uuid" ja:"ID"`
	Name     string  `query:"name" json:"name" validate:"required,is-messi" ja:"ユーザー名"`
	Age      *int    `query:"age" json:"age" validate:"omitempty,number" ja:"年齢"`
	Email    string  `query:"email" json:"email" validate:"required,email" ja:"メールアドレス"`
	Gender   *string `query:"gender" json:"gender" validate:"omitempty,oneof=男性 女性 その他" ja:"性別"`
	IsActive *bool   `query:"active" json:"active" validate:"omitempty,boolean" ja:"アクティブ"`
}

type User2 struct {
	Id       string  `query:"id" param:"id" json:"id" validate:"required,uuid" ja:"ID"`
	Name     string  `query:"name" json:"name" validate:"required,is-messi" ja:"ユーザー名"`
	Age      *int    `query:"age" json:"age" validate:"omitempty,number" ja:"年齢"`
	Email    string  `query:"email" json:"email" validate:"required,email" ja:"メールアドレス"`
	Gender   *string `query:"gender" json:"gender" validate:"omitempty,oneof=男性 女性 その他" ja:"性別"`
	IsActive *bool   `query:"active" json:"active" validate:"omitempty,boolean" ja:"アクティブ"`
	Nums     []int   `json:"nums"`
	Friend   *User
	Friends  []*User
	Aiu      Aiu
}

type Aiu struct {
	Aa string `json:"aa" validate:"required"`
}

type Sample struct {
	ID   string `json:"id" validate:"omitempty" ja:"ID1"`
	Name string `json:"name" validate:"omitempty" ja:"name1"`
	// point型なら、nullかkeyなしでゼロ値（nil）となるため、required/mitemptyが機能する
	// omitemptyの場合はその構造体のフィールドのバリデーションまで行かない
	// {}ではゼロ値でなくなってしまうので注意
	Sample2 *Sample2 `json:"sample2" validate:"omitempty"`
	// 値型の場合、nullでもkeyなしでもrequiredを必ず通過し、omitemptyには絶対当てはまらず、意図した挙動にならない
	Sample3 Sample3 `json:"sample3" validate:"omitempty"`
}
type Sample2 struct {
	ID   string `json:"id2" validate:"required" ja:"ID2"`
	Name string `json:"name2" validate:"omitempty" ja:"name2"`
}
type Sample3 struct {
	ID   *string `json:"id3" validate:"required" ja:"ID3"`
	Name *string `json:"name3" validate:"omitempty" ja:"name3"`
}
