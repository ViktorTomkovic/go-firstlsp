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
func DecodeMessage(msg []byte) (string, int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", 0, errors.New("Did not find separator")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", 0, err
	}
	fmt.Println(content)
	fmt.Println(contentLength)
	fmt.Println(content[:contentLength])
	var baseMessage BaseMessage
	err = json.Unmarshal(content[:contentLength], &baseMessage)
	if err != nil {
		return "", 0, nil
	}
	return baseMessage.Method, contentLength, nil
}
