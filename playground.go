package main

import (
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
	"log"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type (
	// リクエストパラメータを埋め込む構造体
	User struct {
		// validate：バリデーションの内容、ja：フィールドの日本語名、query：GETのクエリストリングのkey、json：POSTリクエストボディのkey
		// is-messiはカスタムバリデーション
		// uuidはstring型でないと正しくチェックできなそうなので、カスタマイズが必要そう
		Id    uuid.UUID `query:"id" json:"id" validate:"required,uuid" ja:"ID"`
		Name  string    `query:"name" json:"name" validate:"required,is-messi" ja:"ユーザー名"`
		Email string    `query:"email" json:"email" validate:"required,email" ja:"メールアドレス"`
	}

	CustomValidator struct {
		trans     ut.Translator
		validator *validator.Validate
	}
)

func InitValidator() echo.Validator {
	jaTrans := ja.New()
	uniTrans := ut.New(jaTrans, jaTrans)
	trans, _ := uniTrans.GetTranslator("ja")

	validate := validator.New()

	// エラーメッセージの日本語化
	if err := jaTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

	// フィールド名の日本語化（jaタグを登録）
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	// カスタム型を登録
	validate.RegisterCustomTypeFunc(ValidateUuidValuer, uuid.UUID{})

	// カスタムバリデーションを登録
	if err := validate.RegisterValidation("is-messi", ValidateIsMessi); err != nil {
		log.Fatal(err)
	}

	// カスタムバリデーションの日本語エラーメッセージを登録
	if err := validate.RegisterTranslation("is-messi", trans, func(ut ut.Translator) error {
		return ut.Add("is-messi", "{0}はmessi以外認めません", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is-messi", fe.Field())

		return t
	}); err != nil {
		log.Fatal(err)
	}

	return &CustomValidator{
		trans:     trans,
		validator: validate,
	}
}

// カスタムバリデーション
func ValidateIsMessi(fl validator.FieldLevel) bool {
	return fl.Field().String() == "messi"
}

// uuid.UUID型を登録
func ValidateUuidValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.UUID); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

// バリデーション処理
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var messages []string
		for _, m := range err.(validator.ValidationErrors).Translate(cv.trans) {
			messages = append(messages, m)
		}
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, messages)
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = InitValidator()

	e.GET("/users", func(c echo.Context) (err error) {
		u := new(User)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	})

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
