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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	errs := make(chan error)
	defer close(errs)

	go func() {
		for err := range errs {
			fmt.Println(err.Error())
		}
	}()

	f, err := getFileDetails(*url)
	if err != nil {
		errs <- err
		return
	}

	if f.name == "" {
		fmt.Println("File name not found in Content-Disposition header. Proceeding with default file name.")
		f.name = generateRandomFileName()
	}

	if !f.acceptRanges {
		fmt.Println("Server does not support range requests. Proceeding with single chunk download.")
		*numberOfChunks = 1
	}

	chunkSize := f.size / *numberOfChunks
	remainder := f.size % *numberOfChunks

	downloadFile(ctx, chunkSize, remainder, errs, f)

	fmt.Printf("Time elapsed: %v\n", time.Since(start).Seconds())
}
