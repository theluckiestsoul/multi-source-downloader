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
multi-source-downloader -url https://link.testfile.org/PDF100MB
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

## Sample file URLs for testing
The following URLs can be used to test the downloader:
- https://link.testfile.org/PDF10MB
- https://link.testfile.org/PDF20MB
- https://link.testfile.org/PDF30MB
- https://link.testfile.org/PDF40MB
- https://link.testfile.org/PDF50MB
- https://link.testfile.org/PDF100MB
- https://link.testfile.org/PDF200MB
- https://examplefile.com/file-download/25

> [!NOTE]  
> Currently, the downloader only supports downloading files from a single source. The file is downloaded in chunks from the same source. The downloader does not support downloading a single file in chunks from multiple sources.

> [!IMPORTANT] 
> The downloader stores the downloaded file in the current working directory. By default the file uses the name provided in the `Content-Disposition` header. If the header is not present, the file is saved with a random name.

> [!WARNING] 
> The `Etag` header verification is not implemented.