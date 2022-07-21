package finder

import (
	"sort"
	"strconv"
)

type FormatConfig struct {
	Size bool
}

// convert a list of FileEntry to lines to be printed to console
// TODO: rename this to FileHashEntryFormatter
func DupesFormatter(dupes []FileHashEntry, config FormatConfig) string {
	var outputStr string
	lines := []string{}
	for _, entry := range dupes {
		var s string
		if config.Size {
			s = entry.Hash + "\t" + strconv.FormatInt(entry.File.Size, 10) + "\t" + entry.File.Path + "\n"
		} else {
			s = entry.Hash + "\t" + entry.File.Path + "\n"
		}
		lines = append(lines, s)
	}
	sort.Strings(lines)
	for _, line := range lines {
		outputStr += line
	}
	return outputStr
}

func FileEntryFormatter(dupes []FileEntry) string {
	var outputStr string
	lines := []string{}
	for _, entry := range dupes {
		var s string
		s = strconv.FormatInt(entry.Size, 10) + "\t" + entry.Path + "\n"
		lines = append(lines, s)
	}
	sort.Strings(lines)
	for _, line := range lines {
		outputStr += line
	}
	return outputStr
}
