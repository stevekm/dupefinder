# dupefinder
find duplicate files in a directory

# Install

Compiling from source requires Go version 1.17+ installed (https://go.dev/dl/).

```
git clone https://github.com/stevekm/dupefinder.git
cd dupefinder
make build
```

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
