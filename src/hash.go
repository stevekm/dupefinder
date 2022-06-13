package finder

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// get the md5 hash of an open file handle
func getFileMD5(inputFile *os.File) string {
	hash := md5.New()
	if _, err := io.Copy(hash, inputFile); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)
	hashStr := hex.EncodeToString(sum[:])
	return hashStr
}
