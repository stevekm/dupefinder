package finder

import (
  "strconv"
)

type FormatConfig struct {
  Size bool
}

// convert a list of FileEntry to lines to be printed to console
func DupesFormatter (hash string, dupes []FileEntry, config FormatConfig) string {
	var outputStr string
	for _, v := range dupes {
    var s string
    if config.Size {
      s = hash + "\t" + strconv.FormatInt(v.Size, 10) + "\t" + v.Path + "\n"
    } else {
      s = hash + "\t" + v.Path + "\n"
    }
		outputStr += s
	}
	return outputStr
}
