package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type ParserStatus int

const (
	Initialized ParserStatus = iota
	Done
)

const newLineCharacter = "\r\n"

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	ParseStatus ParserStatus
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

const bufferInitialSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferInitialSize)
	readToIndex := 0
	result := Request{
		ParseStatus: Initialized,
	}
	for result.ParseStatus != Done {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				result.ParseStatus = Done
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := result.parse(buffer[:readToIndex]) //we keep calling parse but it's only really parsing when there is a full line with '\r\n'
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return &result, nil
}

func parseRequestLine(rawData []byte) (RequestLine, int, error) {
	newLineIndex := bytes.Index(rawData, []byte(newLineCharacter))
	numBytesConsumed := 0
	result := RequestLine{}

	if newLineIndex == -1 {
		return result, numBytesConsumed, nil
	}
	reqLineString := string(rawData[:newLineIndex])
	result, err := extractRequestLineFromString(reqLineString)
	if err != nil {
		return result, numBytesConsumed, err
	}
	numBytesConsumed = newLineIndex + 2

	return result, numBytesConsumed, nil
}

func extractRequestLineFromString(line string) (RequestLine, error) {
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

func (r *Request) parse(data []byte) (int, error) {
	// when the parseRequestLine return a 'numberOfBytes > 0' we read the full line.
	switch r.ParseStatus {
	case Initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.ParseStatus = Done
		return n, nil
	case Done:
		return 0, fmt.Errorf("error: trying to resad data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
