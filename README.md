# dedup

A Golang CLI utility to deduplicate files.

## Install

### Go install

Run `go install` to download the binary to the go's binary folder:

```bash
go install github.com/nu12/dedup@latest
```

Note: go's binary folder (tipically `~/go/bin`) should be added to your PATH.

### From release

Download a tagged release binary for your OS (ubuntu, macos) placing it in a folder in your PATH and make it executable (may require elevated permissions):

```bash
wget -O /usr/local/bin/dedup https://github.com/nu12/dedup/releases/download/vX.Y.Z/dedup-linux-amd64.zip
unzip dedup-linux-amd64.zip
chmod +x dedup
mv dedup /usr/local/bin/dedup
```

Note: replace `X.Y.Z` with a valid version from the repository's releases and `linux-amd64` with the appropriate OS.

### From source

Clone this repo and compile the source code:

```bash
git clone github.com/nu12/dedup
cd dedup
go build -o dedup main.go
```

Move binary to a bin folder in your PATH (may require elevated permissions):
```bash
mv dedup /usr/local/bin/
```

## Usage

General usage for all commands is `dedup [flags]`. Find out all available commands with `dedup -h`:

```
Deduplicate files in a given directory.
Examples: 

List duplicated files in the current directory
dedup -s . --list

Move duplicated files from the current directory to another one
dedup -s . --move -d ../destination-folder

Usage:
  dedup [flags]
  dedup [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Show current version

Flags:
  -d, --destination string   Destination folder (for duplicated files)
  -h, --help                 help for dedup
      --list                 List duplicates
      --move                 Move duplicates
  -s, --source string        Source folder (to be dedup'ed)

Use "dedup [command] --help" for more information about a command.
```

## Release

```
git tag $(go run main.go version) && git push --tags
```