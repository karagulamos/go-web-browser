package network

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Headers map[string]string

func (h Headers) Get(key string) string {
	return h[key]
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

type HttpResponse struct {
	StatusCode int
	StatusText string
	Headers    Headers
	Body       io.ReadCloser
}

func AsHttpResponse(socket Socket) (HttpResponse, error) {
	response := HttpResponse{}

	statusLine, err := socket.ReadLine()
	if err != nil {
		return response, err
	}

	statusTokens := strings.SplitN(statusLine, " ", 3)
	if len(statusTokens) < 3 {
		return response, fmt.Errorf("invalid status line")
	}

	response.StatusCode, _ = strconv.Atoi(statusTokens[1])
	response.StatusText = statusTokens[2]

	headers, err := parseHeaders(socket)
	if err != nil {
		return response, err
	}

	response.Headers = headers
	response.Body = socket

	return response, nil
}

func parseHeaders(socket Socket) (Headers, error) {
	headers := make(Headers)

	for {
		line, err := socket.ReadLine()
		if err != nil {
			return nil, err
		}

		if line == "" {
			break
		}

		if keyValue := strings.SplitN(line, ": ", 2); len(keyValue) == 2 {
			headers.Set(strings.ToLower(keyValue[0]), keyValue[1])
		}
	}

	return headers, nil
}
