package finder

import (
	"fmt"
	"log"
	"testing"
	"os"
	"path/filepath"
	"strconv"
	"github.com/google/go-cmp/cmp"
)

// $ DIR_TEST=True make test
func skipDirTest(t *testing.T){
	var dirTestVar = os.Getenv("DIR_TEST")
	var runDirTest bool
	if dirTestVar != "" {
		runDirTest = true
	}
	if ! runDirTest {
    t.Skip(">>> Skipping dir test")
  }
}

// create a temp file in a dir and write something to its contents
func createTempFile(tempdir string, filename string, contents string) *os.File {
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

	return tempfile
}

// generate a list of numbers
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
			a[i] = min + i
		}
	return a
}



func TestUtil(t *testing.T){
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1 := createTempFile(tempdir, "f.", "writes\n")

	t.Run("Get md5 hash", func(t *testing.T) {
		got := getFileMD5(tempfile1)
		want := "9d365f59076828add0b000414583cb33"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})

}







func TestFinder(t *testing.T) {
	// set up temp dirs for tests
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	subdir1 := filepath.Join(tempdir, "subdir.1")
	err := os.MkdirAll(subdir1, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	subdir2 := filepath.Join(tempdir, "subdir.2")
	err = os.MkdirAll(subdir2, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	subdir3 := filepath.Join(tempdir, "subdir.3")
	err = os.MkdirAll(subdir3, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// set up tempfiles for the test
	// tempfile1 has unique contents
	// tempfile2 and tempfile3 have different names but same contents (empty)
	// tempfile3 and tempfile4 have same names and same contents (empty) but different directories
	// tempfile5 is in the subdir to skip and has same size as tempfile2, tempfile3

	tempfile1 := createTempFile(subdir1, "file1.", "writes\n")
	defer tempfile1.Close()

	tempfile2 := createTempFile(subdir2, "file2.", "")
	defer tempfile2.Close()

	tempfile3 := createTempFile(tempdir, "file3.", "")
	defer tempfile3.Close()

	fi, err := tempfile3.Stat()
	if err != nil {
		log.Fatal(err)
	}
	tempfile3Basename := fi.Name()

	tempfile4Path := filepath.Join(subdir2, tempfile3Basename)
	tempfile4, err := os.Create(tempfile4Path)
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile4.Close()

	tempfile5 := createTempFile(subdir3, "file5.", "")
	defer tempfile5.Close()

	t.Run("Test find dupes", func(t *testing.T) {
		var skipDirs = []string{subdir3}
		got := FindDupes(tempdir, skipDirs)
		wantHash := "d41d8cd98f00b204e9800998ecf8427e"
		want := map[string][]FileEntry{
			wantHash: []FileEntry{
				NewFileEntryFromPath(tempfile3.Name()),
				NewFileEntryFromPath(tempfile2.Name()),
				NewFileEntryFromPath(tempfile4.Name()),
			},
		}
		// test that we found the expected duplicate files
		// NOTE: might need to revise this test to not depend on order of items in the list!
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

		// test that the console formatter prints them in the expected format
		config := FormatConfig{Size: false}
		gotFormat := DupesFormatter(wantHash, got[wantHash], config)
		var wantFormat string
		for _, entry := range want[wantHash] {
			wantFormat += wantHash + "\t" + entry.Path + "\n"
		}

		if diff := cmp.Diff(wantFormat, gotFormat); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

	})
}



func TestTooManyFiles(t *testing.T){
	// only run this test if DIR_TEST env var was enabled because it creates a lot of files
	skipDirTest(t)

	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end

	t.Run("Find dupes without exceeding the limit on number of open files", func(t *testing.T) {
		// make a large number of temp files, each with different contents
		ints := makeRange(0, 20000)
		for i := range ints {
			i_str := strconv.Itoa(i)
			t := createTempFile(tempdir, i_str, i_str)
			t.Close()
		}
		// create two temp files with the same contents
		tempfile1 := createTempFile(tempdir, "f.", "foo")
		defer tempfile1.Close()
		tempfile2 := createTempFile(tempdir, "f2.", "foo")
		defer tempfile2.Close()

		var skipDirs = []string{}
		got := FindDupes(tempdir, skipDirs)
		want := map[string][]string{
			"acbd18db4cc2f85cedef654fccc4a4d8": []string{
				tempfile1.Name(),
				tempfile2.Name(),
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}
	})
}


func TestPermissionsError(t *testing.T){
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end

	// create two temp files with the same contents
	tempfile1 := createTempFile(tempdir, "f.", "foo")
	defer tempfile1.Close()
	tempfile2 := createTempFile(tempdir, "f2.", "foo")
	defer tempfile2.Close()

	// remove read permissions from one file
	err := os.Chmod(tempfile2.Name(), 0000)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("Find dupes while avoiding files with permissions errors", func(t *testing.T) {
		var skipDirs = []string{}
		got := FindDupes(tempdir, skipDirs)
		want := map[string][]FileEntry{}
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
		tempfile3 := createTempFile(subdir1, "f3.", "foo\n")
		tempfile4 := createTempFile(subdir1, "f4.", "foo\n")

		// remove read permissions from the directory file
		err = os.Chmod(subdir2, 0000)
		if err != nil {
			log.Fatal(err)
		}

		var skipDirs = []string{}
		got := FindDupes(subdir1, skipDirs)
		want := map[string][]FileEntry{
			"d3b07384d113edec49eaa6238ad5ff00": []FileEntry{
				NewFileEntryFromPath(tempfile3.Name()),
				NewFileEntryFromPath(tempfile4.Name()),
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
