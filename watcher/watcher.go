package watcher

import (
	"fileserv/handlers"
	"fileserv/indexer"
	"fmt"
	"reflect"
)

func FileWatcher(h *handlers.Handler) {
	newFilesList, err := indexer.Indexfiles(h.BaseDir)
	if err != nil {
		fmt.Printf("error while indexing file %v", err)
	}

	h.Mu.RLock()
	unchanged := reflect.DeepEqual(newFilesList, h.FilesList)
	h.Mu.RUnlock()

	if unchanged {
		fmt.Println("no changes detected.")
		return
	}

	h.Mu.Lock()
	h.FilesList = newFilesList
	h.Mu.Unlock()
	fmt.Println("file list updated")
}
