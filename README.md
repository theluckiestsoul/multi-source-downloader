# multi-source-downloader

Multi-source downloader is a simple tool to download a single file in chunks from a single source. It is written in `go 1.21.5` and uses `go modules` for dependency management.

## Setup and Run Locally
Install go 1.21.5 or higher. Clone the repository and run the following command in the root directory of the project:

### Build
```bash
make build
```

### Install
```bash
make install
```

Alternatively, you can run the following command to install the binary in your `$GOPATH/bin` directory:
```bash
go install github.com/theluckiestsoul/multi-source-downloader@latest
```

Once installed, you can run the binary using the following command:
```bash
multi-source-downloader -url https://examplefile.com/file-download/25
```

## Options
The following options are available:
```bash
Usage of multi-source-downloader:
  -url string (required)
        URL to download the file from
  -chunks int (optional)
        Number of chunks to download concurrently (default 8)
```
