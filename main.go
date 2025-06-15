package main

import (
	"bufio"
	"dalsp/analysis"
	"dalsp/lsp"
	"dalsp/rpc"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	logger := getLogger("/home/vaibhav/Files/dev/projects/dalsp/log.txt")
	logger.Println("Hey, i started!")

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
		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Recieved msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldnt parse this: %s", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewIntializeResponse(request.ID)

		writeResponse(writer, msg)

		logger.Printf("Sent the reply")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}
		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
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

		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("hey, you didnt give me a good file")
	}

	return log.New(logfile, "[daLSP]", log.Ldate|log.Ltime|log.Lshortfile)
}
