package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	start := time.Now()
	parseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	errs := make(chan error)
	go handleErrors(errs)

	f, err := getFileDetails(*url)
	if err != nil {
		errs <- err
		return
	}

	if f.name == "" {
		f.name = generateRandomFileName(f)
		fmt.Printf("File name not found in Content-Disposition header. Proceeding with %s as file name.\n", f.name)
	}

	if !f.acceptRanges {
		fmt.Println("Server does not support range requests. Proceeding with single chunk download.")
		*numberOfChunks = 1
	}

	chunkSize := f.size / *numberOfChunks
	remainder := f.size % *numberOfChunks

	downloadFile(ctx, chunkSize, remainder, errs, f)

	fmt.Printf("Time elapsed: %.2f seconds. \n", time.Since(start).Seconds())
}

func handleErrors(errs <-chan error) {
	for err := range errs {
		fmt.Println(err.Error())
	}
}
