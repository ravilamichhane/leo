package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/ravilmc/leo/kana"

	"github.com/go-playground/validator/v10"
)

var AvailableCountryCodes = []string{
	"NP",
	"IN",
	"VN",
	"PK",
	"HL",
}

var AvailableLanguages = []string{
	"EN_US",
	"JA_JP",
	"VI_VN",
	"NE_NP",
}

func ValidateKataKana(fl validator.FieldLevel) bool {
	katakana := fl.Field().String()

	return kana.IsKatakana(katakana)
}

func ValidateCountryCode(fl validator.FieldLevel) bool {
	country := fl.Field().String()

	return slices.Contains(AvailableCountryCodes, strings.ToUpper(country))
}

func ValidateCountryCodeSlice(fl validator.FieldLevel) bool {
	countries := fl.Field().Interface().([]string)

	for _, country := range countries {
		if !slices.Contains(AvailableCountryCodes, strings.ToUpper(country)) {
			return false
		}
	}

	return true
}

func ValidateLanguageCode(fl validator.FieldLevel) bool {
	language := fl.Field().String()

	return slices.Contains(AvailableLanguages, strings.ToUpper(language))
}

func ValidateLanguageCodeSlice(fl validator.FieldLevel) bool {
	languages := fl.Field().Interface().([]string)

	for _, language := range languages {
		if !slices.Contains(AvailableLanguages, strings.ToUpper(language)) {
			return false
		}
	}

	return true
}

type PhoneCountry struct {
	Name *string `json:"name"`
}

type PhoneNumberValidationResponse struct {
	Valid   bool          `json:"valid"`
	Country *PhoneCountry `json:"country"`
}

func ValidateJapanesePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	regexp := regexp.MustCompile(`^(0([1-9]{1}-?[1-9]\d{3}|[1-9]{2}-?\d{3}|[1-9]{2}\d{1}-?\d{2}|[1-9]{2}\d{2}-?\d{1})-?\d{4}|0[789]0-?\d{4}-?\d{4}|050-?\d{4}-?\d{4})$`)
	isvalidpattern := regexp.MatchString(phoneNumber)

	if !isvalidpattern {
		return false
	}

	apiKey := GetEnv("ABSTRACT_API_KEY", "")

	resp, err := http.Get(fmt.Sprintf("https://phonevalidation.abstractapi.com/v1/?api_key=%s&phone=+81%s", apiKey, phoneNumber))
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var phoneResponse PhoneNumberValidationResponse

	if err := json.Unmarshal(body, &phoneResponse); err != nil {
		return false
	}

	if !phoneResponse.Valid {
		return false
	}

	if phoneResponse.Country.Name == nil {
		return false
	}

	return *phoneResponse.Country.Name == "Japan"
}
