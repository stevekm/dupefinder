# dupefinder

Program to find duplicate files in a directory tree. 

`dupefinder` follows this basic search methodology;

- search the directory tree for files with the same byte size ("size dupes")
- search "size dupe" file list to find files with the same hash
- print results to console

`dupefinder` has several options availble to configure its behavior;

- hash multiple files in parallel (default 2)
- choose from different hashing algorithms (default md5, also available sha1, sha256, xxhash)
- print file byte sizes along with hashes for easier sorting
- consider only files that meet minimum or maximum file size parameters
- hash only the first `n` bytes of each file

`dupefinder` will automatically skip any files or directories that cannot be read or encounter errors. 

-----

Potential future features:
- supply a list of directory and file name patterns to exclude from search
- web interface

# Usage

```
./dupefinder /path/to/dir
```

Example:

```
$ ./dupefinder .
122641c2d78877cd166493bf15c80c4b	.git/refs/heads/master
122641c2d78877cd166493bf15c80c4b	.git/refs/remotes/origin/master
```

Include file sizes in output and sort to find the largest duplicates:

```
$ ./dupefinder --print-size ./ | sort -k2,2n
```

# Install

Download and run a pre-built binary from a release: https://github.com/stevekm/dupefinder/releases

Or compile from source, requires Go 1.17+ installed (https://go.dev/dl/).

```
git clone https://github.com/stevekm/dupefinder.git
cd dupefinder
make build
```
