package main

import (
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	ja_translations "gopkg.in/go-playground/validator.v9/translations/ja"
)

type (
	User struct {
		Name  string `json:"name" validate:"required" ja:"ユーザー名"`
		Email string `json:"email" validate:"required,email" ja:"メールアドレス"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

var (
	uni         *ut.UniversalTranslator
	validateObj *validator.Validate
	trans       ut.Translator
)

func Init() {
	ja := ja.New()
	uni = ut.New(ja, ja)
	t, _ := uni.GetTranslator("ja")
	trans = t
	validateObj = validator.New()
	validateObj.RegisterTagNameFunc(func(fld reflect.StructField) string {
		fieldName := fld.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})
	ja_translations.RegisterDefaultTranslations(validateObj, trans)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, GetErrorMessages(err))
	}
	return nil
}

func main() {
	Init()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validateObj}
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

func GetErrorMessages(err error) []string {
	if err == nil {
		return []string{}
	}
	var messages []string
	for _, m := range err.(validator.ValidationErrors).Translate(trans) {
		messages = append(messages, m)
	}
	return messages
}
