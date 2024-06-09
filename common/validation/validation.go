package validation

import (
	"errors"
	"go-dating-app/app/dto"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
)

var (
	//nolint:gochecknoglobals // singleton.
	validation *validator.Validate

	//nolint:gochecknoglobals // singleton.
	validationTrans ut.Translator
)

func NewValidation() error {
	validation = validator.New(validator.WithRequiredStructEnabled())

	// RegisterTagNameFunc is used to customize the name of the field with tag
	// taken from json instead struct property name.
	validation.RegisterTagNameFunc(func(fld reflect.StructField) string {
		splitNumber := 2
		name := strings.SplitN(fld.Tag.Get("json"), ",", splitNumber)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Set translation to english.
	english := en.New()
	uni := ut.New(english, english)
	validationTrans, _ = uni.GetTranslator("en")
	if err := enTrans.RegisterDefaultTranslations(validation, validationTrans); err != nil {
		return err
	}
	return nil
}

// ValidateVar validate a variable and returns error if any.
func ValidateVar(s any, tag string) error {
	return validation.Var(s, tag)
}

// ValidateStruct validate struct and returns error if any.
func ValidateStruct(s any) error {
	return validation.Struct(s)
}

// FormatStructErrors translate errors message from ValidateStruct() to user friendly messages.
func FormatStructErrors(errs error) []dto.ValidationErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		responses := make([]dto.ValidationErrorResponse, len(ve))
		for i, e := range ve {
			responses[i] = dto.ValidationErrorResponse{
				Field:   e.Field(),
				Message: e.Translate(validationTrans),
			}
		}
		return responses
	}
	return nil
}
