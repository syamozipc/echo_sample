package main

import (
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"reflect"

	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	"github.com/labstack/echo/v4"
)

type (
	User struct {
		Name  string `json:"name" validate:"required" ja:"ユーザー名"`
		Email string `json:"email" validate:"required,email" ja:"メールアドレス"`
	}

	CustomValidator struct {
		trans     ut.Translator
		validator *validator.Validate
	}
)

func InitValidator() echo.Validator {
	ja := ja.New()
	uni := ut.New(ja, ja)
	trans, _ := uni.GetTranslator("ja")

	validate := validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})
	if err := ja_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatalf(err.Error())
	}

	return &CustomValidator{
		trans:     trans,
		validator: validate,
	}
}

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
		u := new(User)

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
