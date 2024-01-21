package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

var (
	url            = flag.String("url", "https://examplefile.com/file-download/25", "file url to download")
	numberOfChunks = flag.Int("chunks", 8, "number of chunks to download concurrently")
)

func init() {
	const usage = `Usage of multi-source-downloader %s:

    multi-source-downloader [options] -url <file url> -chunks <number of chunks> -verify

Options:
    - url: file url to download (required)
    - chunks: number of chunks to download concurrently (default %d)
    - verify: verify file integrity (default false)
`
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, getVersion(), *numberOfChunks)
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	validateFlags()
}

func validateFlags() {
	if *url == "" {
		fmt.Println("URL is required")
		flag.Usage()
	}
	if *numberOfChunks < 1 {
		fmt.Println("Number of chunks must be greater than 0")
		flag.Usage()
	}
}

func getVersion() string {
	if i, ok := debug.ReadBuildInfo(); ok {
		return i.Main.Version
	}

	return "(unknown)"
}
