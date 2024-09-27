package web

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterValidation("imagefile", func(fl validator.FieldLevel) bool {

		file, ok := fl.Field().Interface().(FileUpload)
		if !ok {
			return true
		}

		if file.Data == nil {
			return true
		}

		if !strings.HasPrefix(file.Mime, "image/") {
			return false
		}

		return true
	}, true)

	validate.RegisterValidation("filerequired", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(FileUpload)
		if !ok {
			return false
		}

		if file.Data == nil {
			return false
		}

		return true
	}, true)

	validate.RegisterValidation("maxsize", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(FileUpload)
		if !ok {
			return true
		}

		if file.Data == nil {
			return true
		}

		size := fl.Param()

		multiplier := 1

		size, hasM := strings.CutSuffix(size, "MB")
		if hasM {
			multiplier = 1000 * 1000
		}

		size, hasK := strings.CutSuffix(size, "KB")

		if hasK {
			multiplier = 1000
		}

		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return true
		}

		if sizeInt*multiplier < len(file.Data) {
			return false
		}

		return true

	})

	validate.RegisterTranslation("maxsize", translator, func(ut ut.Translator) error {
		return ut.Add("maxsize", "{0} must be less than {1}", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("maxsize", fe.Field(), fe.Param())
		return t
	})

	validate.RegisterTranslation("imagefile", translator, func(ut ut.Translator) error {
		return ut.Add("imagefile", "{0} must be a valid image file", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("imagefile", fe.Field())
		return t
	})

	validate.RegisterTranslation("filerequired", translator, func(ut ut.Translator) error {
		return ut.Add("filerequired", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filerequired", fe.Field())
		return t
	})

	validate.RegisterValidation("katakana", ValidateKataKana)
	validate.RegisterTranslation("katakana", translator, func(ut ut.Translator) error {
		return ut.Add("katakana", "{0} must be a valid katakana", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("katakana", fe.Field())
		return t
	})

	validate.RegisterValidation("countrycode", ValidateCountryCode)
	validate.RegisterTranslation("countrycode", translator, func(ut ut.Translator) error {
		return ut.Add("countrycode", "{0} must be a valid country code", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("countrycode", fe.Field())
		return t
	})

	validate.RegisterValidation("countrycodes", ValidateCountryCodeSlice)
	validate.RegisterTranslation("countrycodes", translator, func(ut ut.Translator) error {
		return ut.Add("countrycodes", "{0} must be a valid list of country codes", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("countrycodes", fe.Field())
		return t
	})

	validate.RegisterValidation("languagecode", ValidateLanguageCode)
	validate.RegisterTranslation("languagecode", translator, func(ut ut.Translator) error {
		return ut.Add("languagecode", "{0} must be a valid language code", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("languagecode", fe.Field())
		return t
	})

	validate.RegisterValidation("languagecodes", ValidateLanguageCodeSlice)
	validate.RegisterTranslation("languagecodes", translator, func(ut ut.Translator) error {
		return ut.Add("languagecodes", "{0} must be a valid list of language codes", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("languagecodes", fe.Field())
		return t
	})

	validate.RegisterValidation("japanesephonenumber", ValidateJapanesePhoneNumber)
	validate.RegisterTranslation("japanesephonenumber", translator, func(ut ut.Translator) error {
		return ut.Add("japanesephonenumber", "{0} must be a valid Japanese phone number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("japanesephonenumber", fe.Field())
		return t
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: validate, translator: translator}
}

func (cv *CustomValidator) Validate(val interface{}) error {
	if err := cv.validator.Struct(val); err != nil {

		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(cv.translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
