package handlers

import (
	"sync"
	"time"
)

type FileMetaData struct {
	FileName string    `json:"file_name"`
	FilePath string    `json:"file_path"`
	ModTime  time.Time `json:"mod_time"`
	FileSize int       `json:"file_size"`
}

type Handler struct {
	FilesList []FileMetaData
	BaseDir   string
	Mu        sync.RWMutex
}

func New(dir string, filesList []FileMetaData) *Handler {
	return &Handler{FilesList: filesList, BaseDir: dir}
}
