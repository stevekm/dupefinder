package finder

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFindAllFiles(t *testing.T) {
	// set up temp dirs for tests
	tempdir := t.TempDir() // automatically gets cleaned up when all tests end

	t.Run("Test find all files", func(t *testing.T) {
		tempDirs, tempFiles := createTempFilesDirs1(tempdir)
		var skipDirs = []string{tempDirs[2]}

		gotFiles, gotNumFiles := FindFilesSizes(tempdir, skipDirs)

		wantFiles := map[int64][]FileEntry{
			0: []FileEntry{
				NewFileEntryFromPath(tempFiles[2].Name()),
				NewFileEntryFromPath(tempFiles[1].Name()),
				NewFileEntryFromPath(tempFiles[3].Name()),
			},
			7: []FileEntry{
				NewFileEntryFromPath(tempFiles[0].Name()),
			},
		}

		var wantNumFiles uint64 = 4

		// test that we found the expected files
		// NOTE: might need to revise this test to not depend on order of items in the list!
		if diff := cmp.Diff(wantFiles, gotFiles); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(wantNumFiles, gotNumFiles); diff != "" {
			t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
		}

	})
}
