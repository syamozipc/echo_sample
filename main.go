package main

import (
	"echo_sample/models"
	"echo_sample/validation"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = validation.InitValidator()

	e.GET("/sample", func(c echo.Context) (err error) {
		u := new(models.Sample)

		if err = c.Bind(u); err != nil {
			fmt.Println("bindエラー")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Printf("%#v", u)

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/users", func(c echo.Context) (err error) {
		u := new(models.User)

		if err = c.Bind(u); err != nil {
			fmt.Println("bindエラー")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Printf("%#v", u)

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	})

	e.GET("/users/:id", func(c echo.Context) (err error) {
		u := new(models.User)

		if err = c.Bind(u); err != nil {
			fmt.Println("bindエラー")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Printf("%#v", u)

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	})

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(models.User2)

		if err = c.Bind(u); err != nil {
			fmt.Println("bindエラー")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		fmt.Printf("%#v", u)

		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
