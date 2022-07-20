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
	HashBytes int64  `help:"number of bytes to hash for each duplicated file; example: 1000 = 1KB, 1000000 = 1MB, 1000000000 = 1GB"`
	Algo      string `help:"hashing algorithm to use. Options: md5, sha1, sha256, xxhash" default:"md5"`
	SizeOnly bool `help:"only look for duplicates based on file size"`
	MinSize   int64  `help:"only include files of minimum size (bytes) or larger when searching"`
	// NOTE: note sure how to get Kong to accept type of *int64 here for MaxSize;
	MaxSize   int64  `help:"only include files of maximum size (bytes) or smaller when searching. Value must be >0, value of 0 = disabled" default:"0"`
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
		cli.SizeOnly,
		cli.MaxSize,
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
	sizeOnly bool,
	maxSize int64,
	debug bool,
) error {

	if enableProfile {
		cpuFile, memFile := finder.StartProfiler()
		defer cpuFile.Close()
		defer memFile.Close()
		defer pprof.StopCPUProfile()
	}

	findConfig := finder.FindConfig{MinSize: minSize} // var skipDirs = []string{} // ignoreFile goes here

	// NOTE: note sure how to get Kong to accept type of *int64 here for MaxSize
	if maxSize > 0 {
		findConfig.MaxSize = &maxSize
	}

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

	if sizeOnly {
		fileSizeMap, _ := finder.FindFilesSizes(inputDir, findConfig)
		sizeDupes := finder.FindSizeDupes(fileSizeMap)
		for _, entries := range sizeDupes {
			format := finder.FileEntryFormatter(entries)
			fmt.Printf("%s", format)
		}
	} else {
		dupes := finder.FindDupes(inputDir, findConfig, hashConfig)
		for _, entries := range dupes {
			format := finder.DupesFormatter(entries, formatConfig)
			fmt.Printf("%s", format) // format has newline embedded at the end
		}
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
