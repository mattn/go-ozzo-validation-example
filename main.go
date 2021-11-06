package main

import (
	"embed"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Comment struct {
	Name    string `json:"name"`
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
		validation.Field(&a.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.Content, validation.Required, validation.Length(5, 50)),
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
