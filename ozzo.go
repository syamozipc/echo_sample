package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User2 struct {
	Name  string
	Email string
}

func (u User2) Validate() error {
	return validation.ValidateStruct(&u,
		// Streetは空を許容せず、5から50までの長さ
		validation.Field(&u.Name, validation.Required.Error("名前は必須です"), validation.Length(5, 50).Error("名前は{min}以上{max}以下です")),
		// Cityは空を許容せず、5から50までの長さ
		validation.Field(&u.Email, validation.Required, validation.Length(5, 50)),
	)
}

func main() {
	e := echo.New()

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User2)

		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := u.Validate(); err != nil {
			if e, ok := err.(validation.InternalError); ok {
				return c.JSON(http.StatusInternalServerError, e.InternalError())
			}
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
