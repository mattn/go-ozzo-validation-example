package main

import (
	"embed"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Comment struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Content string `json:"content"`
}

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func (a Comment) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Name,
			validation.Required.Error("名前は必須入力です"),
			validation.RuneLength(5, 20).Error("名前は 5～20 文字です"),
			is.PrintableASCII.Error("名前はASCIIで入力して下さい"),
		),
		validation.Field(
			&a.Email,
			validation.Required.Error("メールアドレスは必須入力です"),
			validation.RuneLength(5, 40).Error("メールアドレスは 5～40 文字です"),
			is.Email.Error("メールアドレスを入力して下さい"),
		),
		validation.Field(
			&a.Content,
			validation.Required.Error("本文は必須入力です"),
			validation.RuneLength(5, 50).Error("本文は 5～50 文字です"),
		),
	)
}

//go:embed static
var localFS embed.FS

func main() {
	e := echo.New()
	e.Debug = true
	e.Validator = &CustomValidator{}
	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Filesystem: http.FS(localFS),
	}))

	e.POST("/api", func(c echo.Context) error {
		var comment Comment
		if err := c.Bind(&comment); err != nil {
			return err
		}
		if err := c.Validate(comment); err != nil {
			errs := err.(validation.Errors)
			for k, err := range errs {
				c.Logger().Error(k + ": " + err.Error())
			}
			return err
		}
		return c.JSON(http.StatusOK, &struct {
			Result string `json:"result"`
		}{
			Result: "OK",
		})
	})
	e.Logger.Fatal(e.Start(":8989"))
}
