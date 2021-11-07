package main

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

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
			validation.RuneLength(5, 20).Error("名前は {min}～{max} 文字です"),
			is.PrintableASCII.Error("名前はASCIIで入力して下さい"),
		),
		validation.Field(
			&a.Email,
			validation.Required.Error("メールアドレスは必須入力です"),
			validation.RuneLength(5, 40).Error("メールアドレスは {min}～{max} 文字です"),
			is.Email.Error("メールアドレスを入力して下さい"),
		),
		validation.Field(
			&a.Content,
			validation.Required.Error("本文は必須入力です"),
			validation.RuneLength(5, 50).Error("本文は {min}～{max} 文字です"),
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
			re := regexp.MustCompile(`{[a-z]+}`)
			for k, err := range errs {
				c.Logger().Error(k + ": " + err.Error())
				verr := err.(validation.Error)
				params := verr.Params()
				errs[k] = errors.New(re.ReplaceAllStringFunc(verr.Message(), func(s string) string {
					return fmt.Sprint(params[strings.Trim(s, "{}")])
				}))
			}
			return errs
		}
		return c.JSON(http.StatusOK, &struct {
			Result string `json:"result"`
		}{
			Result: "OK",
		})
	})
	e.Logger.Fatal(e.Start(":8989"))
}
