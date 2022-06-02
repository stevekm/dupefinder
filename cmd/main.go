package main

import (
  "log"
  "fmt"
  "github.com/alecthomas/kong"
  "dupefinder/src" // "dupefinder/src" as finder
)

type CLI struct {
  InputDir string `help:"path to input file to search" arg:""`
  IgnoreFile string `help:"path to file of dir paths to ignore"`
}

func (cli *CLI) Run () error {
  err := run(cli.InputDir, cli.IgnoreFile)
  if err != nil {
		log.Fatalln(err)
	}
  return nil
}

func run (inputDir string, ignoreFile string) error {
  // ignoreFile goes here
  var skipDirs = []string{}
  dupes := finder.FindDupes(inputDir, skipDirs)
  for hash, paths := range dupes {
    format := finder.DupesFormatter(hash, paths)
    // fmt.Println(format)
    fmt.Printf("%s", format) // format has newline embedded at the end
  }
  return nil
}

func main () {
  var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("Duplicate File Finder"),
		kong.Description("Program for finding duplicate files in a directory"))

	ctx.FatalIfErrorf(ctx.Run())

}
