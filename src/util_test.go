package finder

import (
	"fmt"
	"log"
	"testing"
	"os"
	"path/filepath"
	"github.com/google/go-cmp/cmp"
)


// NOTE: some of these methods are already depricated
func TestUtil(t *testing.T){
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1, err := os.CreateTemp(tempdir, "file1.") // *os.File
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile1.Close()

	// write to the file
	nbytesWritten, err := tempfile1.WriteString("writes\n")
	if err != nil {
		fmt.Println(nbytesWritten)
		log.Fatal(err)
	}

	// need to reset the cursor after writing
	i, err := tempfile1.Seek(0, 0)
	if err != nil {
		fmt.Println("Error", err, i)
		log.Fatal(err)
	}

	t.Run("Get md5 hash", func(t *testing.T) {
		got := getFileMD5(tempfile1)
		want := "9d365f59076828add0b000414583cb33"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})

	t.Run("Get file size", func(t *testing.T) {
		got := getFileSize(tempfile1)
		var want int64 = 7
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})
}







func TestFinder(t *testing.T) {
	// set up temp dirs for files
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
	tempfile1, err := os.CreateTemp(subdir1, "file1.") // *os.File
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile1.Close()

	tempfile2, err := os.CreateTemp(subdir2, "file2.")
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile2.Close()

	tempfile3, err := os.CreateTemp(tempdir, "file3.")
	if err != nil {
		log.Fatal(err)
	}
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

	tempfile5, err := os.CreateTemp(subdir3, "file5.")
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile5.Close()

	// write to the file
	nbytesWritten, err := tempfile1.WriteString("writes\n")
	if err != nil {
		fmt.Println(nbytesWritten)
		log.Fatal(err)
	}

	t.Run("Test find dupes", func(t *testing.T) {
		var skipDirs = []string{subdir3}
		got := FindDupes(tempdir, skipDirs)
		// fmt.Printf("Duplicate file sizes: %v\n", got)
		want := map[string][]string{
			"d41d8cd98f00b204e9800998ecf8427e": []string{
				tempfile3.Name(),
				tempfile2.Name(),
				tempfile4.Name(),
			},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

		for hash, paths := range got {
			gotFormat := DupesFormatter(hash, paths)
			// wantFormat :=
			// if diff := cmp.Diff(want, gotFormat); diff != "" {
			// 	t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
			// }
			fmt.Println(gotFormat)
		}






		// got := IsUncalledFilter(&mutation)
		// want := true
		// if !cmp.Equal(got, want) {
		// 	t.Errorf("got %v is not the same as %v", got, want)
		// }
		// if diff := cmp.Diff(want, got); diff != "" {
		// 	t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		// }

	})
}
