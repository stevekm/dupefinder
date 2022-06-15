package main

import (
	"dupefinder/src" // "dupefinder/src" as finder
	"fmt"
	"github.com/alecthomas/kong"
	"log"
	"runtime/pprof"
)

type CLI struct {
	InputDir   string `help:"path to input file to search" arg:""`
	IgnoreFile string `help:"path to file of dir paths to ignore"`
	PrintSize  bool   `help:"print the file size"`
	Parallel   int    `help:"number of items to process in parallel" default:"2"`
	Profile    bool   `help:"enable profiling and outputs files for use with
'go tool pprof cpu.prof'"`
	HashBytes int64  `help:"number of bytes to hash for each duplicated file; example: 500000 = 500KB, 1000000 = 1MB"`
	Algo      string `help:"hashing algorithm to use. Options: md5, sha1, sha256, xxhash" default:"md5"`
	MinSize   int64  `help:"only include files of minimum size (bytes) when searching"`
	Debug     bool   `help:"only used for dev debug purposes! Don't use this option it doesnt do anything"`
}

func (cli *CLI) Run() error {
	err := run(
		cli.InputDir,
		cli.IgnoreFile,
		cli.PrintSize,
		cli.Parallel,
		cli.Profile,
		cli.HashBytes,
		cli.Algo,
		cli.MinSize,
		cli.Debug,
	)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func run(
	inputDir string,
	ignoreFile string,
	printSize bool,
	numWorkers int,
	enableProfile bool,
	hashBytes int64,
	algo string,
	minSize int64,
	debug bool,
) error {

	if enableProfile {
		cpuFile, memFile := finder.StartProfiler()
		defer cpuFile.Close()
		defer memFile.Close()
		defer pprof.StopCPUProfile()
	}

	findConfig := finder.FindConfig{MinSize: minSize} // var skipDirs = []string{} // ignoreFile goes here
	hashConfig := finder.HashConfig{NumWorkers: numWorkers, Algo: algo}
	if hashBytes > 0 {
		hashConfig.Partial = true
		hashConfig.NumBytes = hashBytes
	}

	formatConfig := finder.FormatConfig{Size: printSize}

	if debug {
		// change the commands here to use when debugging and benchmarking stuff, etc..
		finder.FindFilesSizes(inputDir, findConfig)
		return nil
	}

	dupes := finder.FindDupes(inputDir, findConfig, hashConfig)
	for _, entries := range dupes {
		format := finder.DupesFormatter(entries, formatConfig)
		fmt.Printf("%s", format) // format has newline embedded at the end
	}
	return nil
}

func main() {
	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("Duplicate File Finder"),
		kong.Description("Program for finding duplicate files in a directory"))

	ctx.FatalIfErrorf(ctx.Run())

}
