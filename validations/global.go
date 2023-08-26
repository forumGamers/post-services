package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func GetValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("text_validate",postTextRegex)

	return validate
}

func postTextRegex(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`/[^a-zA-Z0-9.,\-\s\@()#*&]%$!?/g`).MatchString(fl.Field().String())
}