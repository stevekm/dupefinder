package finder

import (
	"os"
	"path/filepath"
	"io/fs"
	"log"
)

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

// find all files in the directory tree and group them by file size
func FindFilesSizes (dirPath string, skipDirs []string) (map[int64][]FileEntry, uint64) {
	fileMap := map[int64][]FileEntry{}
	var numFiles uint64

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		// skip item that cannot be read
		if os.IsPermission(err) {
			logger.Printf("Skipping path that could not be read %q: %v\n", path, err)
			return filepath.SkipDir
		}
		// generic handling for other errors
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
		if info.Mode().IsRegular() {
			size := info.Size()
			fileEntry := NewFileEntryFromPathInfo(path, info)
			fileMap[size] = append(fileMap[size], fileEntry)
			numFiles += 1
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dirPath, err)
	}

	return fileMap, numFiles
}



// find all the duplicate files in the dir
// Duplicates = same file size, same hash value
// TODO: this might need to be broken up to aid garbage collection ??
func FindDupes(dirPath string, skipDirs []string) map[string][]FileHashEntry {
	fileSizeMap, _ := FindFilesSizes(dirPath, skipDirs)
	hashDupes := FindHashDupes(fileSizeMap)
	return hashDupes
}
