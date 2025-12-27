package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}()

func Validate[T any](req T) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errs := map[string]string{}
	for _, e := range err.(validator.ValidationErrors) {
		errs[e.Field()] = humanMessage(e)
	}
	return errs
}

func humanMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	default:
		return "invalid value"
	}
}
