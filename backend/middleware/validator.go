package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en_GB"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// Validate a struct, translating field errors into human readable messages
func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		messages := make([]string, len(fieldErrors))
		for index, fe := range fieldErrors {
			messages[index] = fe.Translate(v.translator)
		}
		return echo.NewHTTPError(http.StatusUnprocessableEntity, strings.Join(messages, "\n"))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

func NewValidator() *Validator {
	v := validator.New()
	locale := en_GB.New()
	translator, _ := ut.New(locale, locale).GetTranslator(locale.Locale())
	if err := enTranslations.RegisterDefaultTranslations(v, translator); err != nil {
		panic(err)
	}
	return &Validator{v, translator}
}
