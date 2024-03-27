package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"errors"
	"bytes"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Did not find separator")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}
	fmt.Println(content)
	fmt.Println(contentLength)
	fmt.Println(content[:contentLength])
	var baseMessage BaseMessage
	err = json.Unmarshal(content[:contentLength], &baseMessage)
	if err != nil {
		return "", nil, nil
	}
	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}
	if len(content) < contentLength {
		// did not receive the whole message
		return 0, nil, nil
	}
	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
