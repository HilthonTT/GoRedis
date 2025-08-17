package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("form")
		name := strings.SplitN(tag, ",", 2)[0]

		if name == "-" || name == "" {
			return ""
		}

		return name
	})

	return v
}
