package finder

import (
	"io/fs"
	"log"
	"os"
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

// use this to create FileEntry if file info has already been called
func NewFileEntryFromPathInfo(filepath string, fileinfo fs.FileInfo) FileEntry {
	entry := FileEntry{
		Path: filepath,
		Name: fileinfo.Name(),
		Size: fileinfo.Size(),
	}

	return entry
}

func NewFileHashEntry(fileEntry FileEntry) FileHashEntry {
	fileHashEntry, err := GetFileHash(fileEntry)
	if err != nil {
		log.Fatalf("Could not convert FileEntry to FileHashEntry %v\n", err)
	}
	return fileHashEntry
}
