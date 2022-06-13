package finder

import (
	"testing"
)

func TestUtil(t *testing.T) {
	// setup test dirs & files
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end
	tempfile1, _ := createTempFile(tempdir, "f.", "writes\n")

	t.Run("Get md5 hash", func(t *testing.T) {
		got := getFileMD5(tempfile1)
		want := "9d365f59076828add0b000414583cb33"
		if got != want {
			t.Errorf("got %v is not the same as %v", got, want)
		}
	})

}
