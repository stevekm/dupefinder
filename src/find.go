package finder

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"runtime"
	"fmt"
)

type FindConfig struct {
	NumWorkers int
}

type FindResult struct {
	File FileEntry
	Err error
}

// get only the top-level subdirectories in the root path and bare files
// so that they can be passed down for parallel processing
func FindRootSubdirsFiles(dirPath string) ([]string, []FileEntry) {
	startingDepth := strings.Count(dirPath, string(os.PathSeparator))
	dirs := []string{}
	fileEntries := []FileEntry{}

	// First dir walk to get the root subdirs
	maxSubdirDepth := startingDepth + 1
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

		// skip it if its the original dir
		if path == dirPath {
			// return filepath.SkipDir
			return nil
		}

		// skip and sub-subdirs
		// fmt.Printf("FindSubdirs: %v %v\n", path, strings.Count(path, string(os.PathSeparator)))
		if info.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxSubdirDepth {
			return filepath.SkipDir
		}

		// collect dirs
		if info.IsDir() {
			dirs = append(dirs, path)
			// return nil
		}
		// else if info.Mode().IsRegular() { // collect files
		// 	fileEntry := NewFileEntryFromPathInfo(path, info)
		// 	fileEntries = append(fileEntries, fileEntry)
		// 	// return nil
		// }
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dirPath, err)
	}




	// second walk to get the root files
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
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

		if path == dirPath {
			return nil
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		// collect files
		if info.Mode().IsRegular() {
			fileEntry := NewFileEntryFromPathInfo(path, info)
			fileEntries = append(fileEntries, fileEntry)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path %q: %v\n", dirPath, err)
	}

	return dirs, fileEntries
}


func FindFilesSizesParallel(dirPath string, skipDirs []string, config FindConfig) (map[int64][]FileEntry, uint64) {
	fileMap := map[int64][]FileEntry{}
	var numFiles uint64

	// get root files and subdirs
	subDirs, rootFiles := FindRootSubdirsFiles(dirPath)

	for _, file := range rootFiles {
		fileMap[file.Size] = append(fileMap[file.Size], file)
		numFiles += 1
	}

	// set up for concurrent parallel processing of file hashing
	var numWorkers int = 1
	if config.NumWorkers > 0 {
		numWorkers = config.NumWorkers
	}
	runtime.GOMAXPROCS(numWorkers)
	work := make(chan string) // subdirs go in here
	results := make(chan FileEntry) // FileEntry find results come out here
	// create worker goroutines
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for subDir := range work {
				fmt.Printf(">>> Starting FindFilesSizes for subDir: %v\n", subDir)
				fileMap, _ := FindFilesSizes(subDir, skipDirs)
				for size, files := range fileMap {
					for _, file := range files {
						results <- file
					}
				}
			}
		}()
	}


}


// find all files in the directory tree and group them by file size
func FindFilesSizes(dirPath string, skipDirs []string) (map[int64][]FileEntry, uint64) {
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
func FindDupes(dirPath string, skipDirs []string, hashConfig HashConfig) map[string][]FileHashEntry {
	fileSizeMap, _ := FindFilesSizes(dirPath, skipDirs)
	hashDupes := FindHashDupes(fileSizeMap, hashConfig)
	return hashDupes
}
