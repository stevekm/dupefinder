package finder

import (
	"log"
	"os"
)

// https://pkg.go.dev/log
// https://pkg.go.dev/log#Logger.SetFlags
// https://pkg.go.dev/log#pkg-constants

// logger to use throughout the package
var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
