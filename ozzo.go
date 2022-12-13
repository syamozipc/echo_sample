package main

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	OzzoUser struct {
		Id    uuid.UUID `query:"id" json:"id"`
		Name  string    `query:"name" json:"name"`
		Email string    `query:"email" json:"email"`
	}

	OzzoCustomValidator struct{}
)

func (cv *OzzoCustomValidator) Validate(i interface{}) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func (u OzzoUser) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Id, validation.Required.Error("IDは必須です"), is.UUID.Error("UUIDは正しい形式で指定してください")),
		// Streetは空を許容せず、5から50までの長さ
		validation.Field(&u.Name, validation.Required.Error("名前は必須です")),
		// Cityは空を許容せず、5から50までの長さ
		validation.Field(&u.Email, validation.Required.Error("メールアドレスは必須です"), is.Email.Error("正しい形式で入力してください")),
	); err != nil {
		return err
	}

	if err := validation.Validate(u.Name, validation.By(OzzoValidateIsMessi)); err != nil {
		return err
	}

	return nil
}

// カスタムバリデーション
func OzzoValidateIsMessi(value interface{}) error {
	s, _ := value.(string)
	if s != "messi" {
		return errors.New("メッシ以外は認めません")
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &OzzoCustomValidator{}

	e.GET("/users", func(c echo.Context) (err error) {
		u := new(OzzoUser)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	})

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(OzzoUser)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
