package finder

// https://pkg.go.dev/runtime/pprof
// https://github.com/google/pprof/blob/main/doc/README.md

// $ go tool pprof cpu.prof
// $ go tool pprof mem.prof

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func StartProfiler() (*os.File, *os.File) {
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	// defer cpuFile.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	// defer pprof.StopCPUProfile()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	// defer memFile.Close() // error handling omitted for example
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	return cpuFile, memFile
}
