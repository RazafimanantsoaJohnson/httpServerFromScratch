package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Header map[string]string

func NewHeaders() Header {
	return map[string]string{}
}

func (h Header) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// done reading
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace((parts[1]))
	key = strings.TrimSpace(key)
	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}
	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Header) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func (h Header) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

var validTokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validTokens(data []byte) bool {
	for _, c := range data {
		if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-') {
			return false
		}
	}
	return true
}
