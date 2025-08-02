package indexer

import (
	"fileserv/handlers"
	"fmt"
	"os"
	"path/filepath"
)

func Indexfiles(dir string) ([]handlers.FileMetaData, error) {
	var filesList []handlers.FileMetaData
	var workingDir string = "."
	var pathPrefix string = ""
	if dir != "." {
		workingDir = dir + "/"
		pathPrefix = workingDir
	}
	entries, err := os.ReadDir(workingDir)
	if err != nil {
		fmt.Printf("error while reading file: %v", err)
		return nil, err
	}

	for _, entry := range entries {
		pathName, err := filepath.Abs(pathPrefix + entry.Name())
		if err != nil {
			fmt.Printf("Error getting file path: [[%v]]\n", err)
			return nil, err
		}

		fileInfo, err := os.Stat(pathName)
		if err != nil {
			fmt.Printf("Error getting file info:%s {{%v}}\n", pathName, err)
			return nil, err
		}
		fileData := handlers.FileMetaData{
			FileName: entry.Name(),
			FilePath: pathName,
			FileSize: int(fileInfo.Size()),
			ModTime:  fileInfo.ModTime(),
		}
		filesList = append(filesList, fileData)
	}
	return filesList, nil
}
