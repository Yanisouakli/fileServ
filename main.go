package main

import (
	"encoding/json"
	"fileserv/handlers"
	"fileserv/indexer"
	"fileserv/watcher"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var workingDir string = "."

func main() {
	args := os.Args
	if len(args) > 1 {
		workingDir = os.Args[1]
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)

	filesList, err := indexer.Indexfiles(workingDir)
	if err != nil {
		fmt.Printf("error when indexing file :%v", err)
	}
	mux := http.NewServeMux()
	handler := handlers.New(workingDir, filesList)

	go func() {
		<-sigs
		fmt.Println("program will shutdown...")
		fmt.Println("performing backup")
		jsonToWrite, err := json.MarshalIndent(handler.FilesList, "", " ")
		if err != nil {
			fmt.Println("error while marshalling files list")
		}
		Werr := os.WriteFile("files_backup.json", jsonToWrite, 0644)
		if Werr != nil {
			fmt.Println("error while creating the backup")
		}
		done <- true
	}()

	go func() {
		for range ticker.C {
			watcher.FileWatcher(handler)
		}
	}()

	mux.HandleFunc("/files", handler.ListFilesHandler)
	mux.HandleFunc("/search", handler.SearchFileHandler)

	log.Println("file indexer server running on 8080")

	go func() {
		ServerErr := http.ListenAndServe(":8080", mux)
		if ServerErr != nil {
			log.Printf("error while statring the server %v", ServerErr)
			done <- true
		}
	}()

	<-done
	fmt.Println("program exited gracefully")
}
