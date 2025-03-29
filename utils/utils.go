package utils

import (
	"errors"
	"strings"
	"unicode"
)

func isAlphanumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isValidSwiftCode(swiftCode string) bool {
	return (len(swiftCode) == 11) && isAlphanumeric(swiftCode)
}

func IsHeadquarter(swiftCode string) (bool, error) {
	if !isValidSwiftCode(swiftCode) {
		return false, errors.New("Invalid swift code " + swiftCode)
	}
	return strings.HasSuffix(swiftCode, "XXX"), nil
}
