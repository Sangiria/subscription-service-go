package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func BindAndValidate(c *echo.Context, req interface{}, tag string) error {
    if err := c.Bind(req); err != nil {
        return err
    }

    v := validator.New()
    v.SetTagName(tag)
    
    return v.Struct(req)
}