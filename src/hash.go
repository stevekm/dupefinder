package finder

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/cespare/xxhash" //https://pkg.go.dev/github.com/cespare/xxhash#section-readme
	"hash"
	"io"
	"os"
	"runtime"
	"sync"
)

// var numWorkers int = 4
type HashConfig struct {
	NumWorkers int
	NumBytes   int64
	Partial    bool
	Algo       string
	Verbose    bool //false by default
}

type HashResult struct {
	Entry FileHashEntry
	Err   error
}

// get the md5 hash of an open file handle
// https://stackoverflow.com/questions/1761607/what-is-the-fastest-hash-algorithm-to-check-if-two-files-are-equal
func getFileMD5(inputFile *os.File, config HashConfig) string {
	algo := config.Algo
	if algo == "" {
		algo = "md5"
	}

	var hashWriter hash.Hash
	switch {
	case algo == "md5":
		hashWriter = md5.New()
	case algo == "sha1":
		hashWriter = sha1.New()
	case algo == "sha256":
		hashWriter = sha256.New()
	case algo == "xxhash":
		hashWriter = xxhash.New()
	default:
		hashWriter = md5.New()
	}

	// optionally hash only part of the file
	if (config.Partial) && (config.NumBytes > 0) {
		numBytesCopied, err := io.CopyN(hashWriter, inputFile, config.NumBytes)
		if err != nil {
			// if we are hashing n bytes then a lot of files will be too small so handle EOF
			if err == io.EOF {
				// dont print this it floods the terminal
				// logger.Printf("Hashed %v bytes from file %v when %v bytes were wanted; continuing...\n", numBytesCopied, inputFile.Name(), config.NumBytes)
			} else {
				logger.Fatalf("Error encountered while hashing %v bytes from file: %v\n", numBytesCopied, err)
			}
		}

	} else {
		_, err := io.Copy(hashWriter, inputFile)
		if err != nil {
			logger.Fatalf("Error encountered while hashing file: %v\n", err)
		}
	}

	sum := hashWriter.Sum(nil)
	hashStr := hex.EncodeToString(sum[:])
	return hashStr
}

// handle the file opening and closing in order to get the file hash
func GetFileHash(fileEntry FileEntry, config HashConfig) (FileHashEntry, error) {
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
	hash := getFileMD5(file, config)
	file.Close()

	fileHashEntry := FileHashEntry{File: fileEntry, Hash: hash}
	return fileHashEntry, err
}

// find files that have the same hash value
func FindHashDupes(fileMap map[int64][]FileEntry, hashConfig HashConfig) map[string][]FileHashEntry {
	hashesMap := map[string][]FileHashEntry{}
	var numFilesHashed int

	// set up for concurrent parallel processing of file hashing
	// https://stackoverflow.com/questions/71458290/how-to-batch-dealing-with-files-using-goroutine/71458664#71458664
	var numWorkers int
	if hashConfig.NumWorkers > 0 {
		numWorkers = hashConfig.NumWorkers
	} else {
		numWorkers = 1
	}

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
				if hashConfig.Verbose {
					logger.Printf("Hashing %v\n", fileEntry.Path)
				}
				fileHashEntry, err := GetFileHash(fileEntry, hashConfig)
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
		numFilesHashed += 1
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

	if hashConfig.Verbose {
		logger.Printf("Hashed %v files\n", numFilesHashed)
	}

	dupesMap := map[string][]FileHashEntry{}
	var numHashDupes int
	for hash, entries := range hashesMap {
		if len(entries) > 1 {
			dupesMap[hash] = entries
			numHashDupes += len(entries)
		}
	}

	if hashConfig.Verbose {
		logger.Printf("Found %v hash duplicates\n", numHashDupes)
	}
	return dupesMap
}
