package finder

import (
	"testing"
	"fmt"
	"log"
)

func TestHash(t *testing.T) {
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1, _ := createTempFile(tempdir, "f.", "writes\n")

	t.Run("Get md5 hash", func(t *testing.T) {
		hashConfig := HashConfig{}
		got := getFileMD5(tempfile1, hashConfig)
		want := "9d365f59076828add0b000414583cb33"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})

}

func TestHashSHA1(t *testing.T) {
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1, _ := createTempFile(tempdir, "f.", "writes\n")

	t.Run("Get sha1 hash", func(t *testing.T) {
		hashConfig := HashConfig{Algo:"sha1"}
		got := getFileMD5(tempfile1, hashConfig)
		want := "67503a007b3829965fde57d51768bdb32bb0389f"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})
}

func TestHashSHA256(t *testing.T) {
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1, _ := createTempFile(tempdir, "f.", "writes\n")

	t.Run("Get sha1 hash", func(t *testing.T) {
		hashConfig := HashConfig{Algo:"sha256"}
		got := getFileMD5(tempfile1, hashConfig)
		want := "fd6e46528c86f5f2a43aa9f013bf64fcc6939606e077bf3a4b14ef09fcb46f59"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})
}


func TestHashN(t *testing.T) {
	tempdir := t.TempDir()
	tempfile, _ := createLargeFile(tempdir, 4e5)
	// tempfileEntry := NewFileEntryFromPathInfo(tempfile.Name(), info)

	t.Run("Hash only the file head", func(t *testing.T){
		// hash the entire file
		hashConfig := HashConfig{}
		got := getFileMD5(tempfile, hashConfig)
		want := "d948f712fa329203f590e91cf6dd3e3e"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}

		// need to reset the cursor after writing
		i, err := tempfile.Seek(0, 0)
		if err != nil {
			fmt.Println("Error", err, i)
			log.Fatal(err)
		}

		// hash only the first 10 bytes
		got = getFileMD5(tempfile, HashConfig{Partial: true, NumBytes: 10})
		want = "a63c90cc3684ad8b0a2176a6a8fe9005"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})
}
