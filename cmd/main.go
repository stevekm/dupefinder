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
	Profile    bool   `help:"enable profiling"`
}

func (cli *CLI) Run() error {
	err := run(cli.InputDir, cli.IgnoreFile, cli.PrintSize, cli.Parallel, cli.Profile)
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
	enableProfile bool) error {

	if enableProfile {
		cpuFile, memFile := finder.StartProfiler()
		defer cpuFile.Close()
		defer memFile.Close()
		defer pprof.StopCPUProfile()
	}

	// ignoreFile goes here
	hashConfig := finder.HashConfig{NumWorkers: numWorkers}
	formatConfig := finder.FormatConfig{Size: printSize}
	var skipDirs = []string{}
	dupes := finder.FindDupes(inputDir, skipDirs, hashConfig)
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
