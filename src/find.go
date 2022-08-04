package finder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type FindConfig struct {
	MinSize  int64
	MaxSize *int64 // zero value nil allows to check if value was set
	SkipDirs []string
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

func containsFileHashEntry(l []FileHashEntry, e FileHashEntry) bool {
	for _, a := range l {
		if a == e {
			return true
		}
	}
	return false
}

// find all files in the directory tree and group them by file size
func FindFilesSizes(dirPath string, config FindConfig) (map[int64][]FileEntry, uint64) {
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
			containsStr(config.SkipDirs, info.Name()) ||
			containsStr(config.SkipDirs, path) {
			logger.Printf("skipping a dir: %+v %v \n", info.Name(), path)
			return filepath.SkipDir
		}

		// if its a file then add it to the list
		if info.Mode().IsRegular() {
			// test for file size filters
			size := info.Size()
			// default to false
			var passMinSize bool
			var passMaxSize bool
			var passSize bool

			if size >= config.MinSize {
				passMinSize = true
			}

			// MaxSize automatically passes if no value was given
			if config.MaxSize == nil {
				passMaxSize = true
				} else {
					if size <= *config.MaxSize {
						passMaxSize = true
				}
			}

			if passMinSize && passMaxSize {
				passSize = true
			}

			if passSize {
				fileEntry := NewFileEntryFromPathInfo(path, info)
				fileMap[size] = append(fileMap[size], fileEntry)
				numFiles += 1
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dirPath, err)
	}

	return fileMap, numFiles
}

func FindSizeDupes(fileSizeMap map[int64][]FileEntry) map[int64][]FileEntry {
	dupesMap := map[int64][]FileEntry{}
	for key, entries := range fileSizeMap {
		if len(entries) > 1 {
			dupesMap[key] = entries
		}
	}
	return dupesMap
}

// find all the duplicate files in the dir
// Duplicates = same file size, same hash value
// TODO: this might need to be broken up to aid garbage collection ??
func FindDupes(dirPath string, findConfig FindConfig, hashConfig HashConfig) map[string][]FileHashEntry {
	fileSizeMap, _ := FindFilesSizes(dirPath, findConfig)
	sizeDupes := FindSizeDupes(fileSizeMap)
	hashDupes := FindHashDupes(sizeDupes, hashConfig)
	return hashDupes
}
