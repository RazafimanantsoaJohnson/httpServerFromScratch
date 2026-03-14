package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/RazafimanantsoaJohnson/httpServer/internal/headers"
)

type ParserStatus int

const (
	Initialized ParserStatus = iota
	ParsingHeader
	ParsingBody
	Done
)

const newLineCharacter = "\r\n"

type Request struct {
	RequestLine RequestLine
	Headers     headers.Header
	Body        []byte
	ParseStatus ParserStatus
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

const bufferInitialSize = 4096
const contentLengthHeader = "Content-Length"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferInitialSize)
	readToIndex := 0
	result := Request{
		ParseStatus: Initialized,
		Headers:     headers.NewHeaders(),
		Body:        make([]byte, 0),
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
				if result.ParseStatus != Done {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", result.ParseStatus, numBytesRead)
				}
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
	totalBytesRead := 0
	for r.ParseStatus != Done {
		if totalBytesRead > len(data) {
			return totalBytesRead, fmt.Errorf("parser consumed beyond buffer: consumed=%d len=%d", totalBytesRead, len(data))
		}
		n, err := r.parseSingleLine(data[totalBytesRead:])

		if err != nil {
			return totalBytesRead, err
		}
		if n == 0 {
			return totalBytesRead, err
		}
		totalBytesRead += n
	}

	return totalBytesRead, nil
}

func (r *Request) parseSingleLine(data []byte) (int, error) {
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
		r.ParseStatus = ParsingHeader
		return n, nil
	case ParsingHeader:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return n, nil
		}
		if !done {
			return n, nil
		}
		r.ParseStatus = ParsingBody
		return n, nil
	case ParsingBody:
		contentLengthStr, isHeaderPresent := r.Headers.Get(contentLengthHeader)
		fmt.Printf(contentLengthStr)
		if !isHeaderPresent {
			r.ParseStatus = Done
			return len(data), nil
		}
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("invalid content length header")
		}

		r.Body = append(r.Body, data...)
		if len(r.Body) > contentLength {
			return len(data), fmt.Errorf("content-length should be equal to the length of body")
		} else if len(r.Body) == contentLength {
			r.ParseStatus = Done
		}

		return len(data), nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
