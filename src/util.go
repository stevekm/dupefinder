package finder

import (
	"fmt"
	"log"
	"os"
	"io"
	"path/filepath"
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"strings"
)

// get the md5 hash of an open file handle
func getFileMD5(inputFile *os.File) string {
	hash := md5.New()
	if _, err := io.Copy(hash, inputFile); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)
	hashStr := hex.EncodeToString(sum[:])
	return hashStr
}

type FileEntry struct {
	Path string
	Name string // basename of the file
	Size int64
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
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		// skip some dirs
		if info.IsDir() &&
			containsStr(skipDirs, info.Name()) ||
			containsStr(skipDirs, path) {
			fmt.Printf("skipping a dir: %+v %v \n", info.Name(), path)
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

// find all the duplicate files in the dir
// Duplicates = same file size, same hash value
// TODO: this might need to be broken up to aid garbage collection ??
func FindDupes(dirPath string, skipDirs []string) map[string][]string {
	fileList := GetFiles(dirPath, skipDirs)

	// first group by files with the same size
	sizes := map[int64][]FileEntry{}
	for _, fileEntry := range fileList {
		sizes[fileEntry.Size] = append(sizes[fileEntry.Size], fileEntry)
	}

	// find sizes with multiple files
	sizeDupes := [][]FileEntry{}
	for _, v := range sizes {
		names := []FileEntry{}
		if len(v) > 1 {
			for _, fileEntry := range v {
				names = append(names, fileEntry)
			}
			sizeDupes = append(sizeDupes, names)
		}
	}

	// check the hashes to determine if they are actually duplicates
	hashes := map[string][]FileEntry{}
	for _, entries := range sizeDupes {
		for _, entry := range entries {
			file, err := os.Open(entry.Path)
			if err != nil {
				log.Fatalf("error opening the path %v\n", err)
			}
			hash := getFileMD5(file)
			hashes[hash] = append(hashes[hash], entry)
		}
	}

	// reduce the list to only the entries with multiple files with the same hash
	hashDupes := map[string][]string{}
	for hash, entries := range hashes {
		if len(entries) > 1 {
			for _, entry := range entries {
				hashDupes[hash] = append(hashDupes[hash], entry.Path)
			}
		}
	}

	return hashDupes
}

func DupesFormatter (hash string, dupes []string) string {
	var outputStr string
	outputStr = hash + "\n"
	outputStr += strings.Join(dupes, "\n")
	outputStr += "\n"
	return outputStr
}
