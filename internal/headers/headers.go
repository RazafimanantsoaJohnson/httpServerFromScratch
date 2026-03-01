package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Header map[string]string

const newLineChar = "\r\n"

func (h Header) Parse(data []byte) (n int, done bool, err error) {

	keyValueLines := bytes.Split(data, []byte(newLineChar))

	if len(keyValueLines) <= 0 {
		return 0, false, nil
	}
	if len(keyValueLines) == 1 || string(keyValueLines[0]) == "" {
		return 0, true, nil
	}
	firstHeader := string(keyValueLines[0])
	n = len(keyValueLines[0])
	firstHeader = strings.TrimSpace(firstHeader)
	firstColonIndex := strings.Index(firstHeader, ":")
	isHeaderValid := (strings.Index(firstHeader, " ") > firstColonIndex)
	if !isHeaderValid {
		return 0, false, fmt.Errorf("wrong header format; spaces after the 'key'")
	}
	key := firstHeader[:firstColonIndex]
	value := strings.TrimSpace(firstHeader[firstColonIndex+1:])
	isKeyValid := isFieldNameValid(key)
	fmt.Println("Is key valid: ", isKeyValid)
	if !isKeyValid {
		return 0, false, fmt.Errorf("the key '%s' is not valid", key)
	}
	headerKey := strings.ToLower(key)
	currentKeyValue, doesKeyExist := h[strings.ToLower(headerKey)]
	if doesKeyExist {
		h[headerKey] = fmt.Sprintf("%s,%s", currentKeyValue, value)
	} else {
		h[strings.ToLower(key)] = value
	}

	return n, false, nil
}

func isFieldNameValid(fieldName string) bool {
	validCharacters := []rune{}

	for c := rune('a'); c <= rune('z'); c++ {
		validCharacters = append(validCharacters, c)
	}
	for c := rune('A'); c <= rune('Z'); c++ {
		validCharacters = append(validCharacters, c)
	}
	for c := rune('0'); c <= rune('9'); c++ {
		validCharacters = append(validCharacters, c)
	}

	punct := "!#$%&'*+-.^_`|~"
	for _, c := range punct {
		validCharacters = append(validCharacters, c)
	}
	if len(fieldName) <= 0 {
		return false
	}
	for _, fieldChar := range fieldName {
		if !(slices.Contains(validCharacters, rune(fieldChar))) {
			return false
		}
	}
	return true
}
