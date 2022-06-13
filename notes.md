
Make a program to track files in the dir tree (disk volume) and look for duplicates

secondary goal; report files that have changed or are missing by tracking file dir tree

-----
Reference:

https://stackoverflow.com/questions/53314863/fastest-algorithm-to-detect-duplicate-files


So the usual algorithm seems to work like this:

    generate a sorted list of all files (path, Size, id)
    group files with the exact same size
    calculate the hash of all the files with a same size and compare the hashes
    same has means identical files - a duplicate is found

Sometimes the speed gets increased by first using a faster hash algorithm (like md5) with more collision probability and second if the hash is the same use a second slower but less collision-a-like algorithm to prove the duplicates. Ano


-----
https://stackoverflow.com/questions/1761607/what-is-the-fastest-hash-algorithm-to-check-if-two-files-are-equal

So to compare two files, use this algorithm:

    Compare sizes
    Compare dates (be careful here: this can give you the wrong answer; you must test whether this is the case for you or not)
    Compare the hashes

-----

Tasks:

- recurse dir tree to find all files
  - get path, basename, byte size, timestamp for each file

- from list of all files, find matches based on:
  - same basename (same file in two dirs)
    - compare sizes


------
// https://pkg.go.dev/path/filepath#Walk
// TODO: look into using https://pkg.go.dev/io/fs#WalkDirFunc , https://pkg.go.dev/path/filepath#WalkDir
// https://github.com/golang/go/issues/16399
// https://pkg.go.dev/io/fs#DirEntry
// NOTE: dont think WalkDir will save much time since I want the extra info anyway??
// https://pkg.go.dev/io/fs#FileInfo
// https://golang.hotexamples.com/examples/os/-/IsPermission/golang-ispermission-function-examples.html
