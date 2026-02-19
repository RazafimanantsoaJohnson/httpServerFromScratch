package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

func RequestFromHeader(reader io.Reader) (*Request, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	messageContent := string(content)
	httpMessage := strings.Split(messageContent, "\r\n")
	requestLine, err := parseRequestLine(httpMessage[0])
	if err != nil {
		return nil, err
	}

	result := Request{
		RequestLine: requestLine,
	}

	return &result, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	result := RequestLine{}
	requestParts := strings.Split(line, " ")

	if len(requestParts) != 3 {
		return result, fmt.Errorf("the request line has a wrong format")
	}
	method, endpoint, version := requestParts[0], requestParts[1], requestParts[2]

	if method != strings.ToUpper(method) {
		return result, fmt.Errorf("unsupported HTTP method")
	}

	versionParts := strings.Split(version, "/")
	if len(versionParts) != 2 {
		return result, fmt.Errorf("unsupported HTTP version")
	} else if versionParts[1] != "1.1" {
		return result, fmt.Errorf("unsupported HTTP version")
	}

	result = RequestLine{
		Method:        method,
		HttpVersion:   versionParts[1],
		RequestTarget: endpoint,
	}

	return result, nil
}
