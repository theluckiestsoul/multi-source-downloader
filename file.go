package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	chunk struct {
		index        int
		tempFileName string
	}

	file struct {
		name         string
		size         int
		etag         string
		acceptRanges bool
	}
)

var (
	client = http.DefaultClient
)

func downloadFile(ctx context.Context, chunkSize int, remainder int, errs chan error, f *file) {
	fmt.Printf("Downloading %s \n", f.name)
	chunks := make(chan chunk, *numberOfChunks)

	var eg errgroup.Group

	for i := 0; i < *numberOfChunks; i++ {
		index := i
		eg.Go(func() error {
			start := index * chunkSize
			end := start + chunkSize - 1
			if index == *numberOfChunks-1 {
				end += remainder
			}
			tmpFileName, err := downloadChunk(ctx, *url, start, end)
			if err != nil {
				return err
			}
			select {
			case chunks <- chunk{index: index, tempFileName: tmpFileName}:
			case <-ctx.Done():
				fmt.Println("Operation canceled")
				return ctx.Err()
			}
			return nil
		})
	}
	go func() {
		err := eg.Wait()
		close(chunks)
		if err != nil {
			errs <- err
		}
	}()

	fileNames := make([]string, *numberOfChunks)
	defer cleanupFiles(fileNames, errs)

	for chunk := range chunks {
		fileNames[chunk.index] = chunk.tempFileName
		fmt.Printf("Chunk %s downloaded\n", chunk.tempFileName)
	}
	if err := mergeFiles(f.name, fileNames); err != nil {
		errs <- err
	} else {
		fmt.Println("File downloaded successfully.")
	}

}

func getFileDetails(url string) (*file, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	cd := resp.Header.Get("Content-Disposition")
	index := strings.Index(cd, "filename=")
	var filename string
	if index > -1 {
		filename = cd[index+len("filename="):]
	}

	etag := resp.Header.Get("ETag")

	acceptRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	f := &file{
		name:         filename,
		size:         size,
		etag:         etag,
		acceptRanges: acceptRanges,
	}

	return f, nil
}

func downloadChunk(ctx context.Context, url string, start, end int) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", rangeHeader)
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return writeToFile(resp.Body)
}

func writeToFile(body io.ReadCloser) (string, error) {
	tmpFile, err := os.CreateTemp("", "chunk")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func mergeFiles(fileName string, fileNames []string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, tmpFileName := range fileNames {
		tmpFile, err := os.OpenFile(tmpFileName, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer tmpFile.Close()

		_, err = io.Copy(file, tmpFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupFiles(files []string, errs chan error) {
	fmt.Println("Cleaning up temporary files...")
	for _, file := range files {
		if file == "" {
			continue
		}
		fmt.Printf("Removing %s\n", file)
		err := os.Remove(file)
		if err != nil {
			errs <- err
		}
	}
}

func generateRandomFileName() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
