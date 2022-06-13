package finder

import (
	"os"
	"log"
)

// logger to use throughout the package
var logger = log.New(os.Stderr, "", 0)
