package finder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// 
// HELPER FUNCTIONS AND METHOD FOR USE WITH TEST SUITE
// 

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


// create a temp file in a dir and write something to its contents
func createTempFile(tempdir string, filename string, contents string) (*os.File, string) {
	tempfile, err := os.CreateTemp(tempdir, filename)
	if err != nil {
		log.Fatal(err)
	}
	// defer tempfile.Close()

	// write to the file
	if contents != "" {
		nbytesWritten, err := tempfile.WriteString(contents)
		if err != nil {
			fmt.Println(nbytesWritten)
			log.Fatal(err)
		}

		// need to reset the cursor after writing
		i, err := tempfile.Seek(0, 0)
		if err != nil {
			fmt.Println("Error", err, i)
			log.Fatal(err)
		}
	}

	// get the randomly generated file basename
	fi, err := tempfile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	basename := fi.Name()

	return tempfile, basename
}

func createSubDir(tempdir string, filename string) string {
	subdir := filepath.Join(tempdir, filename)
	err := os.MkdirAll(subdir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return subdir
}

// set up a bunch of temp files and subdirs to use in test cases
// tempfile1 has unique contents
// tempfile2 and tempfile3 have different names but same contents (empty)
// tempfile3 and tempfile4 have same names and same contents (empty) but different directories
// tempfile5 is in the subdir to skip and has same size as tempfile2, tempfile3
func createTempFilesDirs1(tempdir string) ([]string, []*os.File) {
	subdir1 := createSubDir(tempdir, "subdir.1")
	subdir2 := createSubDir(tempdir, "subdir.2")
	subdir3 := createSubDir(tempdir, "subdir.3")

	tempfile1, _ := createTempFile(subdir1, "file1.", "writes\n")
	// defer tempfile1.Close()

	tempfile2, _ := createTempFile(subdir2, "file2.", "")
	// defer tempfile2.Close()

	tempfile3, tempfile3Basename := createTempFile(tempdir, "file3.", "")
	// defer tempfile3.Close()

	tempfile4, _ := createTempFile(subdir2, tempfile3Basename, "")
	// defer tempfile4.Close()

	tempfile5, _ := createTempFile(subdir3, "file5.", "")
	// defer tempfile5.Close()

	tempDirs := []string{subdir1, subdir2, subdir3}
	tempFiles := []*os.File{tempfile1, tempfile2, tempfile3, tempfile4, tempfile5}

	return tempDirs, tempFiles
}

func createLargeFile(tempdir string, size int64) (*os.File, os.FileInfo) {
	// for _, v := range []int64{3e8, 4e8, 1e8, 1e7, 4e5, 8e4,3e8, 4e8,} {
	// 	tempfile, info := createLargeFile(tempdir, v)
	// 	fileEntries = append(fileEntries, FileEntry{Path: tempfile.Name(), Size: info.Size(), Name: info.Name()})
	// }
	tempfile, err := os.CreateTemp(tempdir, "fooo.txt")
	if err != nil {
		log.Fatal(err)
	}

	if err := tempfile.Truncate(size); err != nil {
		log.Fatal(err)
	}

	info, err := tempfile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return tempfile, info
}

// generate a list of numbers
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
