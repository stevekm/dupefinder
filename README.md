# dupefinder
find duplicate files in a directory

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

# Install

Download and run a pre-built binary from a release: https://github.com/stevekm/dupefinder/releases

Or compile from source, requires Go 1.17+ installed (https://go.dev/dl/).

```
git clone https://github.com/stevekm/dupefinder.git
cd dupefinder
make build
```

