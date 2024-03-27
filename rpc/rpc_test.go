package rpc_test

import (
	"testing"
	"github.com/ViktorTomkovic/go-firstlsp/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing:true})
	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	_, contentLength, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}
	if contentLength != 16 {
		t.Fatalf("Expected: 16, Got: %d", contentLength)
	}
}

func TestDecodeMethod(t *testing.T) {
	incomingMessage := "Content-Length: 36\r\n\r\n{\"Method\":\"textDocument/completion\"}"
	method, contentLength, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}
	if contentLength != 36 {
		t.Fatalf("Expected: 36, Got: %d", contentLength)
	}
	if method != "textDocument/completion" {
		t.Fatalf("Expected: \"textDocument/completion\", Got: %s", method)
	}
}
