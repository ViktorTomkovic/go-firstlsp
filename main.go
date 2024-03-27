package main

import (
	"bufio"
	"log"
	"os"
	"github.com/ViktorTomkovic/go-firstlsp/rpc"
)

func main() {
	logger := getLogger("firstLSP.log")
	logger.Println("Started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("Cannot create logger")
	}
	return log.New(logfile, "[firstLSP]", log.Ldate|log.Ltime|log.Lshortfile)
}


