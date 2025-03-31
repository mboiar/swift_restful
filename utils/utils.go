// Utils implements utility functions for processing SWIFT data.
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

func IsValidSwiftCode(swiftCode string) bool {
	return (len(swiftCode) == 11) && isAlphanumeric(swiftCode)
}

func IsHeadquarter(swiftCode string) (bool, error) {
	if !IsValidSwiftCode(swiftCode) {
		return false, errors.New("Invalid swift code " + swiftCode)
	}
	return strings.HasSuffix(swiftCode, "XXX"), nil
}
