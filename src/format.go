package finder

import (
	"strconv"
)

type FormatConfig struct {
	Size bool
}

// convert a list of FileEntry to lines to be printed to console
func DupesFormatter(dupes []FileHashEntry, config FormatConfig) string {
	var outputStr string
	for _, entry := range dupes {
		var s string
		if config.Size {
			s = entry.Hash + "\t" + strconv.FormatInt(entry.File.Size, 10) + "\t" + entry.File.Path + "\n"
		} else {
			s = entry.Hash + "\t" + entry.File.Path + "\n"
		}
		outputStr += s
	}
	return outputStr
}
