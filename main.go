package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"io"

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
	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, writer, &state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Cannot parse initialize request: %s", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		// reply to init request
		response := lsp.NewInitializeResponse(request.ID)
		writeResponse(logger, writer, response)
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
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
		}
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(logger, writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
		}
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(logger, writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return
		}
		response := state.CodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(logger, writer, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/completion: %s", err)
			return
		}
		response := state.Completion(request.ID, request.Params.TextDocument.URI)
		writeResponse(logger, writer, response)
	}
}

func writeResponse(logger *log.Logger, writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	_, err := writer.Write([]byte(reply))
	if err != nil {
		logger.Printf("Could not reply: %s", err)
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("Cannot create logger")
	}
	return log.New(logfile, "[firstLSP]", log.Ldate|log.Ltime|log.Lshortfile)
}
