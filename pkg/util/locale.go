package util

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

func TranslateValidationErrors(vE validator.ValidationErrors) map[string]string {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale, en.New())

	translator, ok := uni.GetTranslator("en")
	if ok {
		return vE.Translate(translator)
	}

	return map[string]string{}
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func TranslateFiberValidationErrors(input any) map[string]string {

	var validate = validator.New()

	var errors []*ErrorResponse
	err := validate.Struct(input)
	plainError := make(map[string]string)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = strings.ToLower(err.Field())
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
			if err.Param() != "" {
				plainError[element.FailedField] = strings.ToLower(err.Field()) + " is required and must be " + err.Param()
			} else if err.Tag() == "e164" {
				plainError[element.FailedField] = "phone is required and must be in e164 format"
			} else if err.Tag() == "lowercase" {
				plainError[element.FailedField] = strings.ToLower(err.Field()) + " is required and must be " + err.Tag()
			} else {
				plainError[element.FailedField] = strings.ToLower(err.Field()) + " is required"
			}
		}
	}

	//return errors
	if err != nil {
		return plainError
	}
	return nil
}
