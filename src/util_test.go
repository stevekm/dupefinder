package finder

import (
	"runtime"
	"sync"
	"testing"
	"time"
)


// NOTE: What is this test even for?? Looks like some old prototype code..
// Looks like it just creates a ton of files and then hashes them all
func TestLargeFileHandling(t *testing.T) {
	skipDirTest(t)
	tempdir := t.TempDir()
	t.Run("Test create large file", func(t *testing.T) {
		fileEntries := []FileEntry{}
		for _, v := range []int64{3e8, 4e8, 1e8, 1e7, 4e5, 8e4, 3e8, 4e8} {
			tempfile, info := createLargeFile(tempdir, v)
			fileEntries = append(fileEntries, FileEntry{Path: tempfile.Name(), Size: info.Size(), Name: info.Name()})
			fileEntries = append(fileEntries, FileEntry{Path: tempfile.Name(), Size: info.Size(), Name: info.Name()})
			fileEntries = append(fileEntries, FileEntry{Path: tempfile.Name(), Size: info.Size(), Name: info.Name()})
			fileEntries = append(fileEntries, FileEntry{Path: tempfile.Name(), Size: info.Size(), Name: info.Name()})
		}
		// time.Sleep(1 * time.Second)
		// log.Printf("\n\n>>>>> %v\n", fileEntries)
		time.Sleep(1 * time.Second)

		runtime.GOMAXPROCS(4)
		work := make(chan FileEntry)
		results := make(chan FileHashEntry)

		// create worker 5 goroutines
		wg := sync.WaitGroup{}
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for entry := range work {
					fileHashEntry := NewFileHashEntry(entry, HashConfig{})
					// log.Printf("%v\n", fileHashEntry)
					results <- fileHashEntry
				}
			}()
		}

		// send the work to the workers
		// this happens in a goroutine in order
		// to not block the main function, once
		// all 5 workers are busy
		go func() {
			for _, entry := range fileEntries {
				work <- entry
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
		allResults := []FileHashEntry{}
		for result := range results {
			// could write the file to disk
			allResults = append(allResults, result)
		}
		// log.Printf("allResults: %v\n", allResults)

	})
}
