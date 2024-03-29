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
	Num1 *int `json:"num1" validate:"omitempty,ltfield_if_Max_is_explicit=Num2" ja:"num1"`
	Num2 *int `json:"num2" validate:"omitempty" ja:"num2"`
}
type Sample2 struct {
	Name *string `json:"name" validate:"required_if=Flg false,exculded_if_for_bool=Flg true" ja:"名前"`
	Flg  *bool   `json:"flg" validate:"omitempty" ja:"フラグ"`
}
type Sample3 struct {
	ID   *string `json:"id3" validate:"required" ja:"ID3"`
	Name *string `json:"name3" validate:"omitempty" ja:"name3"`
}

/*
レスポンスの考え方
- オプショナル項目だけomitemptyをつける
- オプショナルかつゼロ値を返す場合はポインタ型にする（omitemptyで弾かれない様にするため）
- オプショナルな構造体は全てポインタ型（値型だとomitemptyが効かない）
- スライスは空配列で返す（models.NestSlice{}のように要素数0のスライスで初期化）
*/
type SampleRes struct {
	// 構造体（コメントはreturn値）
	Nest        Nest  `json:"nest"`                  // Nest構造体
	NestPtr     *Nest `json:"nestPtr"`               // 未代入（nil）：ゼロ値であるnull、他：Nest構造体
	NestOmit    Nest  `json:"nestOmit,omitempty"`    // Nest構造体（値型struct自体はゼロ値にならないのでプロパティ削除は起こらない）
	NestOmitPtr *Nest `json:"nestOmitPtr,omitempty"` // 未代入（nil）：プロパティ自体を削除、他：Nest構造体
	// スライス（コメントはreturn値）
	NestSclie        NestSlice  `json:"nestSlice"`                  // 未代入（nil）：nullが返る、他：Nestスライス
	NestScliePtr     *NestSlice `json:"nestSlicePtr"`               // 未代入（nil）：nullが返る、他：Nestスライス
	NestSclieOmit    NestSlice  `json:"nestSliceOmit,omitempty"`    // 未代入（nil）：プロパティ自体を削除、他：Nest構造体
	NestSclieOmitPtr *NestSlice `json:"nestSliceOmitPtr,omitempty"` // 未代入（nil）：プロパティ自体を削除、他：Nest構造体
}
type Nest struct {
	Content    string  `json:"content"`
	PtrContent *string `json:"ptrContent"`
}

type NestSlice []Nest
