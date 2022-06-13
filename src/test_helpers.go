package finder

import (
	"os"
	"log"
	"fmt"
	"path/filepath"
)

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

// generate a list of numbers
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
