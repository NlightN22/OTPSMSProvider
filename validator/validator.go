package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var e164Regexp = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func RegisterCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("e164", func(fl validator.FieldLevel) bool {
			return e164Regexp.MatchString(fl.Field().String())
		})
	}
}
