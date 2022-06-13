package finder

import (
	"os"
)

// handle the file opening and closing in order to get the file hash
func GetFileHash(fileEntry FileEntry) (FileHashEntry, error) {
	file, err := os.Open(fileEntry.Path)
	// if file read permission is denied, skip this file
	if os.IsPermission(err) {
		// logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", err)
		return FileHashEntry{}, err
	}

	if err != nil {
		// logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
		return FileHashEntry{}, err
	}
	hash := getFileMD5(file)
	file.Close()

	fileHashEntry := FileHashEntry{File: fileEntry, Hash: hash}
	return fileHashEntry, err
}

// find files that have the same hash value
func FindHashDupes(fileMap map[int64][]FileEntry) map[string][]FileHashEntry {
	hashesMap := map[string][]FileHashEntry{}
	for _, entries := range fileMap {
		for _, entry := range entries {
			fileHashEntry, err := GetFileHash(entry)
			if os.IsPermission(err) {
				logger.Printf("WARNING: Skipping file that could not be opened due to permissions error: %v\n", err)
				continue
			}

			if err != nil {
				logger.Printf("WARNING: Skipping file that could not be opened: %v\n", err)
				continue
			}
			hashesMap[fileHashEntry.Hash] = append(hashesMap[fileHashEntry.Hash], fileHashEntry)
		}
	}
	dupesMap := map[string][]FileHashEntry{}
	for hash, entries := range hashesMap {
		if len(entries) >1 {
			dupesMap[hash] = entries
		}
	}
	return dupesMap
}

// find all the duplicate files in the dir
// Duplicates = same file size, same hash value
// TODO: this might need to be broken up to aid garbage collection ??
func FindDupes(dirPath string, skipDirs []string) map[string][]FileHashEntry {
	fileSizeMap, _ := FindFilesSizes(dirPath, skipDirs)
	hashDupes := FindHashDupes(fileSizeMap)
	return hashDupes
}
