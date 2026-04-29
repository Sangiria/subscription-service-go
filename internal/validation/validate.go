package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(st any, tag *string) error {
    v := validator.New()
    if tag != nil {
        v.SetTagName(*tag)
    }

    err := v.Struct(st)
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