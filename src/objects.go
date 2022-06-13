package finder

import (
	"os"
	"log"
)

// basic file entry
type FileEntry struct {
	Path string
	Name string // basename of the file
	Size int64
}

// file entry with hash
type FileHashEntry struct {
	File FileEntry
	Hash string
}

// method for creating a new FileEntry when we have only the filepath available
func NewFileEntryFromPath(filepath string) FileEntry {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("error opening the path %v\n", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	entry := FileEntry{
		Path: filepath,
		Name: info.Name(),
		Size: info.Size(),
	}

	return entry
}
