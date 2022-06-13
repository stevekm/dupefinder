package finder

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

var numWorkers int = 4

type HashResult struct {
	Entry FileHashEntry
	Err   error
}

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

// handle the file opening and closing in order to get the file hash
func GetFileHash(fileEntry FileEntry) (FileHashEntry, error) {
	file, err := os.Open(fileEntry.Path)
	// if file read permission is denied, skip this file
	if os.IsPermission(err) {
		// logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", err)
		return FileHashEntry{}, err
	}

	if err != nil {
		// logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
		return FileHashEntry{}, err
	}
	hash := getFileMD5(file)
	file.Close()

	fileHashEntry := FileHashEntry{File: fileEntry, Hash: hash}
	return fileHashEntry, err
}

// find files that have the same hash value
func FindHashDupes(fileMap map[int64][]FileEntry) map[string][]FileHashEntry {
	hashesMap := map[string][]FileHashEntry{}

	// set up for concurrent parallel processing of file hashing
	// https://stackoverflow.com/questions/71458290/how-to-batch-dealing-with-files-using-goroutine/71458664#71458664
	runtime.GOMAXPROCS(numWorkers)
	work := make(chan FileEntry)
	results := make(chan HashResult)
	// create worker goroutines
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileEntry := range work {
				fileHashEntry, err := GetFileHash(fileEntry)
				result := HashResult{Entry: fileHashEntry, Err: err}
				results <- result
			}
		}()
	}

	// send the work to the workers
	// this happens in a goroutine in order
	// to not block the main function, once
	// all 5 workers are busy
	go func() {
		for _, entries := range fileMap {
			for _, entry := range entries {
				work <- entry
			}
		}
		// close the work channel after
		// all the work has been send
		close(work)

		// wait for the workers to finish
		// then close the results channel
		wg.Wait()
		close(results)
	}()

	// collect the results
	// the iteration stops if the results
	// channel is closed and the last value
	// has been received
	for result := range results {
		// allResults = append(allResults, result)
		if os.IsPermission(result.Err) {
			logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", result.Err)
			continue
		}

		if result.Err != nil {
			logger.Printf("WARNING: Skipping file that could not be opened: %v\n", result.Err)
			continue
		}
		hashesMap[result.Entry.Hash] = append(hashesMap[result.Entry.Hash], result.Entry)
	}

	dupesMap := map[string][]FileHashEntry{}
	for hash, entries := range hashesMap {
		if len(entries) > 1 {
			dupesMap[hash] = entries
		}
	}
	return dupesMap
}
