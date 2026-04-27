package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func BindAndValidate(c *echo.Context, req interface{}, tag string) error {
    if err := c.Bind(req); err != nil {
        return fmt.Errorf("data parsing error")
    }

    v := validator.New()
    v.SetTagName(tag)

    err := v.Struct(req)
    if err == nil {
        return nil
    }

	if vErrors, ok := err.(validator.ValidationErrors); ok {
        var errMsgs []string
        for _, ve := range vErrors {
            errMsgs = append(errMsgs, fmt.Sprintf("field %s: %s", ve.Field(), ve.Tag()))
        }
        return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
    }

    return err
}