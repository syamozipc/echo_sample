package main

import (
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
	"log"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type (
	// リクエストパラメータを埋め込む構造体
	CreateUser struct {
		// validateにはバリデーションの内容を、jaにはフィールドの日本語名を入れる
		// is-messiはカスタムバリデーション
		Name  string `json:"name" validate:"required,is-messi,min=1,max=1" ja:"ユーザー名"`
		Email string `json:"email" validate:"required,email" ja:"メールアドレス"`
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

	// フィールド名の日本語化（jaタグを登録）
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	// エラーメッセージの日本語化
	if err := jaTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

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

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(CreateUser)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
