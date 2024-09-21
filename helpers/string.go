package helpers

import (
	"strings"
)

func ToSnakeCase(s string) string {
	splitted := strings.Split(s, "")
	newString := ""
	for _, v := range splitted {
		if strings.ToUpper(v) == v {
			newString += "_" + strings.ToLower(v)
		} else {
			newString += v
		}
	}

	return newString
}

func Capitalise(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ToCamelCase(s string, upperAll bool) string {
	splitted := strings.Split(s, "_")
	newString := ""
	for i, v := range splitted {
		if !upperAll && i == 0 {
			newString += strings.ToLower(v)
			continue
		}
		newString += strings.ToUpper(string(v[0])) + v[1:]
	}

	return newString
}
