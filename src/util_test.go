package finder

import (
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// skip the test with the big dir with lots of files;
// $ DIR_TEST=True make test
func skipDirTest(t *testing.T) {
	var dirTestVar = os.Getenv("DIR_TEST")
	var runDirTest bool
	if dirTestVar != "" {
		runDirTest = true
	}
	if !runDirTest {
		t.Skip(">>> Skipping dir test")
	}
}

func TestFinder(t *testing.T) {
	// set up temp dirs for tests
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end

	t.Run("Test find dupes", func(t *testing.T) {
		tempDirs, tempFiles := createTempFilesDirs1(tempdir)
		hashConfig := HashConfig{NumWorkers: 2}
		findConfig := FindConfig{SkipDirs: []string{tempDirs[2]}}
		got := FindDupes(tempdir, findConfig, hashConfig)
		wantHash := "d41d8cd98f00b204e9800998ecf8427e"
		want := map[string][]FileHashEntry{
			wantHash: []FileHashEntry{
				NewFileHashEntry(NewFileEntryFromPath(tempFiles[2].Name()), hashConfig),
				NewFileHashEntry(NewFileEntryFromPath(tempFiles[1].Name()), hashConfig),
				NewFileHashEntry(NewFileEntryFromPath(tempFiles[3].Name()), hashConfig),
			},
		}
		// test that we found the expected duplicate files
		if len(got) != len(want) {
			t.Errorf("got %v is not the same as %v", got, want)
		}
		if len(got[wantHash]) != len(want[wantHash]) {
			t.Errorf("got %v is not the same as %v", got, want)
		}
		for _, entry := range want[wantHash] {
			if !containsFileHashEntry(got[wantHash], entry) {
				t.Errorf("%v not in list %v", entry, got[wantHash])
			}
		}

		// test that the console formatter prints them in the expected format
		config := FormatConfig{Size: false}
		gotFormat := DupesFormatter(got[wantHash], config)
		var wantFormat string
		for _, entry := range want[wantHash] {
			wantFormat += wantHash + "\t" + entry.File.Path + "\n"
		}

		if diff := cmp.Diff(wantFormat, gotFormat); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

	})
}

func TestTooManyFiles(t *testing.T) {
	// only run this test if DIR_TEST env var was enabled because it creates a lot of files
	skipDirTest(t)

	// setup test dirs & files
	tempdir := t.TempDir()

	t.Run("Find dupes without exceeding the limit on number of open files", func(t *testing.T) {
		// make a large number of temp files, each with different contents
		ints := makeRange(0, 20000)
		for i := range ints {
			i_str := strconv.Itoa(i)
			t, _ := createTempFile(tempdir, i_str, i_str)
			t.Close()
		}
		// create two temp files with the same contents
		tempfile1, _ := createTempFile(tempdir, "f.", "foo")
		defer tempfile1.Close()
		tempfile2, _ := createTempFile(tempdir, "f2.", "foo")
		defer tempfile2.Close()

		// var skipDirs = []string{}
		hashConfig := HashConfig{NumWorkers: 2}
		findConfig := FindConfig{SkipDirs: []string{}}
		got := FindDupes(tempdir, findConfig, hashConfig)
		want := map[string][]FileHashEntry{
			"acbd18db4cc2f85cedef654fccc4a4d8": []FileHashEntry{
				NewFileHashEntry(NewFileEntryFromPath(tempfile1.Name()), hashConfig),
				NewFileHashEntry(NewFileEntryFromPath(tempfile2.Name()), hashConfig),
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestPermissionsError(t *testing.T) {
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end

	// create two temp files with the same contents
	tempfile1, _ := createTempFile(tempdir, "f.", "foo")
	defer tempfile1.Close()
	tempfile2, _ := createTempFile(tempdir, "f2.", "foo")
	defer tempfile2.Close()

	// remove read permissions from one file
	err := os.Chmod(tempfile2.Name(), 0000)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Find dupes while avoiding files with permissions errors", func(t *testing.T) {
		findConfig := FindConfig{SkipDirs: []string{}}
		hashConfig := HashConfig{NumWorkers: 2}
		got := FindDupes(tempdir, findConfig, hashConfig)
		want := map[string][]FileHashEntry{}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Find dupes while skipping directories with permissions errors", func(t *testing.T) {
		// make a subdir to hold some files
		subdir1 := filepath.Join(tempdir, "subdir.1")
		err := os.MkdirAll(subdir1, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// make another sub-subdir that we are gonna mess with permissions on
		subdir2 := filepath.Join(tempdir, "subdir.2")
		err = os.MkdirAll(subdir2, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		// add some files with duplicate contents in the sub-subdir
		createTempFile(subdir2, "f.", "foo\n")
		createTempFile(subdir2, "f2.", "foo\n")

		// add similar files in the parent dir
		tempfile3, _ := createTempFile(subdir1, "f3.", "foo\n")
		tempfile4, _ := createTempFile(subdir1, "f4.", "foo\n")

		// remove read permissions from the directory file
		err = os.Chmod(subdir2, 0000)
		if err != nil {
			log.Fatal(err)
		}

		hashConfig := HashConfig{NumWorkers: 2}
		findConfig := FindConfig{SkipDirs: []string{}}
		got := FindDupes(subdir1, findConfig, hashConfig)
		want := map[string][]FileHashEntry{
			"d3b07384d113edec49eaa6238ad5ff00": []FileHashEntry{
				NewFileHashEntry(NewFileEntryFromPath(tempfile3.Name()), hashConfig),
				NewFileHashEntry(NewFileEntryFromPath(tempfile4.Name()), hashConfig),
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

		// fix the permissons so we can cleanup
		err = os.Chmod(subdir2, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	})
}

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
