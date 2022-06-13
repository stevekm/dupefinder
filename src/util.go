package finder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var logger = log.New(os.Stderr, "", 0)

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

// walk the directory tree recursively to find all files
// skip dirs that are in the skipDirs list
func GetFiles(dirPath string, skipDirs []string) []FileEntry {
	fileList := []FileEntry{}
	// https://pkg.go.dev/path/filepath#Walk
	// TODO: look into using https://pkg.go.dev/io/fs#WalkDirFunc , https://pkg.go.dev/path/filepath#WalkDir
	// https://github.com/golang/go/issues/16399
	// https://pkg.go.dev/io/fs#DirEntry
	// NOTE: dont think WalkDir will save much time since I want the extra info anyway??
	// https://pkg.go.dev/io/fs#FileInfo
	// https://golang.hotexamples.com/examples/os/-/IsPermission/golang-ispermission-function-examples.html
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		// skip item that cannot be read
		if os.IsPermission(err) {
			logger.Printf("Skipping path that could not be read %q: %v\n", path, err)
			return filepath.SkipDir
		}

		if err != nil {
			logger.Printf("Error encountered when accessing path %q: %v\n", path, err)
			return err
		}
		// skip some dirs
		if info.IsDir() &&
			containsStr(skipDirs, info.Name()) ||
			containsStr(skipDirs, path) {
			logger.Printf("skipping a dir: %+v %v \n", info.Name(), path)
			return filepath.SkipDir
		}

		// if its a file then add it to the list
		// https://pkg.go.dev/io/fs#FileMode.IsRegular
		if info.Mode().IsRegular() {
			file := FileEntry{
				Path: path,
				Name: info.Name(),
				Size: info.Size(),
			}
			fileList = append(fileList, file)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dirPath, err)
		// return
	}

	return fileList
}

// check if a slice contains a specific string
// TODO: update to Go 1.18 so we dont have to do this anymore
func containsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// group all files with the same size value
func GroupBySize(fileList []FileEntry) map[int64][]FileEntry {
	sizes := map[int64][]FileEntry{}
	for _, fileEntry := range fileList {
		sizes[fileEntry.Size] = append(sizes[fileEntry.Size], fileEntry)
	}
	return sizes
}

// find size values that have multiple associated files
func FindSizeDupes(fileList []FileEntry) []FileEntry {
	sizeMap := GroupBySize(fileList)
	sizeDupesList := []FileEntry{}
	for _, v := range sizeMap {
		if len(v) > 1 {
			for _, fileEntry := range v {
				sizeDupesList = append(sizeDupesList, fileEntry)
			}
		}
	}
	return sizeDupesList
}

// arrange file list into groupings based on hash value
func GroupByHash(fileList []FileEntry) map[string][]FileEntry {
	hashMap := map[string][]FileEntry{}
	for _, entry := range fileList {
		file, err := os.Open(entry.Path)
		// if file read permission is denied, skip this file
		if os.IsPermission(err) {
			logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", err)
			continue
		}

		if err != nil {
			logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
			continue
		}
		hash := getFileMD5(file)
		file.Close()
		hashMap[hash] = append(hashMap[hash], entry)
	}
	return hashMap
}

// handle the file opening and closing in order to get the file hash
func GetFileHash(fileEntry FileEntry) (FileHashEntry, error) {
	file, err := os.Open(fileEntry.Path)
	// if file read permission is denied, skip this file
	if os.IsPermission(err) {
		logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", err)
		// continue
	}

	if err != nil {
		logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
		// continue
	}
	hash := getFileMD5(file)
	file.Close()

	fileHashEntry := FileHashEntry{File: fileEntry, Hash: hash}
	return fileHashEntry, err
}

// get all hashes for all files in the list
func GetFilesHashes(fileList []FileEntry) []FileHashEntry {
	fileHashes := []FileHashEntry{}
	for _, entry := range fileList {
		fileHashEntry, err := GetFileHash(entry)
		if err != nil {
			// logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
			continue
		}
		fileHashes = append(fileHashes, fileHashEntry)
	}
	return fileHashes
}

// find files that have the same hash value
func FindHashDupes(fileList []FileEntry) map[string][]FileEntry {
	hashMap := GroupByHash(fileList)
	hashDupes := map[string][]FileEntry{}
	for hash, entries := range hashMap {
		if len(entries) > 1 {
			for _, entry := range entries {
				hashDupes[hash] = append(hashDupes[hash], entry)
			}
		}
	}
	return hashDupes
}

// find all the duplicate files in the dir
// Duplicates = same file size, same hash value
// TODO: this might need to be broken up to aid garbage collection ??
func FindDupes(dirPath string, skipDirs []string) map[string][]FileEntry {
	fileList := GetFiles(dirPath, skipDirs)

	// find sizes with multiple files
	sizeDupesList := FindSizeDupes(fileList)

	// reduce the list to only the entries with multiple files with the same hash
	hashDupes := FindHashDupes(sizeDupesList)

	return hashDupes
}
