package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonFiles, err := json.MarshalIndent(h.FilesList, "", " ")
	if err != nil {
		fmt.Fprintf(w, "an error occured while marshaling files %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFiles)
}

func (h *Handler) SearchFileHandler(w http.ResponseWriter, r *http.Request) {
	var returnResult []FileMetaData
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	term := r.URL.Query().Get("q")
	if len(term) == 0 {
		http.Error(w, "Please provide a search term", http.StatusBadRequest)
		return
	}
	for _, file := range h.FilesList {
		if strings.Contains(file.FileName, term) {
			returnResult = append(returnResult, file)
		}
	}

	if len(returnResult) == 0 {
		http.Error(w, "no file mathes this string", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(returnResult); err != nil {
		log.Printf("Failed to encode response %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
