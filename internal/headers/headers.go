package headers

import (
	"bytes"
	"fmt"
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
	h[key] = value

	fmt.Println("Header =>", h[key], ":", h[value])
	return n, false, nil
}
