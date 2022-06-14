package finder

import (
	"log"
	"os"
)

// logger to use throughout the package
var logger = log.New(os.Stderr, "", 0)

// logger.Fatalf
