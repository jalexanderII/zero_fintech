package middleware

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

// ValidateStruct is a generic validator for all serializers
func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate = validator.New()
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
