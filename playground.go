package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	User struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	CustomValidatorEn struct {
		validatorEn *validator.Validate
	}
)

func (cv *CustomValidatorEn) Validate(i interface{}) error {
	if err := cv.validatorEn.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidatorEn{validatorEn: validator.New()}
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
