package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/ViktorTomkovic/go-firstlsp/analysis"
	"github.com/ViktorTomkovic/go-firstlsp/lsp"
	"github.com/ViktorTomkovic/go-firstlsp/rpc"
)

func main() {
	logger := getLogger("firstLSP.log")
	logger.Println("Started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, &state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state *analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Cannot parse initialize request: %s", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		// reply to init request
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		_, err := writer.Write([]byte(reply))
		if err != nil {
			logger.Printf("Could not reply init: %s", err)
		}
		logger.Printf("Sent response initialize")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Cannot parse didOpen notification: %s", err)
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Cannot parse didChange notification: %s", err)
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}

	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("Cannot create logger")
	}
	return log.New(logfile, "[firstLSP]", log.Ldate|log.Ltime|log.Lshortfile)
}
